package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/indeedhat/harmony/internal/device"
	. "github.com/indeedhat/harmony/internal/logger"
	"github.com/indeedhat/harmony/internal/net"
	"github.com/indeedhat/harmony/internal/net/discovery"
	"github.com/indeedhat/harmony/internal/net/server/router"
	"github.com/vmihailenco/msgpack/v5"
)

type Harmony struct {
	ctx *common.Context
	// discovery service used to locate existing servers on startup
	discover *discovery.Service
	// os independent device manager
	dev *device.DeviceManager
	// local window server manager
	vdu device.Vdu
	// if this peer is in server mode or not
	serverMode bool
	// if this peer has active control of another peer or not
	active bool
	// client connected to the socket server
	client *net.Client
	// uuid to identify this peer over the network
	uuid uuid.UUID
	// transition zones are used to define screen edges that 'transition' to other peers
	tZones []common.TransitionZone
	// cache the times of the last n (config.AltEscapeCount) alt key up events
	// if enough events happen in a specified time frame then all clients peers
	// will be told to release focus and exclusive access locks on all devices
	altCache []time.Time
}

// New sets up a new Harmony instance
func New(ctx *common.Context) (*Harmony, error) {
	Log("app", "hid discovery")
	dev, err := device.NewDeviceManager(ctx)
	if err != nil {
		return nil, err
	}

	Log("app", "vdu discovery")
	vdu, err := device.NewVdu()
	if err != nil {
		return nil, err
	}

	Log("app", "starting peer discovery")
	discover, err := discovery.New(ctx)
	if err != nil {
		return nil, err
	}

	return &Harmony{
		ctx:      ctx,
		discover: discover,
		dev:      dev,
		uuid:     uuid.New(),
		vdu:      vdu,
	}, nil
}

// Run the application
func (app *Harmony) Run() error {
	app.discover.Run()

	defer app.ctx.Cancel()
	defer app.discover.Close()

	// need to block until we have a client connected
	server := <-app.discover.Server
	Log("app", "handling discovery event")
	if err := app.handleDiscoveryMessage(server); err != nil {
		return err
	}

	for {
		select {
		case ev := <-app.dev.Events:
			app.handleInputEvent(ev)

		case ev := <-app.client.Events:
			app.handleServerEvent(ev)

		case <-app.ctx.Done():
			return nil
		}
	}
}

func (app *Harmony) handleServerEvent(data []byte) {
	switch common.MsgType(data[0]) {
	case common.MsgTypeInputEvent:
		if event := unmarshalEvent[common.InputEvent](data[2:]); event != nil {
			app.dev.Input <- event
		}

	case common.MsgTypeReleaseFouces:
		Log("app", "handling release focus")

		if event := unmarshalEvent[common.ReleaseFocus](data[2:]); event != nil {
			Log("app", "release focus")
			app.dev.ReleaseAccess()
			app.active = false
			// TODO: move the cursor to the middle of the main monitor
		}

	case common.MsgTypeFocusRecieved:
		Log("app", "handling focus recieved")
		if event := unmarshalEvent[common.FocusRecieved](data[2:]); event == nil {
			return
		}

		Log("app", "focus recieved")
		app.active = true
		// TODO: calculate the desired
		var x, y int
		cursor, err := app.vdu.CursorPos()
		if err != nil {
			// no point in moving cursor if we dont know its current position
			return
		}

		app.dev.MoveCursor(x-cursor.X, y-cursor.Y)

	case common.MsgTypeTrasitionAssigned:
		Log("app", "handling new transition zones")
		if event := unmarshalEvent[common.TransitionZoneAssigned](data[2:]); event != nil {
			Log("app", "recieved new transition zones")
			app.tZones = *event
		}

	default:
		Logf("app", "unknown message type: %d", data[0])
	}
}

func (app *Harmony) handleDiscoveryMessage(server discovery.Server) error {
	if server.IpAddress == "" {
		Log("app", "starting server")
		app.runServer()
	}

	Log("app", "connecting as peer")
	return app.startClient(server.IpAddress)
}

func (app *Harmony) handleInputEvent(event *common.InputEvent) {
	app.handleEmergancyRelease(event)
	if !app.active {
		return
	}

	app.client.Input <- event
}

func (app *Harmony) handleEmergancyRelease(event *common.InputEvent) {
	if !device.IsAltUpEvent(event) {
		return
	}

	// cache last 3 timestamps
	if len(app.altCache) == config.AltEscapeCount {
		app.altCache = append(app.altCache[:config.AltEscapeCount-1], time.Now())
	} else {
		app.altCache = append(app.altCache, time.Now())
	}

	if len(app.altCache) != config.AltEscapeCount {
		return
	}

	diff := app.altCache[config.AltEscapeCount-1].Sub(app.altCache[0])
	if diff <= time.Second*config.AltEscapeTimeframe {
		Log("app", "emergancy release")
		app.altCache = []time.Time{}
		app.client.Input <- &common.ReleaseFocus{}
	}
}

func (app *Harmony) runServer() {
	if app.serverMode {
		return
	}

	app.serverMode = true
	go func() {
		r := router.New(app.ctx, app.uuid, nil)
		r.Run(fmt.Sprint(":", config.ServerPort))
	}()

	// bit dirty but need to make sure the server satrts before connecting as a peer
	Log("app", "waiting for server start")
	time.Sleep(time.Second * 2)
}

func (app *Harmony) startClient(ip string) error {
	if ip == "" {
		ip = "127.0.0.1"
	}

	client, err := net.NewClient(app.ctx, app.uuid, ip)
	if err != nil {
		return err
	}

	app.client = client
	return nil
}

func unmarshalEvent[T any](data []byte) *T {
	var event T

	if err := msgpack.Unmarshal(data, &event); err != nil {
		return nil
	}

	return &event
}
