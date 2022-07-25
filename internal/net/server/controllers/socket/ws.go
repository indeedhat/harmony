package socket

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/vmihailenco/msgpack/v5"
)

// Ws handler for server - client communications
func (soc *Socket) Ws() gin.HandlerFunc {
	var upgrader = websocket.Upgrader{}

	return func(ctx *gin.Context) {
		ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		done := make(chan struct{})
		go soc.readFromSocket(ws, done)
		go ping(soc.appCtx, ws)

		<-done
	}
}

// readFromSocket and process/forward the messages
func (soc *Socket) readFromSocket(ws *websocket.Conn, done chan struct{}) {
	defer close(done)

	ws.SetReadLimit(config.MaxMessageSize)
	ws.SetReadDeadline(time.Now().Add(config.PongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(config.PongWait))
		return nil
	})

	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			break
		}

		if messageType != websocket.BinaryMessage {
			continue
		}

		switch data[0] {
		case byte(common.MsgTypeConnect):
			var msg common.ClientConnect

			if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
				log.Print("ws: failed to unmarshal message")
			}

			soc.clients[msg.UUID] = ws

		case byte(common.MsgTypeChangeFoucs):
			// TODO: change the focus to the new active client
			// send message to active client

		case byte(common.MsgTypeFocusRecieved):
			// TODO: unlock inputs
		}
	}
}

// ping the client to keep the connection alive
func ping(ctx *common.Context, ws *websocket.Conn) {
	ticker := time.NewTicker(config.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().
				Add(config.WriteWait))

			if err != nil && errors.Is(err, websocket.ErrCloseSent) {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
