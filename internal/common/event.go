package common

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
