package socket

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/net"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	writeWait        = 10 * time.Second
	maxMessageSize   = 8192
	pongWait         = 60 * time.Second
	pingPeriod       = (pongWait * 9) / 10
	closeGracePeriod = 10 * time.Second
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

		clientIndex := len(soc.appCtx.ClientPool)
		client := common.NewClient(ws, clientIndex)
		soc.appCtx.AddClient(client)
		defer soc.appCtx.RemoveClient(client)

		done := make(chan struct{})
		go readFromSocket(client, ws, done)
		go ping(client, ws)

		<-done
	}
}

// readFromSocket and process/forward the messages
func readFromSocket(client *common.Client, ws *websocket.Conn, done chan struct{}) {
	defer close(done)

	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
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
		case byte(net.MsgTypeCConnect):
			var msg net.ClientConnect
			if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
				log.Print("ws: failed to unmarshal message")
			}

			client.Config = &msg

		case byte(net.MsgTypeCReleaseControl):
			// TODO: make this happen
		}
	}
}

// ping the client to keep the connection alive
func ping(client *common.Client, ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				if err.Error() == "websocket: close sent" {
					return
				}
			}

		case <-client.Done():
			return
		}
	}
}
