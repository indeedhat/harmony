package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/indeedhat/harmony/internal/screens"
)

// Index controller
func (ui *UI) Index() gin.HandlerFunc {
	type DisplayGroup struct {
		Width       int
		Height      int
		UUID        uuid.UUID
		Hostname    string
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
				UUID:     peer.UUID,
				Hostname: peer.Hostname,
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
				zone.Bounds.X = zone.Bounds.X / config.UIScaleFactor
				zone.Bounds.Y = zone.Bounds.X / config.UIScaleFactor
				zone.Bounds.W = zone.Bounds.W / config.UIScaleFactor
				zone.Bounds.Z = zone.Bounds.Z / config.UIScaleFactor

				zone.Target.Bounds.X = zone.Target.Bounds.X / config.UIScaleFactor
				zone.Target.Bounds.Y = zone.Target.Bounds.X / config.UIScaleFactor
				zone.Target.Bounds.W = zone.Target.Bounds.W / config.UIScaleFactor
				zone.Target.Bounds.Z = zone.Target.Bounds.Z / config.UIScaleFactor

				group.Transitions = append(group.Transitions, zone)
			}

			groups[i] = group
		}

		ctx.HTML(http.StatusOK, "index", gin.H{
			"groups": groups,
		})
	}
}
