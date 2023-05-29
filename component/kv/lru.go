// Author: huaxr
// Time: 2022/6/13 3:37 下午
// Git: huaxr

package kv

import (
	"sync"
	"time"
)

type lruCache struct {
	size       int
	capacity   int
	mutex      sync.RWMutex
	cache      map[string]*node
	Head, Tail *node
}

type node struct {
	key       string
	value     interface{}
	Pre, Next *node
}

func initLinkNode(key string, value interface{}) *node {
	return &node{key, value, nil, nil}
}

func InitLruCache(capacity int) Cache {
	l := lruCache{
		0,
		capacity,
		sync.RWMutex{},
		map[string]*node{},
		initLinkNode("head", 0),
		initLinkNode("tail", 0),
	}
	l.Head.Next = l.Tail
	l.Tail.Pre = l.Head
	return &l
}

func (l *lruCache) KVSize() int {
	return l.size
}

func (l *lruCache) KVGet(key string) (interface{}, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if _, ok := l.cache[key]; !ok {
		return nil, false
	}
	node := l.cache[key]
	l.updateToHead(node)
	return node.value, true
}

func (l *lruCache) KVSet(key string, value interface{}, duration time.Duration) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, ok := l.cache[key]; !ok {
		node := initLinkNode(key, value)
		for l.size >= l.capacity {
			l.deleteLast()
		}
		l.cache[key] = node
		l.insertNewHead(node)
	} else {
		node := l.cache[key]
		node.value = value
		l.updateToHead(node)
	}
}

func (l *lruCache) updateToHead(node *node) {
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	temp := l.Head.Next
	l.Head.Next = node
	node.Pre = l.Head
	node.Next = temp
	temp.Pre = node

}

func (l *lruCache) deleteLast() {
	node := l.Tail.Pre
	l.Tail.Pre = node.Pre
	node.Pre.Next = node.Next
	node.Pre = nil
	node.Next = nil
	l.size--
	delete(l.cache, node.key)
}

func (l *lruCache) insertNewHead(node *node) {
	temp := l.Head.Next
	l.Head.Next = node
	node.Pre = l.Head
	temp.Pre = node
	node.Next = temp
	l.size++
}
