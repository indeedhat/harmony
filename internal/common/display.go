package common

// DisplayBounds provides a common return type for displays
// regardless of platform
type DisplayBounds struct {
	Position Vector2
	Width    int
	Height   int
}
