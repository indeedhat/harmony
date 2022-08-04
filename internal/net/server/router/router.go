package router

import (
	"encoding/json"
	"html/template"
	"mime"

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
	mime.AddExtensionType(".js", "application/javascript")
	router := gin.Default()

	screenManager := screens.NewScreenManager()

	_ = ui.New(router, screenManager)
	_ = socket.New(ctx, serverUUID, router, screenManager)

	viewsConfig := goview.DefaultConfig
	viewsConfig.Root = "./web/views"
	viewsConfig.DisableCache = true
	viewsConfig.Funcs = template.FuncMap{
		"json": func(data interface{}) string {
			bytes, err := json.Marshal(data)
			if err != nil {
				return ""
			}

			return string(bytes)
		},
	}
	router.HTMLRender = ginview.New(viewsConfig)

	router.Static("/js", "./web/public/js")

	return router
}
