package events

import (
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/screens"
)

// ClientConnect is sent from the client on connect to inform the server about itself
type ClientConnect struct {
	Hostname string    `msgpack:"h"`
	UUID     uuid.UUID `msgpack:"u"`
	Displays []screens.DisplayBounds
}

// Marshal ClientConnect struct into a byte array for sending via websocket
func (ev *ClientConnect) Marshal() ([]byte, error) {
	return marshalEvent(ev, MsgTypeConnect)
}

// String gives the string name of the event type
func (ev *ClientConnect) String() string {
	return "ClientConnect"
}

var _ WsMessage = (*ClientConnect)(nil)
