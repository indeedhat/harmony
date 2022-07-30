package events

import "syscall"

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
