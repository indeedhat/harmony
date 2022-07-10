package net

import (
	"syscall"

	"github.com/vmihailenco/msgpack/v5"
)

type WsMessage interface {
	Marshal() ([]byte, error)
}

type MsgType byte

const (
	MsgTypeCConnect MsgType = iota
	MsgTypeCReleaseControl
	MsgTypeSHidEvent
)

var _ WsMessage = (*ClientConnect)(nil)

// ClientConnect is sent from the client on connect to inform the server about itself
type ClientConnect struct {
	Hostname string `msgpack:"h"`
}

// Marshal ClientConnect struct into a byte array for sending via websocket
func (cc *ClientConnect) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(cc)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeCConnect), ';'}
	return append(base, data...), nil
}

var _ WsMessage = (*ClientReleaseControl)(nil)

// ClientReleaseControl informs the server to take control back over the hid devices
type ClientReleaseControl struct {
	X uint `msgpack:"x"`
	Y uint `msgpack:"y"`
}

// Marshal ClientReleaseControl struct into a byte array for sending via websocket
func (crc *ClientReleaseControl) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(crc)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeCReleaseControl), ';'}
	return append(base, data...), nil
}

var _ WsMessage = (*ServerHidEvent)(nil)

// ServerSendAction sends a HID event to be processed by the client
// this is just an alias of the evdev inputEvent struct
type ServerHidEvent struct {
	Time  syscall.Timeval `msgpack:"u"`
	Type  uint16          `msgpack:"t"`
	Code  uint16          `msgpack:"c"`
	Value int32           `msgpack:"v"`
}

// Marshal ServerHidEvent struct into a byte array for sending via websocket
func (shi *ServerHidEvent) Marshal() ([]byte, error) {
	data, err := msgpack.Marshal(shi)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(MsgTypeSHidEvent), ';'}
	return append(base, data...), nil
}
