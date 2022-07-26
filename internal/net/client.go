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
)

type Client struct {
	// Events coming from the server
	Events chan []byte
	// Input for events to be sent to the server
	Input chan common.WsMessage

	ctx  *common.Context
	ws   *websocket.Conn
	uuid uuid.UUID
}

// NewClient harmony client
func NewClient(ctx *common.Context, uuid uuid.UUID, ip string) (*Client, error) {
	serverAddress := fmt.Sprintf("%s:%d", ip, config.ServerPort)
	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/ws"}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := Client{
		ctx:    ctx,
		ws:     ws,
		Events: make(chan []byte),
		Input:  make(chan common.WsMessage),
	}

	if err := client.sendConnect(); err != nil {
		client.Close()
		return nil, err
	}

	go client.handleIncommingEvents()
	go client.readEventsFromServer()

	return &client, nil
}

// Close the client
func (cnt *Client) Close() error {
	return cnt.ws.Close()
}

// SendMessage over the websocket connection
func (cnt *Client) SendMessage(msg common.WsMessage) error {
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	return cnt.ws.WriteMessage(websocket.BinaryMessage, data)
}

func (cnt *Client) sendConnect() error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	msg := &common.ClientConnect{
		Hostname: hostname,
		UUID:     cnt.uuid,
	}

	return cnt.SendMessage(msg)
}

// readEventsFromServer and pass the hid events out to the application via the InputEvents chanel
func (cnt *Client) readEventsFromServer() {
	for {
		_, data, err := cnt.ws.ReadMessage()
		if err != nil {
			if errors.Is(err, websocket.ErrCloseSent) {
				return
			}
			continue
		}

		cnt.Events <- data
	}
}

// autoClose the websocket/client when context ends
func (cnt *Client) handleIncommingEvents() {
	for {
		select {
		case <-cnt.ctx.Done():
			cnt.Close()
			return

		case ev := <-cnt.Input:
			data, err := ev.Marshal()
			if err != nil {
				cnt.ws.WriteMessage(websocket.BinaryMessage, data)
			}
		}
	}
}
