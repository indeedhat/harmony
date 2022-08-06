package socket

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/screens"
)

type ConnectionWrapper struct {
	ctx   context.Context
	Soc   *websocket.Conn
	Input chan []byte
}

func NewConn(ctx context.Context, ws *websocket.Conn) *ConnectionWrapper {
	con := &ConnectionWrapper{
		ctx:   ctx,
		Soc:   ws,
		Input: make(chan []byte),
	}

	go con.consumeIncommingMessages()

	return con
}

func (con *ConnectionWrapper) Close() error {
	return con.Soc.Close()
}

func (con *ConnectionWrapper) consumeIncommingMessages() {
	for {
		select {
		case <-con.ctx.Done():
			return
		case data := <-con.Input:
			con.Soc.WriteMessage(websocket.BinaryMessage, data)

		}
	}
}

type Socket struct {
	appCtx        *common.Context
	clients       map[uuid.UUID]*ConnectionWrapper
	activeClient  *uuid.UUID
	serverUUID    uuid.UUID
	screenManager *screens.ScreenManager
}

// New UI controller
func New(ctx *common.Context, serverUUID uuid.UUID, router *gin.Engine, screenManager *screens.ScreenManager) *Socket {
	socket := &Socket{
		appCtx:        ctx,
		clients:       make(map[uuid.UUID]*ConnectionWrapper),
		serverUUID:    serverUUID,
		screenManager: screenManager,
	}

	socket.routes(router)

	return socket
}

func (soc *Socket) routes(router *gin.Engine) {
	router.GET("/ws", soc.Ws())
}
