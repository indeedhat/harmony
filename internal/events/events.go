package events

import "github.com/vmihailenco/msgpack/v5"

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

func marshalEvent(ev any, typ MsgType) ([]byte, error) {
	data, err := msgpack.Marshal(ev)
	if err != nil {
		return nil, err
	}

	base := []byte{byte(typ), ';'}
	return append(base, data...), nil
}

// Unmarshal an event from its byte array
func Unmarshal[T any](data []byte) *T {
	var event T

	if err := msgpack.Unmarshal(data, &event); err != nil {
		return nil
	}

	return &event
}
