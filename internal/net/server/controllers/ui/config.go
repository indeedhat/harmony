package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/device"
)

// Index controller
func (ui *UI) Index() gin.HandlerFunc {
	type DisplayGroup struct {
		Width    int
		Height   int
		Displays []device.DisplayBounds
	}

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	return func(ctx *gin.Context) {
		groups := make([]DisplayGroup, 2)
		for _, display := range ui.displays {
			screen := device.DisplayBounds{
				X:      display.X / 4,
				Y:      display.Y / 4,
				Width:  display.Width / 4,
				Height: display.Height / 4,
			}

			groups[0].Displays = append(groups[0].Displays, screen)
			groups[0].Width = max(groups[0].Width, screen.X+screen.Width)
			groups[0].Height = max(groups[0].Height, screen.Y+screen.Height)

			groups[1].Displays = append(groups[1].Displays, screen)
			groups[1].Width = max(groups[1].Width, screen.X+screen.Width)
			groups[1].Height = max(groups[1].Height, screen.Y+screen.Height)
		}
		ctx.HTML(http.StatusOK, "index", gin.H{
			"groups": groups,
		})
	}
}
