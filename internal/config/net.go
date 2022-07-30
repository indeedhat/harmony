package config

import "time"

// Socket server
const (
	MaxMessageSize = 8192
	PongWait       = 60 * time.Second
	PingPeriod     = (PongWait * 9) / 10
)
