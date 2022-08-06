package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/screens"
)

// Index controller
func (ui *UI) Index() gin.HandlerFunc {
	type DisplayGroup struct {
		Width       int
		Height      int
		UUID        uuid.UUID
		Displays    []screens.DisplayBounds
		Transitions []screens.TransitionZone
	}

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	return func(ctx *gin.Context) {
		// TODO: make this actually work from peer display config
		groups := make([]DisplayGroup, len(ui.screenManager.Peers))
		zones := ui.screenManager.CalculateTransitionZones()

		for i, peer := range ui.screenManager.Peers {
			group := DisplayGroup{
				UUID: peer.UUID,
			}

			for _, display := range peer.Displays {
				screen := screens.DisplayBounds{
					Position: common.Vector2{
						X: display.Position.X / 4,
						Y: display.Position.Y / 4,
					},
					Width:  display.Width / 4,
					Height: display.Height / 4,
				}

				group.Width = max(group.Width, screen.Position.X+screen.Width)
				group.Height = max(group.Height, screen.Position.Y+screen.Height)

				group.Displays = append(group.Displays, screen)
			}

			for _, zone := range zones[peer.UUID] {
				group.Transitions = append(group.Transitions, screens.TransitionZone{
					UUID:      zone.UUID,
					Direction: zone.Direction,
					Bounds: [2]common.Vector2{
						{
							X: zone.Bounds[0].X / 4,
							Y: zone.Bounds[0].Y / 4,
						},
						{
							X: zone.Bounds[1].X / 4,
							Y: zone.Bounds[1].Y / 4,
						},
					},
				})
			}

			groups[i] = group
		}

		ctx.HTML(http.StatusOK, "index", gin.H{
			"groups": groups,
		})
	}
}
