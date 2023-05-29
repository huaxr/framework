// Author: huaxr
// Time: 2022-11-29 11:13
// Git: huaxr

package transport

import (
	"sync"
	"time"
)

var defaultTimeout = 8 * time.Second
var once sync.Once

// for users setting global default timeout when calling grpc.
func SetDefaultTimeout(t time.Duration) {
	once.Do(func() {
		defaultTimeout = t
	})
}
