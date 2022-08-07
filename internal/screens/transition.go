package screens

import (
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
)

type TransitionTarget struct {
	// This wiss be passed between shared between the peers on both side of the TransitionZone
	UUID uuid.UUID
	// Bounds of the transition zone on the target machine
	Bounds common.Vector4
}

type TransitionZone struct {
	Target TransitionTarget
	// Bounds of the TransitionZone
	Bounds common.Vector4
	// Direction of travel required to trigger the transition
	Direction common.Direction
}

// ShouldTransition calculates if the peer should transition foucs based on the defined zone
func (zone *TransitionZone) ShouldTransition(current, previous common.Vector2) bool {
	if current.X < zone.Bounds.X || current.X > zone.Bounds.W {
		return false
	}

	if current.Y < zone.Bounds.Y || current.Y > zone.Bounds.Z {
		return false
	}

	delta := current.Sub(previous)

	switch zone.Direction {
	case common.DirectionDown:
		return delta.Y < 0
	case common.DirectionLeft:
		return delta.X < 0
	case common.DirectionRight:
		return delta.X > 0
	case common.DirectionUp:
		return delta.Y > 0
	default:
		return false
	}
}
