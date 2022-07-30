package events

import (
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
)

// ChangeFocus from the active client to a peer
type ChangeFocus struct {
	UUID uuid.UUID `msgpack:"u"`
	Pos  common.Vector2   `msgpack:"p"`
}

// Marshal ChangeFocus struct into a byte array for sending via websocket
func (ev *ChangeFocus) Marshal() ([]byte, error) {
	return marshalEvent(ev, MsgTypeChangeFoucs)
}

// String gives the string name of the event type
func (ev *ChangeFocus) String() string {
	return "ChangeFocus"
}

var _ WsMessage = (*ChangeFocus)(nil)

// FocusRecieved from a peer
// this message will be sent to the active client to inform them they now have focus
type FocusRecieved struct {
	// ID of the transition zone that triggerd the focus
	ID  uuid.UUID
	Pos uint `msgpack:"x"`
}

// Marshal FocusRecieved struct into a byte array for sending via websocket
func (ev *FocusRecieved) Marshal() ([]byte, error) {
	return marshalEvent(ev, MsgTypeFocusRecieved)
}

// String gives the string name of the event type
func (ev *FocusRecieved) String() string {
	return "FocusRecieved"
}

var _ WsMessage = (*FocusRecieved)(nil)

// ReleaseFocus from all peers
// When this message is sent from any peer all peers will have their focus removed making all hid
// devices operatie for their local client again
type ReleaseFocus struct {
}

// Marshal FocusRecieved struct into a byte array for sending via websocket
func (ev *ReleaseFocus) Marshal() ([]byte, error) {
	return marshalEvent(ev, MsgTypeReleaseFouces)
}

// String gives the string name of the event type
func (ev *ReleaseFocus) String() string {
	return "ReleaseFocus"
}

var _ WsMessage = (*ReleaseFocus)(nil)
