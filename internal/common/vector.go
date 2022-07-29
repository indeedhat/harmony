package common

// Vector2 represents a fixed point in 2d space
type Vector2 struct {
	X int `msgpack:"x"`
	Y int `msgpack:"y"`
}

// Sub subtracts the given vector from this one and returns a new one as a result
func (v2 Vector2) Sub(point Vector2) Vector2 {
	return Vector2{
		X: v2.X - point.X,
		Y: v2.Y - point.Y,
	}
}
