package app

import (
	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/device"
	"github.com/indeedhat/harmony/internal/net"
	"github.com/indeedhat/harmony/internal/net/discovery"
	"github.com/indeedhat/harmony/internal/net/server/router"
)

type Harmony struct {
	ctx        *common.Context
	discover   *discovery.Service
	dev        *device.DeviceManager
	serverMode bool
	active     bool
	client     *net.Client
	uuid       uuid.UUID
}

// New sets up a new Harmony instance
func New(ctx *common.Context) (*Harmony, error) {
	dev, err := device.NewDeviceManager(ctx)
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
	}, nil
}

// Run the application
func (app *Harmony) Run() error {
	app.discover.Run()
	defer app.ctx.Cancel()

	for {
		select {
		case ev := <-app.client.Events:
			// pass incomming events from peers to the vdev
			app.dev.Input <- ev

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
		r := router.New(app.ctx, nil)
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
