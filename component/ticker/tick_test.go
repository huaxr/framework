// Author: huaxr
// Time: 2022/7/5 11:02 上午
// Git: huaxr

package ticker

import (
	"context"
	"testing"
	"time"
)

type testTick struct{}

func (t testTick) TickHeartbeat() *T {
	return NewT(nil, "test", time.NewTicker(1*time.Second))
}

func TestRegister(t *testing.T) {
	RegisterTick(testTick{})
}

func TestTick(t *testing.T) {
	ctx := context.Background()
	tick1 := time.NewTicker(1 * time.Second)
	job1 := func() {
		t.Log("hello")
	}
	a := newJob(ctx, "tick1", tick1, job1, time.Now().Add(10*time.Second))

	manager.Register(a)

	tick2 := time.NewTicker(2 * time.Second)
	job2 := func() {
		t.Log("world")
	}
	b := newJob(ctx, "tick2", tick2, job2)
	manager.Register(b)

	time.Sleep(3 * time.Second)
	manager.Revoke("tick2")

	time.Sleep(2 * time.Second)
	tick3 := time.NewTicker(2 * time.Second)
	job3 := func() {
		t.Log("tick3")
	}
	c := newJob(ctx, "tick3", tick3, job3)
	manager.Register(c)
	select {}
}
