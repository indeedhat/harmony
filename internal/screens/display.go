package screens

import "github.com/indeedhat/harmony/internal/common"

// DisplayBounds provides a common return type for displays
// regardless of platform
type DisplayBounds struct {
	Position common.Vector2
	Width    int
	Height   int
}
