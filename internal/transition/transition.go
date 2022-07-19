package transition

type Direction uint8

const (
	Down Direction = iota
	Left
	Right
	Up
)

type TransitionZone struct {
	// X will always be the left most pixel of the zone
	X int
	// Y will always be the bottom most pixel of the zone
	Y         int
	Width     int
	height    int
	Direction Direction
}

// ShouldTransition calculates if the peer should transition foucs based on the defined zone
func (tz *TransitionZone) ShouldTransition(x, y, deltaX, deltaY int) bool {
	if x < tz.X || x > tz.X+tz.Width {
		return false
	}

	if y < tz.Y || y > tz.Y+tz.Width {
		return false
	}

	switch tz.Direction {
	case Down:
		return deltaY < 0
	case Left:
		return deltaX < 0
	case Right:
		return deltaX > 0
	case Up:
		return deltaY > 0
	default:
		return false
	}
}
