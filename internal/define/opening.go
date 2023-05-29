package define

import (
	"sync/atomic"
)

type Run interface {
	Run() error
}

type Opening struct {
	running int32
}

// if not a pointer would not change running flag
func (d *Opening) Opened() bool {
	if atomic.LoadInt32(&d.running) == 1 {
		return true
	}
	atomic.CompareAndSwapInt32(&d.running, 0, 1)
	return false
}

// if not a pointer would not change running flag
func (d *Opening) Running() bool {
	return d.running == 1
}
