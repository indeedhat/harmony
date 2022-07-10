package main

import (
	"log"
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/net"
	"github.com/indeedhat/harmony/internal/net/router"
	"github.com/jezek/xgb/xproto"
)

func main() {
	ctx := common.NewContext()
	defer ctx.Cancel()

	devPath := "/dev/input/event25"
	dev, err := evdev.Open(devPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	go eventConsumer(ctx)
	go eventListener(ctx, dev)
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

// eventListener listenes for events from a single device and forwards them to the global queue
func eventListener(ctx *common.Context, dev *evdev.InputDevice) {
	for {
		event, err := dev.ReadOne()
		if err != nil {
			return
		}

		if ctx.ActiveClient == nil {
			continue
		}

		// TODO:
		// listen for alt enter
		// send alt up event to client
		// release devices
		// disable active client
		// move cursor back to center

		// if

		ctx.EventQ <- event
	}
}

// eventConsumer consumes events from the global queue and forwards them to the active client
// if any
func eventConsumer(ctx *common.Context) {
	for ev := range ctx.EventQ {
		if ctx.ActiveClient == nil {
			continue
		}

		ctx.ActiveClient.SendEvent(&net.ServerHidEvent{
			Time:  ev.Time,
			Type:  uint16(ev.Type),
			Code:  uint16(ev.Code),
			Value: ev.Value,
		})
	}
}

func inputLock(ctx *common.Context, dev *evdev.InputDevice) {
	for {
		select {
		case idx := <-ctx.ReleaseQ:
			log.Print("release q")
			if ctx.ActiveClient != nil && idx == ctx.ActiveClient.Idx {
				dev.Ungrab()
			}

		case idx := <-ctx.GrabQ:
			log.Print("grab q")
			if ctx.ActiveClient != nil || len(ctx.ClientPool) <= idx {
				continue
			}

			ctx.ActiveClient = ctx.ClientPool[idx]
			dev.Grab()
		}
	}
}
