package ui

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/device"
)

type UI struct {
	displays []device.DisplayBounds
}

// New UI controller
func New(router *gin.Engine, displays []device.DisplayBounds) *UI {
	ui := &UI{
		displays,
	}

	ui.routes(router)

	return ui
}

func (ui *UI) routes(router *gin.Engine) {
	router.GET("", ui.Index())
}
