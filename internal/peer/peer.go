package peer

import (
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/device"
	"github.com/indeedhat/harmony/internal/net/discovery"
)

type Peer struct {
	ctx      *common.Context
	discover *discovery.Service
	dev      *device.DeviceManager
}

// New sets up a new peer instance
func New(ctx *common.Context) (*Peer, error) {
	dev, err := device.NewDeviceManager(ctx)
	if err != nil {
		return nil, err
	}

	discover, err := discovery.New(ctx)
	if err != nil {
		return nil, err
	}

	return &Peer{
		ctx:      ctx,
		discover: discover,
		dev:      dev,
	}, nil
}

// Run the peer
func (pr *Peer) Run() error {
	pr.discover.Run()

	for {
		select {
		case server := <-pr.discover.Server:
			pr.handleDiscoveryMessage(server)

		case ev := <-pr.dev.Events:
			pr.handleInputEvent(ev)

		case <-pr.ctx.Done():
			return nil
		}
	}
}

func (pr *Peer) handleDiscoveryMessage(server discovery.Server) {
	if server.IpAddress == "" {
		// TODO: spin up the server
	} else {
		// TODO: connect to the server
	}
}

func (pr *Peer) handleInputEvent(event *common.InputEvent) {
	// TODO: pass event to the active server client
}
