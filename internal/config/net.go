package config

import "time"

// Web Server
const (
	ServerPort uint16 = 4283
)

// Peer Discovery
const (
	MulticastAddress   = "239.2.3.239:2399"
	DiscoveryPollCount = 3
	DiscoveryInterval  = 2
)

// Socket server
const (
	WriteWait        = 10 * time.Second
	MaxMessageSize   = 8192
	PongWait         = 60 * time.Second
	PingPeriod       = (PongWait * 9) / 10
	CloseGracePeriod = 10 * time.Second
)
