package device

import (
	"github.com/holoplot/go-evdev"
	"github.com/indeedhat/harmony/internal/common"
)

// IsAltUpEvent checks if the given input event is a keyup for the alt key
//
// This might seem like an odly specific floating function but its required for
// emergancy release
func IsAltUpEvent(ev *common.InputEvent) bool {
	return ev.Value == 1 && // key up
		ev.Type == evdev.EV_KEY &&
		(ev.Code == evdev.KEY_LEFTALT || ev.Code == evdev.KEY_RIGHTALT)
}
