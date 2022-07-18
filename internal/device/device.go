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

type DeviceManager struct {
	Events chan *common.InputEvent
	Input  chan *common.InputEvent

	activeDevice *Device
	devices      []Device
	virtualDev   Device
	mux          sync.Mutex
	ctx          *common.Context
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

	go dm.consumeEvents()

	return dm, nil
}

// GrabAccess exclusive access to all the devices being watched
// this will stop rative input events being handled by any other program/service on the machine
func (dm *DeviceManager) GrabAccess() error {
	var err error

	for _, dev := range dm.devices {
		err = common.WrapError(err, dev.Grab())
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

		dm.devices = append(dm.devices[:i], dm.devices[i+1:]...)
	}
}

// Select a watched device to be active
func (dm *DeviceManager) Select(id string) bool {
	dm.mux.Lock()
	defer dm.mux.Unlock()

	dm.ClearSelection()

	for i, dev := range dm.devices {
		if dev.ID() != id {
			continue
		}

		dm.activeDevice = &dm.devices[i]
		return true
	}

	return false
}

// ClearSelection of active device
func (dm *DeviceManager) ClearSelection() {
	dm.activeDevice = nil
}

func (dm *DeviceManager) trackEvents(dev Device) {
	for {
		event, err := dev.Read()
		if err != nil {
			return
		}

		dm.Events <- event
	}

	dm.Forget(dev)
}

func (dm *DeviceManager) consumeEvents() {
	for {
		select {
		case ev := <-dm.Events:
			if dm.activeDevice == nil {
				continue
			}

			(*dm.activeDevice).Write(ev)

		case <-dm.ctx.Done():
			return
		}
	}
}
