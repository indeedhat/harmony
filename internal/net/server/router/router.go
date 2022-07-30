package router

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/net/server/controllers/socket"
	"github.com/indeedhat/harmony/internal/net/server/controllers/ui"
	"github.com/indeedhat/harmony/internal/screens"
)

// New UI controller group
func New(ctx *common.Context, serverUUID uuid.UUID, displays []screens.DisplayBounds) *gin.Engine {
	router := gin.Default()

	_ = ui.New(router, displays)
	_ = socket.New(ctx, serverUUID, router)

	viewsConfig := goview.DefaultConfig
	viewsConfig.Root = "./web/views"
	viewsConfig.DisableCache = true
	router.HTMLRender = ginview.New(viewsConfig)

	return router
}
