// Author: huaxr
// Time: 2022/6/13 3:41 下午
// Git: huaxr

package kv

import (
	"testing"
)

func TestLru(t *testing.T) {
	c := InitLruCache(2)
	c.KVSet("1", "1", 0)
	c.KVSet("2", "2", 0)
	c.KVSet("3", "3", 0)

	v, ok := c.KVGet("1")
	t.Log(v, ok)

	v, ok = c.KVGet("2")
	t.Log(v, ok)

	v, ok = c.KVGet("3")
	t.Log(v, ok)
}
