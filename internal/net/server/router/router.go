package router

import (
	"mime"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/net/server/socket"
	"github.com/indeedhat/harmony/internal/net/server/ui"
	"github.com/indeedhat/harmony/internal/screens"
)

// New UI controller group
func New(ctx *common.Context, serverUUID uuid.UUID, displays []screens.DisplayBounds) *gin.Engine {
	mime.AddExtensionType(".js", "application/javascript")
	router := gin.Default()

	screenManager := screens.NewScreenManager()

	_ = ui.New(router, screenManager)
	_ = socket.New(ctx, serverUUID, router, screenManager)

	viewsConfig := goview.DefaultConfig
	viewsConfig.Root = "./web/views"
	viewsConfig.DisableCache = true
	router.HTMLRender = ginview.New(viewsConfig)

	router.Static("/js", "./web/public/js")

	return router
}
