package common

import (
	"syscall"

	"github.com/vmihailenco/msgpack/v5"
)

type MsgType byte

const (
	MsgTypeCConnect MsgType = iota
	MsgTypeCReleaseControl
	MsgTypeSHidEvent
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

	base := []byte{byte(MsgTypeSHidEvent), ';'}
	return append(base, data...), nil
}

// ClientConnect is sent from the client on connect to inform the server about itself
type ClientConnect struct {
	Hostname string `msgpack:"h"`
}

var _ WsMessage = (*ClientConnect)(nil)

// Marshal ClientConnect struct into a byte array for sending via websocket
func (cc *ClientConnect) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(cc)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeCConnect), ';'}
	return append(base, data...), nil
}

// ClientReleaseControl informs the server to take control back over the hid devices
type ClientReleaseControl struct {
	X uint `msgpack:"x"`
	Y uint `msgpack:"y"`
}

var _ WsMessage = (*ClientReleaseControl)(nil)

// Marshal ClientReleaseControl struct into a byte array for sending via websocket
func (crc *ClientReleaseControl) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(crc)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeCReleaseControl), ';'}
	return append(base, data...), nil
}
