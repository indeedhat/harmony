package discovery

import (
	"fmt"
	"net"
	"time"

	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	"github.com/vmihailenco/msgpack/v5"
)

type msgType uint8

const (
	queryPeers msgType = iota
	serverResponse
)

type peerState uint8

const (
	stateDiscovery peerState = iota
	statePeer
	stateServer
)

type message struct {
	Type       msgType `msgpack:"t"`
	ApiVersion uint8   `msgpack:"v"`
	ServerPort uint16  `msgpack:"p"`
}

type Server struct {
	IpAddress string
	Port      uint16
}

type Service struct {
	// When servers are discoverd their details will be sent over this chanel
	Server chan Server

	state peerState
	ctx   *common.Context
	addr  *net.UDPAddr
	con   *net.UDPConn
}

// New creates a new instance of the Sirvice struct
// Service is used to handle udp multicast for peer discovery and master negotiations
func New(ctx *common.Context) (*Service, error) {
	addr, err := net.ResolveUDPAddr("udp4", config.MulticastAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve udp address")
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to cennect to multicast address: %w", err)
	}

	return &Service{
		state: stateDiscovery,
		ctx:   ctx,
		addr:  addr,
		con:   conn,
	}, nil
}

// Run the discovery service
func (svc *Service) Run() {

}

// discover peers running in server mode
func (svc *Service) discover() {
	var (
		pollCount int
		ticker    = time.NewTicker(time.Duration(config.DiscoveryInterval) * time.Second)
	)

	for {
		select {
		case <-svc.ctx.Done():
			return

		case <-ticker.C:
			if svc.state != stateDiscovery {
				return
			}

			if pollCount >= config.DiscoveryPollCount {
				return
			}

			pollCount++
			data, err := discoveryMsg()

			if err != nil {
				continue
			}

			svc.con.WriteToUDP(data, svc.addr)
		}
	}
}

// listen for messages from peers and respond appropriately
func (svc *Service) listen() {
	for {
		buf := make([]byte, 14)
		_, addr, err := svc.con.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var msg message
		if err := msgpack.Unmarshal(buf, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case serverResponse:
			if svc.state != stateDiscovery {
				continue
			}

			svc.state = stateDiscovery
			svc.Server <- Server{
				IpAddress: addr.IP.String(),
				Port:      msg.ServerPort,
			}
			return

		case queryPeers:
			if svc.state != stateServer {
				continue
			}

			data, err := serverMsg()
			if err != nil {
				continue
			}

			svc.con.WriteToUDP(data, svc.addr)
		}
	}
}

func serverMsg() ([]byte, error) {
	return msgpack.Marshal(&message{
		Type:       serverResponse,
		ApiVersion: config.ApiVersion,
		ServerPort: config.ServerPort,
	})
}

func discoveryMsg() ([]byte, error) {
	return msgpack.Marshal(&message{
		Type:       queryPeers,
		ApiVersion: config.ApiVersion,
	})
}
