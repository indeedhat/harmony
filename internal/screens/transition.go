package screens

import (
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
)

type Direction uint8

const (
	Up Direction = iota
	Right
	Down
	Left
)

type TransitionZone struct {
	// This wiss be passed between shared between the peers on both side of the TransitionZone
	UUID   uuid.UUID
	Bounds [2]common.Vector2
	// Direction of travel required to trigger the transition
	Direction Direction
}

// ShouldTransition calculates if the peer should transition foucs based on the defined zone
func (zone *TransitionZone) ShouldTransition(current, previous common.Vector2) bool {
	if current.X < zone.Bounds[0].X || current.X > zone.Bounds[1].X {
		return false
	}

	if current.Y < zone.Bounds[0].Y || current.Y > zone.Bounds[1].Y {
		return false
	}

	delta := current.Sub(previous)

	switch zone.Direction {
	case Down:
		return delta.Y < 0
	case Left:
		return delta.X < 0
	case Right:
		return delta.X > 0
	case Up:
		return delta.Y > 0
	default:
		return false
	}
}
