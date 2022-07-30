package net

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/indeedhat/harmony/internal/events"
	. "github.com/indeedhat/harmony/internal/logger"
	"github.com/indeedhat/harmony/internal/screens"
)

type Client struct {
	// Events coming from the server
	Events chan []byte

	ctx  *common.Context
	ws   *websocket.Conn
	uuid uuid.UUID
}

// NewClient harmony client
func NewClient(ctx *common.Context, uuid uuid.UUID, ip string, screens []screens.DisplayBounds) (*Client, error) {
	serverAddress := fmt.Sprintf("%s:%d", ip, config.ServerPort)
	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/ws"}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := Client{
		uuid:   uuid,
		ctx:    ctx,
		ws:     ws,
		Events: make(chan []byte),
	}

	if err := client.sendConnect(screens); err != nil {
		client.Close()
		return nil, err
	}

	go client.readEventsFromServer()

	return &client, nil
}

// Close the client
func (cnt *Client) Close() error {
	return cnt.ws.Close()
}

// SendMessage over the websocket connection
func (cnt *Client) SendMessage(msg events.WsMessage) error {
	data, err := msg.Marshal()
	if err != nil {
		Logf("client", "failed to marshal event: %s", err)
		return err
	}

	return cnt.ws.WriteMessage(websocket.BinaryMessage, data)
}

func (cnt *Client) sendConnect(screens []screens.DisplayBounds) error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	msg := &events.ClientConnect{
		Hostname: hostname,
		UUID:     cnt.uuid,
		Displays: screens,
	}

	return cnt.SendMessage(msg)
}

// readEventsFromServer and pass the hid events out to the application via the InputEvents chanel
func (cnt *Client) readEventsFromServer() {
	for {
		_, data, err := cnt.ws.ReadMessage()
		if err != nil {
			if errors.Is(err, websocket.ErrCloseSent) {
				cnt.ctx.Cancel()
				return
			}
			Logf("client", "read error: %s", err)
			continue
		}

		Logf("client", "event recieved: %d", data[0])
		cnt.Events <- data
	}
}
