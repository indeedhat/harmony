package discovery

import (
	"fmt"
	"net"
	"time"

	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/config"
	. "github.com/indeedhat/harmony/internal/logger"
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
	ClusterId  string  `msgpack:"i"`
	StartTime  int64   `msgpack:"s"`
}

type Server struct {
	IpAddress string
	Port      uint16
}

type Service struct {
	// When servers are discoverd their details will be sent over this chanel
	Server chan Server

	startTime int64
	clusterId string
	state     peerState
	ctx       *common.Context
	addr      *net.UDPAddr
	con       *net.UDPConn
}

// New creates a new instance of the Sirvice struct
// Service is used to handle udp multicast for peer discovery and master negotiations
func New(ctx *common.Context) (*Service, error) {
	Log("discovery", "resolving address")
	addr, err := net.ResolveUDPAddr("udp4", ctx.Config.Discovery.MulticastAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve udp address")
	}

	Log("discovery", "starting listener")
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to cennect to multicast address: %w", err)
	}

	return &Service{
		Server:    make(chan Server),
		clusterId: ctx.Config.Discovery.ClusterId,
		state:     stateDiscovery,
		ctx:       ctx,
		addr:      addr,
		con:       conn,
	}, nil
}

// Run the discovery service
func (svc *Service) Run() {
	svc.startTime = time.Now().UnixMilli()
	svc.state = stateDiscovery

	go svc.discover()
	go svc.listen()
}

// Close the discovery service
func (svc *Service) Close() {
	svc.con.Close()
}

// discover peers running in server mode
func (svc *Service) discover() {
	var (
		pollCount int
		ticker    = time.NewTicker(time.Duration(svc.ctx.Config.Discovery.PollIntervalSeconds) * time.Second)
	)
	defer ticker.Stop()

	for {
		select {
		case <-svc.ctx.Done():
			Log("discovery", "shutting down")
			return

		case <-ticker.C:
			if svc.state != stateDiscovery {
				return
			}

			if pollCount >= svc.ctx.Config.Discovery.PollCaunt {
				Log("discovery", "poll limit reached, requesting server start")
				svc.state = stateServer
				svc.Server <- Server{}
				return
			}

			pollCount++
			data, err := svc.discoveryMsg()

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
		buf := make([]byte, 44)
		_, addr, err := svc.con.ReadFromUDP(buf)
		if err != nil {
			Log("discovery", "read failed")
			continue
		}

		var msg message
		if err := msgpack.Unmarshal(buf, &msg); err != nil {
			Logf("discovery", "unmarshal failed: %s", err)
			continue
		}

		Logf("discovery", "handling message: %d", msg.Type)
		switch msg.Type {
		case serverResponse:
			if svc.state != stateDiscovery {
				continue
			}

			Log("discovery", "server found")
			svc.state = statePeer
			svc.Server <- Server{
				IpAddress: addr.IP.String(),
				Port:      msg.ServerPort,
			}
			return

		case queryPeers:
			if svc.state != stateServer {
				Log("discovery", "not in server mode")
				continue
			}

			if msg.ApiVersion != config.ApiVersion || msg.ClusterId != svc.clusterId {
				Logf("discovery", "api(%d, %d) id(%s, %s)",
					msg.ApiVersion,
					config.ApiVersion,
					msg.ClusterId,
					svc.clusterId,
				)
				continue
			}

			data, err := svc.serverMsg()
			if err != nil {
				Logf("discovery", "msg error: %s", err)
				continue
			}

			Log("discovery", "responding to peer discovery request")
			svc.con.WriteToUDP(data, svc.addr)
		}
	}
}

func (svc *Service) serverMsg() ([]byte, error) {
	return msgpack.Marshal(&message{
		StartTime:  svc.startTime,
		Type:       serverResponse,
		ApiVersion: config.ApiVersion,
		ServerPort: uint16(svc.ctx.Config.Server.Port),
		ClusterId:  svc.ctx.Config.Discovery.ClusterId,
	})
}

func (svc *Service) discoveryMsg() ([]byte, error) {
	return msgpack.Marshal(&message{
		StartTime:  svc.startTime,
		Type:       queryPeers,
		ApiVersion: config.ApiVersion,
		ClusterId:  svc.ctx.Config.Discovery.ClusterId,
		ServerPort: 0,
	})
}
