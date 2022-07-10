package main

import (
	"errors"
	"log"
	"strings"
	// "time"

	// "github.com/davecgh/go-spew/spew"
	evdev "github.com/gvalkov/golang-evdev"
	// "github.com/indeedhat/kvm/internal/common"
	// "github.com/indeedhat/kvm/internal/mouse"
	// "github.com/jezek/xgb/xproto"
)

// func main() {
// 	var (
// 		// g502
// 		devPath = "/dev/input/event25"
// 	// devPath = "/dev/input/event26"
// 	// hhkb
// 	// devPath = "/dev/input/event27"
// 	// devPath = "/dev/input/event28"
// 	// devPath = "/dev/input/event29"
// 	// devPath = "/dev/input/event30"
// 	// devPath = "/dev/input/event31"
// 	)

// 	// spew.Dump(filterDevices("HHKB"))
// 	// for _, dev := range filterDevices("G502") {
// 	// 	spew.Dump(dev)
// 	// 	spew.Dump(dev.Capabilities)
// 	// }
// 	// log.Fatal(watchDevice(devPath))
// 	// log.Print(mouse.GetCursorPos())
// 	log.Fatal(grabAtSide(devPath))
// }

func grabAtSide(path string) error {
	var (
		grabbed bool
		// leeway  int
	)
	device, err := evdev.Open(path)
	if err != nil {
		return errors.New("connect: " + err.Error())
	}

	defer func() {
		if grabbed {
			device.Release()
		}
	}()

	// con, err := common.InitXCon()
	// if err != nil {
	// 	return errors.New("xcon: " + err.Error())
	// }

	// setup := xproto.Setup(con)
	// window := setup.DefaultScreen(con).Root

	for {
		_, err := device.Read()
		if err != nil {
			return errors.New("read: " + err.Error())
		}
		continue

		// if grabbed {
		// 	log.Printf("%s: %d", ev.String(), ev.Code)
		// 	continue
		// }

		// if leeway > 0 {
		// 	// leeway--
		// 	continue
		// }

		// repl, err := xproto.QueryPointer(con, window).Reply()
		// if repl.RootX != 0 {
		// 	leeway = 10
		// }
		// if err := device.Grab(); err != nil {
		// 	return errors.New("grab: " + err.Error())
		// }

		// grabbed = true
		// go func() {
		// 	<-time.After(time.Second * 5)
		// 	device.Release()
		// 	leeway = 30
		// 	grabbed = false
		// }()
	}
}

func filterDevices(substr string) []*evdev.InputDevice {
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return nil
	}

	if substr == "" {
		return devices
	}

	var filtered []*evdev.InputDevice

	for _, dev := range devices {
		if strings.Contains(dev.Name, substr) {
			filtered = append(filtered, dev)
		}
	}

	return filtered
}

func watchDevice(path string) error {
	device, err := evdev.Open(path)
	if err != nil {
		return errors.New("connect: " + err.Error())
	}

	if err := device.Grab(); err != nil {
		return errors.New("grob: " + err.Error())
	}
	defer device.Release()

	for {
		ev, err := device.ReadOne()
		if err != nil {
			return errors.New("read: " + err.Error())
		}

		log.Printf("%s: %d", ev.String(), ev.Code)
	}
}

func findInupt0s() []*evdev.InputDevice {
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return nil
	}

	var filtered []*evdev.InputDevice

	for _, dev := range devices {
		if strings.Contains(dev.Phys, "/input0") {
			filtered = append(filtered, dev)
		}
	}

	return filtered
}
