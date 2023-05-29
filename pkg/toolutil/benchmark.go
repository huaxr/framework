// Author: huaxinrui@tal.com
// Time:   2021/7/28 下午4:11
// Git:    huaxr

package toolutil

import (
	"fmt"
	"math"
	"time"
)

func Benchmark(f func()) {
	best := math.MaxFloat64
	for i := 0; i < 100; i++ {
		start := time.Now()
		f()
		end := time.Now()
		elapsed := end.Sub(start).Seconds()
		if elapsed < best {
			best = elapsed
		}
	}
	fmt.Printf("%.0fms\n", best*1000)
}
