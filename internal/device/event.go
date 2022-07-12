package device

import "syscall"

type InputEvent struct {
	Time  syscall.Timeval `msgpack:"u"`
	Type  uint16          `msgpack:"t"`
	Code  uint16          `msgpack:"c"`
	Value int32           `msgpack:"v"`
}
