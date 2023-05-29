// Author: huaxinrui@tal.com
// Time: 2022-12-15 12:09
// Git: huaxr
package kv

import (
	"testing"
	"time"
)

var cache = InitExpireCache()

func Test_Map(t *testing.T) {
	cache.KVSet("a", []int{1, 2, 3}, 2*time.Second)
	time.Sleep(1 * time.Second)
	v, ok := cache.KVGet("a")
	t.Log(v, ok)
}

// once Get+Set about 279.4 ns/op
func BenchmarkLogL(b *testing.B) {
	b.Run("expire", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.KVSet("a", []int{1, 2, 3}, 1*time.Second)
			cache.KVSet("b", []int{1, 2, 3}, 10*time.Second)
			cache.KVGet("a")
			cache.KVGet("b")
		}
	})

	time.Sleep(5 * time.Second)
	b.Log(cache.KVSize())

	time.Sleep(5 * time.Second)
	b.Log(cache.KVSize())
	b.StopTimer()
}

func TestTime(t *testing.T) {
	t1 := time.Now().Unix()
	t.Log(t1)
	a := 1 * time.Millisecond
	t.Log(a.Seconds(), int64(a.Seconds()))
	b := 1 * time.Second
	t.Log(b.Seconds(), int64(b.Seconds()))
	c := 1 * time.Minute
	t.Log(c.Seconds(), int64(c.Seconds()))
	time.Sleep(1 * time.Second)
	t2 := time.Now().Unix()
	t.Log(t2)
}
