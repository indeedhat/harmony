package client

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/indeedhat/harmony/internal/net"
	"golang.org/x/net/context"
)

type Client struct {
	MessageQ chan<- net.WsMessage
	Idx      int
	Config   *net.ClientConnect
	Socket   *websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc
}

// New constructor
func New(ws *websocket.Conn, idx int) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		Socket:   ws,
		MessageQ: make(chan<- net.WsMessage),
		Idx:      idx,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// SendEvent to the client
func (c *Client) SendEvent(event *net.ServerHidEvent) {
	data, err := event.Marshal()
	if err != nil {
		log.Print("failed to marshal hid event")
	}

	c.Socket.WriteMessage(websocket.BinaryMessage, data)
}

// Close the client
func (c *Client) Close() {
	// close ws
	c.Socket.Close()
}

// Done aliases the client contexts Done method
func (c *Client) Done() <-chan struct{} {
	return c.ctx.Done()
}
