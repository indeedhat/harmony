package device

import (
	"errors"
	"fmt"
	"sync"

	"github.com/indeedhat/harmony/internal/common"
)

type Device interface {
	// Read an input event from the device
	Read() (*common.InputEvent, error)
	// Write an input event to the device
	Write(*common.InputEvent) error
	// Grab exclisive accell to device
	Grab() error
	// Release exclusive access on the device
	Release() error
	// Close the device handle
	Close() error
	// Conform to fmt.Stringer
	String() string
	// ID of the device, where this comes from will be system independent but should always be uniqe
	ID() string
}

type DevicePlus interface {
	Device

	// MoveCursor relative to its current position
	MoveCursor(x, y int)
}

type DeviceManager struct {
	// Events stream from grabbed devices to be consumed externally
	Events chan *common.InputEvent
	// Input events from external server to be passed to the vdev
	Input chan *common.InputEvent

	// grabbed state of watched devices
	grabbed bool
	// devices currently being watched
	devices []Device
	// virtual device used for incomming events from peers
	virtualDev DevicePlus
	mux        sync.Mutex
	ctx        *common.Context
}

// NewDeviceManager constructor
func NewDeviceManager(ctx *common.Context) (*DeviceManager, error) {
	devices := FindObservableDevices()
	if len(devices) == 0 {
		return nil, errors.New("no observable devices found")
	}

	vdev, err := CreateVirtualDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual device: %w", err)
	}

	dm := &DeviceManager{
		Events: make(chan *common.InputEvent),
		Input:  make(chan *common.InputEvent),

		ctx:        ctx,
		devices:    devices,
		virtualDev: vdev,
	}

	for _, dev := range dm.devices {
		go dm.trackEvents(dev)
	}

	return dm, nil
}

// GrabAccess exclusive access to all the devices being watched
// this will stop rative input events being handled by any other program/service on the machine
func (dm *DeviceManager) GrabAccess() error {
	var err error
	dm.grabbed = true

	for _, dev := range dm.devices {
		err = common.WrapError(err, dev.Grab())
	}

	if err != nil {
		dm.ReleaseAccess()
	}

	return err
}

// ReleaseAccess exclusive access from all the devices being watched
// this will return input devices to their natural state where native input events can be
// handled by any other program/service on the machine
func (dm *DeviceManager) ReleaseAccess() error {
	var err error

	for _, dev := range dm.devices {
		err = common.WrapError(err, dev.Release())
	}

	dm.grabbed = false

	return err
}

// Close all the wacthed devices
func (dm *DeviceManager) Close() error {
	var err error

	for _, dev := range dm.devices {
		err = common.WrapError(err, dev.Release())
		err = common.WrapError(err, dev.Close())
	}

	err = common.WrapError(err, dm.virtualDev.Release())
	err = common.WrapError(err, dm.virtualDev.Close())

	return err
}

// Watch an aditional device
func (dm *DeviceManager) Watch(newDev Device) {
	dm.mux.Lock()
	defer dm.mux.Unlock()

	newId := newDev.ID()
	for _, dev := range dm.devices {
		if dev.ID() != newId {
			return
		}
	}

	dm.devices = append(dm.devices, newDev)
	if dm.grabbed {
		newDev.Grab()
	}
}

// Forget a watched device
func (dm *DeviceManager) Forget(newDev Device) {
	dm.mux.Lock()
	defer dm.mux.Unlock()

	newId := newDev.ID()
	for i, dev := range dm.devices {
		if dev.ID() != newId {
			continue
		}

		if dm.grabbed {
			dev.Release()
		}

		dm.devices = append(dm.devices[:i], dm.devices[i+1:]...)
	}
}

// MoveCursor relative to its current position
func (dm *DeviceManager) MoveCursor(x, y int) {
	dm.virtualDev.MoveCursor(x, y)
}

func (dm *DeviceManager) trackEvents(dev Device) {
	defer dm.Forget(dev)

	for {
		event, err := dev.Read()
		if err != nil {
			return
		}

		dm.Events <- event
	}
}

func (dm *DeviceManager) consumeIncommingEvents() {
	for {
		select {
		case <-dm.ctx.Done():
			return

		case ev := <-dm.Input:
			dm.virtualDev.Write(ev)
		}
	}
}
