package socket

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/indeedhat/harmony/internal/events"
	. "github.com/indeedhat/harmony/internal/logger"
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

// broadcast a message to all active clients
func (soc *Socket) broadcast(msg events.WsMessage) {
	Logf("server", "broadcast: %s", msg)

	data, err := msg.Marshal()
	if err != nil {
		Logf("server", "broadcast marshal failure: %s", err)
		return
	}

	for _, ws := range soc.clients {
		ws.WriteMessage(websocket.BinaryMessage, data)
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

	var conUUID *uuid.UUID
	defer func() {
		if conUUID == nil {
			return
		}

		delete(soc.clients, *conUUID)
		if soc.activeClient == conUUID {
			soc.activeClient = nil
		}

		soc.broadcast(&events.ReleaseFocus{})
	}()

	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			break
		}

		if messageType != websocket.BinaryMessage {
			continue
		}

		// only allow connect messages until the uuid is set
		if conUUID == nil && data[0] != byte(events.MsgTypeConnect) {
			continue
		}

		// handle events
		switch events.MsgType(data[0]) {
		case events.MsgTypeConnect:
			Log("server", "client connect")
			var msg events.ClientConnect

			if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
				log.Print("ws: failed to unmarshal message")
				continue
			}

			soc.clients[msg.UUID] = ws
			conUUID = &msg.UUID

		case events.MsgTypeInputEvent:
			if soc.activeClient == nil {
				continue
			}

			var msg events.InputEvent
			if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
				log.Print("ws: failed to unmarshal message")
			}

		case events.MsgTypeChangeFoucs:
			Log("server", "change focus")
			if conUUID != soc.activeClient {
				continue
			}

			var msg events.ChangeFocus
			if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
				log.Print("ws: failed to unmarshal message")
			}

			if _, ok := soc.clients[msg.UUID]; !ok {
				continue
			}

			soc.activeClient = &msg.UUID

			recMessage := events.FocusRecieved{}
			data, err := recMessage.Marshal()
			if err == nil {
				soc.clients[msg.UUID].WriteMessage(websocket.BinaryMessage, data)
			}

		case events.MsgTypeReleaseFouces:
			Log("server", "release focus")
			soc.activeClient = nil
			soc.broadcast(&events.ReleaseFocus{})

		default:
			Logf("server", "unknown message type: %s", data[0])
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
