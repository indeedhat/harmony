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
	"github.com/indeedhat/harmony/internal/screens"
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
		go soc.readFromSocket(NewConn(ctx, ws), done)
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
		ws.Input <- data
	}
}

// readFromSocket and process/forward the messages
func (soc *Socket) readFromSocket(con *ConnectionWrapper, done chan struct{}) {
	defer close(done)

	con.Soc.SetReadLimit(config.MaxMessageSize)
	con.Soc.SetReadDeadline(time.Now().Add(config.PongWait))
	con.Soc.SetPongHandler(func(string) error {
		con.Soc.SetReadDeadline(time.Now().Add(config.PongWait))
		return nil
	})

	var conUUID *uuid.UUID
	defer soc.handleDisconnect(conUUID)

	for {
		messageType, data, err := con.Soc.ReadMessage()
		if err != nil {
			break
		}

		if messageType != websocket.BinaryMessage ||
			(conUUID == nil && data[0] != byte(events.MsgTypeConnect)) {

			continue
		}

		// handle events
		switch events.MsgType(data[0]) {
		case events.MsgTypeConnect:
			conUUID = soc.handleConnect(con, data)

		case events.MsgTypeInputEvent:
			soc.handleInputEvent(data)

		case events.MsgTypeChangeFoucs:
			soc.handleChangeFocus(conUUID, data)

		case events.MsgTypeReleaseFouces:
			soc.handleReleaseFocus()

		default:
			Logf("server", "unknown message type: %s", data[0])
		}
	}
}

// handleDisconnect cleans up the peer data on connection close
func (soc *Socket) handleDisconnect(conUUID *uuid.UUID) {
	if conUUID == nil {
		return
	}

	delete(soc.clients, *conUUID)

	if soc.activeClient == conUUID {
		soc.activeClient = nil
		soc.broadcast(&events.ReleaseFocus{})
	}

	zones := soc.screenManager.RemovePeer(*conUUID)
	soc.distributeTransitionZones(zones)
}

// handleReleaseFocus broadcasts the event out to all clients on force release recieved
func (soc *Socket) handleReleaseFocus() {
	Log("server", "release focus")
	soc.activeClient = nil
	soc.broadcast(&events.ReleaseFocus{})
}

// handleChangeFocus lets the appropriate client know it has focus
func (soc *Socket) handleChangeFocus(conUUID *uuid.UUID, data []byte) {
	Log("server", "change focus")
	var msg events.ChangeFocus
	if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
		log.Print("ws: failed to unmarshal message")
	}

	if _, ok := soc.clients[msg.UUID]; !ok {
		return
	}

	soc.activeClient = &msg.UUID

	recMessage := events.FocusRecieved{}
	data, err := recMessage.Marshal()
	if err == nil {
		soc.clients[msg.UUID].Input <- data
	}
}

// handleInputEvent forwards thi hid event to the appropriate peer
func (soc *Socket) handleInputEvent(data []byte) {
	if soc.activeClient == nil {
		Log("server", "no active client")
		return
	}

	var msg events.InputEvent
	if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
		log.Print("ws: failed to unmarshal message")
		return
	}

	client, ok := soc.clients[*soc.activeClient]
	if !ok {
		Log("server", "bad active client")
		return
	}

	client.Input <- data
}

// handleConnect handles a connection event and rebuilds the transition zones
func (soc *Socket) handleConnect(con *ConnectionWrapper, data []byte) *uuid.UUID {
	Log("server", "client connect")
	var msg events.ClientConnect

	if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
		return nil
	}

	soc.clients[msg.UUID] = con

	zones := soc.screenManager.AddPeer(msg.UUID, msg.Displays)
	soc.distributeTransitionZones(zones)

	return &msg.UUID
}

// distributeTransitionZones to the appropriate peers
func (soc *Socket) distributeTransitionZones(zones map[uuid.UUID][]screens.TransitionZone) {
	Log("server", "distribute tzones")
	// send updated transition zones
	for id, zones := range zones {
		con, ok := soc.clients[id]
		if !ok {
			continue
		}

		tzMessage := events.TransitionZoneAssigned(zones)
		data, err := tzMessage.Marshal()
		if err == nil {
			Log("server", "sending t zones")
			con.Input <- data
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
