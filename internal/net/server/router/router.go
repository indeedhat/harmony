package router

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/device"
	"github.com/indeedhat/harmony/internal/net/controllers/socket"
	"github.com/indeedhat/harmony/internal/net/controllers/ui"
)

// New UI controller group
func New(ctx *common.Context, displays []device.DisplayBounds) *gin.Engine {
	router := gin.Default()

	_ = ui.New(router, displays)
	_ = socket.New(ctx, router)

	viewsConfig := goview.DefaultConfig
	viewsConfig.Root = "./web/views"
	viewsConfig.DisableCache = true
	router.HTMLRender = ginview.New(viewsConfig)

	return router
}
