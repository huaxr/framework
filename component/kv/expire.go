// Author: huaxr
// Time: 2022-12-15 12:07
// Git: huaxr
package kv

import (
	"sync"
	"time"
)

type val struct {
	data        interface{}
	expiredTime int64
}

type expireCache struct {
	m       map[string]*val
	timeMap map[int64][]string
	lock    *sync.Mutex
	// signal for graceful
	stop   chan struct{}
	delMsg chan *delMsg
}

type delMsg struct {
	keys []string
	t    int64
}

func InitExpireCache() Cache {
	e := expireCache{
		m:       make(map[string]*val),
		lock:    new(sync.Mutex),
		timeMap: make(map[int64][]string),
		stop:    make(chan struct{}),
		delMsg:  make(chan *delMsg, 1<<12),
	}
	go e.run(time.Now().Unix())
	return &e
}

func (e *expireCache) run(now int64) {
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()
	go func() {
		for v := range e.delMsg {
			e.multiDelete(v.keys, v.t)
		}
	}()
	for {
		select {
		case <-t.C:
			now++ // time.Now().Unix()++ verse each secs
			e.lock.Lock()
			if keys, found := e.timeMap[now]; found {
				e.lock.Unlock()
				e.delMsg <- &delMsg{keys: keys, t: now}
			} else {
				e.lock.Unlock()
			}
		case <-e.stop:
			close(e.delMsg)
			return
		}
	}
}

func (e *expireCache) KVSet(key string, value interface{}, expireSeconds time.Duration) {
	// int64(time.Minute.Seconds()) = 60
	if int64(expireSeconds.Seconds()) <= 0 {
		return
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	expiredTime := time.Now().Unix() + int64(expireSeconds.Seconds())
	e.m[key] = &val{
		data:        value,
		expiredTime: expiredTime,
	}
	e.timeMap[expiredTime] = append(e.timeMap[expiredTime], key)
}

func (e *expireCache) KVSize() int {
	return len(e.m)
}

func (e *expireCache) KVGet(key string) (value interface{}, found bool) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if val, ok := e.m[key]; ok {
		if val.expiredTime <= time.Now().Unix() {
			delete(e.m, key)
			return
		}
		return val.data, true
	}
	return
}

func (e *expireCache) multiDelete(keys []string, t int64) {
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.timeMap, t)
	for _, key := range keys {
		delete(e.m, key)
	}
}
