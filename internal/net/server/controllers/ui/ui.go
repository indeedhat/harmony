package ui

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/common"
)

type UI struct {
	displays []common.DisplayBounds
}

// New UI controller
func New(router *gin.Engine, displays []common.DisplayBounds) *UI {
	ui := &UI{
		displays,
	}

	ui.routes(router)

	return ui
}

func (ui *UI) routes(router *gin.Engine) {
	router.GET("", ui.Index())
}
