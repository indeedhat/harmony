package ui

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/screens"
)

type UI struct {
	screenManager *screens.ScreenManager
}

// New UI controller
func New(router *gin.Engine, screenManager *screens.ScreenManager) *UI {
	ui := &UI{
		screenManager: screenManager,
	}

	ui.routes(router)

	return ui
}

func (ui *UI) routes(router *gin.Engine) {
	router.GET("", ui.Index())
}
