package net

import (
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
	"golang.org/x/net/context"
)

type Client struct {
	Input chan events.WsMessage
	// Events coming from the server
	Events chan []byte

	ctx      context.Context
	ctxClose context.CancelFunc
	config   *config.Config
	ws       *websocket.Conn
	uuid     uuid.UUID
}

// NewClient harmony client
func NewClient(ctx *common.Context, uuid uuid.UUID, ip string, screens []screens.DisplayBounds) (*Client, error) {
	serverAddress := fmt.Sprintf("%s:%d", ip, ctx.Config.Server.Port)
	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/ws"}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	cntx, ctxClose := context.WithCancel(context.Background())
	client := &Client{
		uuid:     uuid,
		ctx:      cntx,
		ctxClose: ctxClose,
		config:   ctx.Config,
		ws:       ws,
		Events:   make(chan []byte),
		Input:    make(chan events.WsMessage),
	}

	go client.readEventsFromServer()
	go client.consumeIncommingMessages()

	client.sendConnect(screens)

	return client, nil
}

// Close the client
func (cnt *Client) Close() error {
	cnt.ctxClose()
	return cnt.ws.Close()
}

// Done exposes the Done method from the clients context
func (cnt *Client) Done() <-chan struct{} {
	return cnt.ctx.Done()
}

// sendConnect message to the server
func (cnt *Client) sendConnect(screens []screens.DisplayBounds) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	msg := &events.ClientConnect{
		Hostname: hostname,
		UUID:     cnt.uuid,
		Displays: screens,
	}

	Log("app", "sending connect")
	cnt.Input <- msg
}

// readEventsFromServer and pass the hid events out to the application via the InputEvents chanel
func (cnt *Client) readEventsFromServer() {
	for {
		_, data, err := cnt.ws.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				cnt.Close()
				return
			}
			Logf("client", "read error: %s", err)
			continue
		}

		cnt.Events <- data
	}
}

func (cnt *Client) consumeIncommingMessages() {
	for {
		select {
		case <-cnt.ctx.Done():
			Log("client", "done")
			return

		case msg := <-cnt.Input:
			data, err := msg.Marshal()
			if err != nil {
				Logf("client", "failed to marshal event: %s", err)
				continue

			}

			if err := cnt.ws.WriteMessage(websocket.BinaryMessage, data); err != nil {
				Logf("client", "ws write failed: %s", err)
			}
		}
	}
}
