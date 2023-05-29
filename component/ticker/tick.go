// Author: huaxr
// Time: 2022/7/5 10:48 上午
// Git: huaxr

package ticker

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/huaxr/framework/internal/define"

	"github.com/huaxr/framework/logx"
)

var (
	manager *jobManager
	lock    sync.Mutex
	once    sync.Once
)

type id string

type job struct {
	ctx    context.Context
	id     id
	ticker *time.Ticker
	// functions that really need to be executed on time
	f func()
	// exits the current task signal
	stop chan struct{}
	// users can register for callbacks when triggered
	callback func()
	start    time.Time

	forever bool
	// deadline kill this job when current time > dead
	dead time.Time
}

type jobManager struct {
	lock  *sync.Mutex
	input chan *job
	count int
	// binding job with id, job stop is a signal notifier
	pool map[id]*job
}

func init() {
	once.Do(func() {
		manager = &jobManager{
			lock:  &sync.Mutex{},
			input: make(chan *job, 0),
			count: 0,
			pool:  make(map[id]*job),
		}
		go manager.Start()
	})
}

func (m *jobManager) Register(jon *job) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.pool[jon.id]; ok {
		logx.L().Errorf("job:%v has already registered", jon.id)
		return
	}

	jon.forever = 1 == 2
	if reflect.DeepEqual(jon.dead, time.Time{}) {
		jon.forever = true
	}

	m.count++
	m.pool[jon.id] = jon
	m.input <- jon
}

func (m *jobManager) Revoke(id id) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if r, ok := m.pool[id]; ok {
		m.count--
		r.stop <- struct{}{}
		delete(m.pool, id)
	}
}

func (m *jobManager) Start() {
	for j := range m.input {
		tmp := j
		logx.L().Debugf("register ticker for %v", tmp.id)
		go func(tmp *job) {
			defer func() {
				if e := recover(); e != nil {
					logx.T(nil, define.ArchFatal).Errorf("panic for jon:%v", tmp.id)
					return
				}
			}()

			for {
				select {
				case <-tmp.ticker.C:
					if !tmp.forever && time.Now().Sub(tmp.dead) > 0 {
						logx.L().Debugf("dead loop, revoke ticker:%v", tmp.id)
						m.Revoke(tmp.id)
						return
					} else {
						if tmp.callback != nil {
							tmp.callback()
						}
						tmp.f()
					}
				case <-tmp.stop:
					logx.L().Debugf("stop tick, %v", tmp.id)
					return
				case <-tmp.ctx.Done():
					logx.L().Warnf("tick context done, %v", tmp.id)
					return
				}
			}
		}(tmp)
	}
}

func newJob(ctx context.Context, name string, t *time.Ticker, f func(), dead ...time.Time) *job {
	lock.Lock()
	defer lock.Unlock()
	var deadline = time.Time{}
	if len(dead) > 0 {
		deadline = dead[0]
	}
	return &job{
		ctx:    ctx,
		id:     id(name),
		ticker: t,
		f:      f,
		stop:   make(chan struct{}),
		start:  time.Now(),
		dead:   deadline,
	}
}
