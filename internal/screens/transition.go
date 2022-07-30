package screens

import (
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
)

type Direction uint8

const (
	Down Direction = iota
	Left
	Right
	Up
)

type TransitionZone struct {
	// This wiss be passed between shared between the peers on both side of the TransitionZone
	UUID   uuid.UUID
	Bounds [2]common.Vector2
	// Direction of travel required to trigger the transition
	Direction Direction
}

// ShouldTransition calculates if the peer should transition foucs based on the defined zone
func (zone *TransitionZone) ShouldTransition(x, y, deltaX, deltaY int) bool {
	if x < zone.Bounds[0].X || x > zone.Bounds[1].X {
		return false
	}

	if y < zone.Bounds[0].Y || y > zone.Bounds[1].Y {
		return false
	}

	switch zone.Direction {
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
