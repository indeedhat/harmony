package screens

import (
	"github.com/indeedhat/harmony/internal/common"
)

type virtualScreenSpace struct {
	Peers []Peer
}

// Width gives the full width of the virtual screen space
func (vss *virtualScreenSpace) Width() int {
	var (
		max = 1 << 16
		min = -max
	)

	for _, peer := range vss.Peers {
		for _, display := range peer.Displays {
			absPos := peer.AbsolutePosition(display)
			min = common.Min(min, absPos.X)
			max = common.Max(max, absPos.X+display.Width)
		}
	}

	return common.Abs(min - max)
}

// Height gives the full height of the virtual screen space
func (vss *virtualScreenSpace) Height() int {
	var (
		max = 1 << 16
		min = -max
	)

	for _, peer := range vss.Peers {
		for _, display := range peer.Displays {
			absPos := peer.AbsolutePosition(display)
			min = common.Min(min, absPos.Y)
			max = common.Max(max, absPos.Y+display.Height)
		}
	}

	return common.Abs(min - max)
}

// GetNewPeerPosition will calculate the default position of a new peer added to the screen space
// it will alway be placed to the right of the top left most screen in the virtual space
func (vss *virtualScreenSpace) GetNewPeerPosition() common.Vector2 {
	var pos common.Vector2

	for _, peer := range vss.Peers {
		for _, display := range peer.Displays {
			absPos := peer.AbsolutePosition(display)
			if pos.X < absPos.X+display.Width {
				pos = absPos.Add(common.Vector2{X: display.Width})
			} else if pos.X == absPos.X+display.Width && pos.Y > absPos.Y {
				pos = absPos.Add(common.Vector2{X: display.Width})
			}
		}
	}

	return pos
}
