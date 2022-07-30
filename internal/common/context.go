package common

import (
	"sync"
	"sync/atomic"

	"github.com/indeedhat/harmony/internal/config"
)

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

// Context provides a common context for the application
// currently it just reimplements the close functionalyty of the
// stdlib context package
type Context struct {
	// Close related things
	done     atomic.Value
	doneMu   sync.Mutex
	canceled bool
	Config   *config.Config
}

// NewContext constructor
func NewContext(conf *config.Config) *Context {
	return &Context{
		Config: conf,
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
