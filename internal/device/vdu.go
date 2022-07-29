package device

import "github.com/indeedhat/harmony/internal/common"

// Vdu provides a common interface for interacting with system dependent display managers
type Vdu interface {
	Close() error
	CursorPos() (*common.Vector2, error)
	DisplayBounds() ([]common.DisplayBounds, error)
	HideCursor() error
	ShowCursor() error
}
