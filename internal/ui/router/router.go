package router

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/screen"
	"github.com/indeedhat/harmony/internal/ui/controllers"
)

func New(displays []screen.DisplayBounds) *gin.Engine {
	router := gin.Default()

	_ = controllers.New(router, displays)

	viewsConfig := goview.DefaultConfig
	viewsConfig.Root = "./web/views"
	viewsConfig.DisableCache = true
	router.HTMLRender = ginview.New(viewsConfig)

	return router
}
