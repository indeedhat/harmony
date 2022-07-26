package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/harmony/internal/common"
)

type Socket struct {
	appCtx       *common.Context
	clients      map[uuid.UUID]*websocket.Conn
	activeClient *uuid.UUID
	serverUUID   uuid.UUID
}

// New UI controller
func New(ctx *common.Context, serverUUID uuid.UUID, router *gin.Engine) *Socket {
	socket := &Socket{
		appCtx:     ctx,
		clients:    make(map[uuid.UUID]*websocket.Conn),
		serverUUID: serverUUID,
	}

	socket.routes(router)

	return socket
}

func (soc *Socket) routes(router *gin.Engine) {
	router.GET("/ws", soc.Ws())
}
