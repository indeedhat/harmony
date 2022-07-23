package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/harmony/internal/common"
)

type Socket struct {
	appCtx *common.Context
}

// New UI controller
func New(ctx *common.Context, router *gin.Engine) *Socket {
	socket := &Socket{
		appCtx: ctx,
	}

	socket.routes(router)

	return socket
}

func (soc *Socket) routes(router *gin.Engine) {
	router.GET("/ws", soc.Ws())
}
