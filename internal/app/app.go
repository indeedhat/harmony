package app

import (
	"github.com/google/uuid"
	"github.com/holoplot/go-evdev"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/device"
	"github.com/indeedhat/harmony/internal/net"
	"github.com/indeedhat/harmony/internal/net/discovery"
	"github.com/indeedhat/harmony/internal/net/server/router"
	"github.com/indeedhat/harmony/internal/transition"
	"github.com/vmihailenco/msgpack/v5"
)

type Harmony struct {
	ctx        *common.Context
	discover   *discovery.Service
	dev        *device.DeviceManager
	vdu        device.Vdu
	serverMode bool
	active     bool
	client     *net.Client
	uuid       uuid.UUID
	tZones     []transition.TransitionZone
}

// New sets up a new Harmony instance
func New(ctx *common.Context) (*Harmony, error) {
	dev, err := device.NewDeviceManager(ctx)
	if err != nil {
		return nil, err
	}

	vdu, err := device.NewVdu()
	if err != nil {
		return nil, err
	}

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

	for {
		select {
		case ev := <-app.client.Events:
			app.handelServerEvent(ev)

		case server := <-app.discover.Server:
			if err := app.handleDiscoveryMessage(server); err != nil {
				return err
			}

		case ev := <-app.dev.Events:
			app.handleInputEvent(ev)

		case <-app.ctx.Done():
			return nil
		}
	}
}

func (app *Harmony) handelServerEvent(data []byte) {
	switch common.MsgType(data[0]) {
	case common.MsgTypeInputEvent:
		if event := unmarshalEvent[common.InputEvent](data[2:]); event != nil {
			app.dev.Input <- event
		}

	case common.MsgTypeReleaseFouces:
		if event := unmarshalEvent[common.ReleaseFocus](data[2:]); event != nil {
			app.dev.ReleaseAccess()
			app.active = false
			// TODO: move the cursor to the middle of the main monitor
		}

	case common.MsgTypeFocusRecieved:
		if event := unmarshalEvent[common.FocusRecieved](data[2:]); event == nil {
			return
		}

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
		if event := unmarshalEvent[common.TransitionZoneAssigned](data[2:]); event != nil {
			app.tZones = *event
		}
	}
}

func (app *Harmony) handleDiscoveryMessage(server discovery.Server) error {
	if server.IpAddress == "" {
		app.runServer()
	}

	return app.startClient(server.IpAddress)
}

func (app *Harmony) handleInputEvent(event *common.InputEvent) {
	if !app.active {
		return
	}

	app.client.Input <- event
}

func (app *Harmony) runServer() {
	if app.serverMode {
		return
	}

	app.serverMode = true
	go func() {
		r := router.New(app.ctx, app.uuid, nil)
		r.Run()
	}()
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
