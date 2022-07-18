package device

// DisplayBounds provides a common return type for displays
// regardless of platform
type DisplayBounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Cursor represents the current coordinates of the cursor
type Cursor struct {
	X int
	Y int
}

// Vdu provides a common interface for interacting with system dependent display managers
type Vdu interface {
	CursorPos() (*Cursor, error)
	DisplayBounds() ([]DisplayBounds, error)
	Close() error
}
