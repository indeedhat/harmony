package common

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/holoplot/go-evdev"
)

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type Context struct {
	// Cilent related things
	poolMu       sync.Mutex
	ClientPool   []*Client
	ActiveClient *Client

	// Event related things
	EventQ chan *evdev.InputEvent

	// Clients that want to release control send on this chanel
	ReleaseQ chan int
	// Clients that want to grab input send on this chanel
	GrabQ chan int

	// Close related things
	done     atomic.Value
	doneMu   sync.Mutex
	canceled bool

	// device related things
	devices []*evdev.InputDevice
}

// NewContext constructor
func NewContext() *Context {
	return &Context{
		EventQ:     make(chan *evdev.InputEvent),
		ReleaseQ:   make(chan int),
		GrabQ:      make(chan int),
		ClientPool: make([]*Client, 0),
	}
}

// Done returns a chan if the context is canceled
func (ctx *Context) Done() <-chan struct{} {
	done := ctx.done.Load()
	if done != nil {
		return done.(chan struct{})
	}

	ctx.doneMu.Lock()
	defer ctx.doneMu.Unlock()

	done = ctx.done.Load()
	if done == nil {
		done = make(chan struct{})
		ctx.done.Store(done)
	}

	return done.(chan struct{})
}

// Cancel the context
func (ctx *Context) Cancel() {
	ctx.doneMu.Lock()
	defer ctx.doneMu.Unlock()

	if ctx.canceled {
		return // already canceled
	}

	ctx.canceled = true

	done, _ := ctx.done.Load().(chan struct{})
	if done == nil {
		ctx.done.Store(closedchan)
	} else {
		close(done)
	}
}

// GrabDevices will stop input events for the local devices from propergating to the system
func (ctx *Context) GrabDevices() {
	for _, dev := range ctx.devices {
		dev.Grab()
	}
}

// ReleaseDevices allowing input events to be processed locally
func (ctx *Context) ReleaseDevices() {
	for _, dev := range ctx.devices {
		dev.Ungrab()
	}
}

// AddClient to the pool
func (ctx *Context) AddClient(client *Client) {
	ctx.poolMu.Lock()
	defer ctx.poolMu.Unlock()

	ctx.ClientPool = append(ctx.ClientPool, client)
}

// RemoveClient from the pool
func (ctx *Context) RemoveClient(client *Client) {
	ctx.poolMu.Lock()
	defer ctx.poolMu.Unlock()

	log.Print("sending release")
	ctx.ReleaseQ <- client.Idx
	time.Sleep(10 * time.Millisecond)

	if ctx.ActiveClient != nil && ctx.ActiveClient.Idx == client.Idx {

		ctx.ActiveClient = nil
	}

	ctx.ClientPool = append(ctx.ClientPool[:client.Idx], ctx.ClientPool[client.Idx:]...)
}
