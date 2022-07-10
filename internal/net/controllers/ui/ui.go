package ui

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/screen"
)

type UI struct {
	displays []screen.DisplayBounds
}

// New UI controller
func New(router *gin.Engine, displays []screen.DisplayBounds) *UI {
	ui := &UI{
		displays,
	}

	ui.routes(router)

	return ui
}

func (ui *UI) routes(router *gin.Engine) {
	router.GET("", ui.Index())
}
