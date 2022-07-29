package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/common"
)

// Index controller
func (ui *UI) Index() gin.HandlerFunc {
	type DisplayGroup struct {
		Width    int
		Height   int
		Displays []common.DisplayBounds
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
			screen := common.DisplayBounds{
				Position: common.Vector2{
					X: display.Position.X / 4,
					Y: display.Position.Y / 4,
				},
				Width:  display.Width / 4,
				Height: display.Height / 4,
			}

			groups[0].Displays = append(groups[0].Displays, screen)
			groups[0].Width = max(groups[0].Width, screen.Position.X+screen.Width)
			groups[0].Height = max(groups[0].Height, screen.Position.Y+screen.Height)

			groups[1].Displays = append(groups[1].Displays, screen)
			groups[1].Width = max(groups[1].Width, screen.Position.X+screen.Width)
			groups[1].Height = max(groups[1].Height, screen.Position.Y+screen.Height)
		}

		ctx.HTML(http.StatusOK, "index", gin.H{
			"groups": groups,
		})
	}
}
