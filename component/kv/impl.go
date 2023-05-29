// Author: huaxr
// Time:   2021/12/6 下午4:55
// Git:    huaxr

package kv

import "time"

// KVCache
type Cache interface {
	KVGet(key string) (val interface{}, exist bool)
	KVSet(key string, val interface{}, duration time.Duration)
	KVSize() int
}
