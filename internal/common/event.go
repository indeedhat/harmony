package common

import (
	"syscall"

	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/transition"
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
	Marshal() ([]byte, error)
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

var _ WsMessage = (*InputEvent)(nil)

// Marshal ServerHidEvent struct into a byte array for sending via websocket
func (ie *InputEvent) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(ie)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeInputEvent), ';'}
	return append(base, data...), nil
}

// ClientConnect is sent from the client on connect to inform the server about itself
type ClientConnect struct {
	Hostname string    `msgpack:"h"`
	UUID     uuid.UUID `msgpack:"u"`
}

var _ WsMessage = (*ClientConnect)(nil)

// Marshal ClientConnect struct into a byte array for sending via websocket
func (cc *ClientConnect) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(cc)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeConnect), ';'}
	return append(base, data...), nil
}

// ChangeFocus from the active client to a peer
type ChangeFocus struct {
	UUID uuid.UUID `msgpack:"u"`
	X    uint      `msgpack:"x"`
	Y    uint      `msgpack:"y"`
}

var _ WsMessage = (*ChangeFocus)(nil)

// Marshal ChangeFocus struct into a byte array for sending via websocket
func (cf *ChangeFocus) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(cf)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeChangeFoucs), ';'}
	return append(base, data...), nil
}

// FocusRecieved from a peer
// this message will be sent to the active client to inform them they now have focus
type FocusRecieved struct {
	// ID of the transition zone that triggerd the focus
	ID  uuid.UUID
	Pos uint `msgpack:"x"`
}

var _ WsMessage = (*FocusRecieved)(nil)

// Marshal FocusRecieved struct into a byte array for sending via websocket
func (fr *FocusRecieved) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(fr)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeFocusRecieved), ';'}
	return append(base, data...), nil
}

// ReleaseFocus from all peers
// When this message is sent from any peer all peers will have their focus removed making all hid
// devices operatie for their local client again
type ReleaseFocus struct {
}

var _ WsMessage = (*ReleaseFocus)(nil)

// Marshal FocusRecieved struct into a byte array for sending via websocket
func (rf *ReleaseFocus) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(rf)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeFocusRecieved), ';'}
	return append(base, data...), nil
}

// TransitionZoneAssigned will be sent to clients on connect and whenever
// the global screen arrangement is updated, it is used to pass the new details
// of their transition zones
type TransitionZoneAssigned []transition.TransitionZone

var _ WsMessage = (*TransitionZoneAssigned)(nil)

// Marshal FocusRecieved struct into a byte array for sending via websocket
func (tza TransitionZoneAssigned) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(tza)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeTrasitionAssigned), ';'}
	return append(base, data...), nil
}
