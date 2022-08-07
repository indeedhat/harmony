package common

type Direction uint8

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
	DirectionNone
)

// Vector2 represents a fixed point in 2d space
type Vector2 struct {
	X int `msgpack:"x"`
	Y int `msgpack:"y"`
}

// Add adds the given vector from this one and returns a new one as a result
func (v2 Vector2) Add(point Vector2) Vector2 {
	return Vector2{
		X: v2.X + point.X,
		Y: v2.Y + point.Y,
	}
}

// Sub subtracts the given vector from this one and returns a new one as a result
func (v2 Vector2) Sub(point Vector2) Vector2 {
	return Vector2{
		X: v2.X - point.X,
		Y: v2.Y - point.Y,
	}
}

// Vector4 is being used to represent a rectangle
// top left is represented by X and Y
// bottom right is represented by W and Z in place of X2 and Y2 respectively
type Vector4 struct {
	W int
	X int
	Y int
	Z int
}

// Overlaps checks if another rectangle overlaps with this one
func (v4 Vector4) Overlaps(rect Vector4) bool {
	if v4.X >= rect.W || rect.X >= v4.W {
		return false
	}

	if v4.Y >= rect.Z || rect.Y >= v4.Z {
		return false
	}

	return true
}

// Touches returns the side that this rectangle touches the given one
// if the rectangles overlap then this will return no touch direction
func (v4 Vector4) Touches(rect Vector4) Direction {
	if v4.Overlaps(rect) {
		return DirectionNone
	}

	if v4.X == rect.W && v4.Y < rect.Z && v4.Z > rect.Y {
		return DirectionLeft
	}

	if v4.W == rect.X && v4.Y < rect.Z && v4.Z > rect.Y {
		return DirectionRight
	}

	if v4.Y == rect.Z && v4.X < rect.W && v4.W > rect.X {
		return DirectionUp
	}

	if v4.Z == rect.Y && v4.X < rect.W && v4.W > rect.X {
		return DirectionDown
	}

	return DirectionNone
}
