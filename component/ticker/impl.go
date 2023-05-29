// Author: huaxr
// Time: 2022/7/6 10:08 上午
// Git: huaxr

package ticker

import (
	"context"
	"time"
)

type tick interface {
	TickHeartbeat() *T
}

type T struct {
	f func()
	n string
	t *time.Ticker
}

func NewT(f func(), name string, t *time.Ticker) *T {
	return &T{
		f: f,
		n: name,
		t: t,
	}
}

func RegisterTick(ti tick) {
	t := ti.TickHeartbeat()
	manager.Register(newJob(context.Background(), t.n, t.t, t.f))
}
