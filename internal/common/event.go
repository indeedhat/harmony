package common

import (
	"syscall"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

type MsgType byte

const (
	MsgTypeConnect MsgType = iota
	MsgTypeFocusRecieved
	MsgTypeChangeFoucs
	MsgTypeReleaseFouces
	MsgTypeInputEvent
	MsgTypeTrasitionAssigned
)

// WsMessage interface describes any message/event that is transmissable
// via the websocket connection
type WsMessage interface {
	// Marshal an event into a byte array using mesgpack
	Marshal() ([]byte, error)
	// String representation of the message type
	String() string
}

// InputEvent provides a system independent format for transfering
// input events between components
// it is currently modeled on... cloned from the linux evdev/uinput event but
// will likely have to change once more os's start to be added
type InputEvent struct {
	Time  syscall.Timeval `msgpack:"u"`
	Type  uint16          `msgpack:"t"`
	Code  uint16          `msgpack:"c"`
	Value int32           `msgpack:"v"`
}

// Marshal ServerHidEvent struct into a byte array for sending via websocket
func (ev *InputEvent) Marshal() ([]byte, error) {
	return marshalEvent(ev, MsgTypeInputEvent)
}

// String gives the string name of the event type
func (ev *InputEvent) String() string {
	return "InputEvent"
}

var _ WsMessage = (*InputEvent)(nil)

// ClientConnect is sent from the client on connect to inform the server about itself
type ClientConnect struct {
	Hostname string    `msgpack:"h"`
	UUID     uuid.UUID `msgpack:"u"`
	Displays []DisplayBounds
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

// ChangeFocus from the active client to a peer
type ChangeFocus struct {
	UUID uuid.UUID `msgpack:"u"`
	Pos  Vector2   `msgpack:"p"`
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

// TransitionZoneAssigned will be sent to clients on connect and whenever
// the global screen arrangement is updated, it is used to pass the new details
// of their transition zones
type TransitionZoneAssigned []TransitionZone

// Marshal FocusRecieved struct into a byte array for sending via websocket
func (ev TransitionZoneAssigned) Marshal() ([]byte, error) {
	return marshalEvent(ev, MsgTypeTrasitionAssigned)
}

// String gives the string name of the event type
func (tza *TransitionZoneAssigned) String() string {
	return "TransitionZoneAssigned"
}

var _ WsMessage = (*TransitionZoneAssigned)(nil)

func marshalEvent(ev any, typ MsgType) ([]byte, error) {
	data, err := msgpack.Marshal(ev)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(typ), ';'}
	return append(base, data...), nil
}
