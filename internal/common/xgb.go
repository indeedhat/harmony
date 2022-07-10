package common

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
)

// InitXCon isitialises the connection to xserver
func InitXCon() (*xgb.Conn, error) {
	con, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}

	if err := xinerama.Init(con); err != nil {
		return nil, err
	}

	return con, nil
}
