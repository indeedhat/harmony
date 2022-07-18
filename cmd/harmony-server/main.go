package main

import (
	"log"
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/device"
	"github.com/indeedhat/harmony/internal/net"
	"github.com/indeedhat/harmony/internal/net/router"
	"github.com/jezek/xgb/xproto"
)

func main() {
	ctx := common.NewContext()
	defer ctx.Cancel()

	dm, err := device.NewDeviceManager(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go eventConsumer(ctx)
	go watchCursor(ctx)
	go inputLock(ctx, dev)

	r := router.New(ctx, nil)
	r.Run()
}

func watchCursor(ctx *common.Context) {
	con, err := common.InitXCon()
	if err != nil {
		log.Fatal("connection to x failed")
		return
	}

	setup := xproto.Setup(con)
	window := setup.DefaultScreen(con).Root

	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
			if len(ctx.ClientPool) == 0 || ctx.ActiveClient != nil {
				continue
			}

			resp, err := xproto.QueryPointer(con, window).Reply()
			if err != nil {
				log.Print("qp error: " + err.Error())
				continue
			}

			// TODO: not hard code this
			if resp.RootY == 1439 && (resp.RootX >= 0 && resp.RootX <= 2559) {
				ctx.GrabQ <- 0
			}
		case <-ctx.Done():
			return
		}
	}
}

func inputLock(ctx *common.Context, dm *device.DeviceManager) {
	for {
		select {
		case idx := <-ctx.ReleaseQ:
			log.Print("release q")
			if ctx.ActiveClient != nil && idx == ctx.ActiveClient.Idx {
				dm.Release()
			}

		case idx := <-ctx.GrabQ:
			log.Print("grab q")
			if ctx.ActiveClient != nil || len(ctx.ClientPool) <= idx {
				log.Print("nope")
				continue
			}

			ctx.ActiveClient = ctx.ClientPool[idx]
			log.Print(dm.Grab())
		}
	}
}
