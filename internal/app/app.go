package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/indeedhat/harmony/internal/device"
	"github.com/indeedhat/harmony/internal/events"
	. "github.com/indeedhat/harmony/internal/logger"
	"github.com/indeedhat/harmony/internal/net"
	"github.com/indeedhat/harmony/internal/net/discovery"
	"github.com/indeedhat/harmony/internal/net/server/router"
	"github.com/indeedhat/harmony/internal/screens"
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
	tZones []screens.TransitionZone
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
	Logf("app", "uuid: %s", app.uuid)
	app.discover.Run()

	defer app.ctx.Cancel()
	defer app.discover.Close()

	// need to block until we have a client connected
	server := <-app.discover.Server
	Log("app", "handling discovery event")
	if err := app.handleDiscoveryMessage(server); err != nil {
		return err
	}

	go app.watchTransitionZones()

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
	switch events.MsgType(data[0]) {
	case events.MsgTypeInputEvent:
		if event := events.Unmarshal[events.InputEvent](data[2:]); event != nil {
			app.dev.Input <- event
		}

	case events.MsgTypeReleaseFouces:
		Log("app", "handling release focus")

		if event := events.Unmarshal[events.ReleaseFocus](data[2:]); event != nil {
			Log("app", "release focus")
			app.dev.ReleaseAccess()
			app.active = false

			cursorPos, err := app.vdu.CursorPos()
			if err != nil {
				return
			}

			displays, err := app.vdu.DisplayBounds()
			if err != nil || len(displays) == 0 {
				return
			}

			desiredPos := common.Vector2{
				X: displays[0].Position.X + displays[0].Width/2,
				Y: displays[0].Position.Y + displays[0].Height/2,
			}

			diff := desiredPos.Sub(*cursorPos)
			app.dev.MoveCursor(diff)
		}

	case events.MsgTypeFocusRecieved:
		Log("app", "handling focus recieved")
		if event := events.Unmarshal[events.FocusRecieved](data[2:]); event == nil {
			return
		}

		Log("app", "focus recieved")
		app.active = false
		app.dev.ReleaseAccess()

		// TODO: move mouse to proper place in transition zone

	case events.MsgTypeTrasitionAssigned:
		Log("app", "handling new transition zones")
		if event := events.Unmarshal[events.TransitionZoneAssigned](data[2:]); event != nil {
			Log("app", "recieved new transition zones")
			Logf("app", "%#v", *event)
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

func (app *Harmony) handleInputEvent(event *events.InputEvent) {
	app.handleEmergancyRelease(event)
	if !app.active {
		return
	}

	Log("app", "sending event to client")
	app.client.SendMessage(event)
}

func (app *Harmony) handleEmergancyRelease(event *events.InputEvent) {
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
		app.client.SendMessage(&events.ReleaseFocus{})
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

	screens, err := app.vdu.DisplayBounds()
	if err != nil {
		return err
	}

	client, err := net.NewClient(app.ctx, app.uuid, ip, screens)
	if err != nil {
		return err
	}

	app.client = client
	return nil
}

func (app *Harmony) watchTransitionZones() {
	var (
		lastPos *common.Vector2
		ticker  = time.NewTicker(time.Millisecond * config.TransitionPollIntervalMs)
	)
	defer ticker.Stop()

	for {
		select {
		case <-app.ctx.Done():
			return

		case <-ticker.C:
			if len(app.tZones) == 0 {
				continue
			}

			pos, err := app.vdu.CursorPos()
			if err != nil {
				continue
			}

			if lastPos == nil {
				lastPos = pos
				continue
			}

			for _, zone := range app.tZones {
				if !zone.ShouldTransition(*pos, *lastPos) {
					continue
				}

				if err := app.dev.GrabAccess(); err != nil {
					continue
				}

				Log("app", "giving up focus")
				app.active = true
				app.client.SendMessage(&events.ChangeFocus{
					UUID: zone.UUID,
					Pos:  *pos,
				})
			}

			lastPos = pos
		}
	}
}
