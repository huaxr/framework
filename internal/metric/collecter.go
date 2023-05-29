// Author: huaxr
// Time: 2022-12-01 15:46
// Git: huaxr

package metric

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/huaxr/framework/pkg/toolutil"
)

type CountWithClearType int

const (
	Mysql CountWithClearType = iota
	Kafka
)

var (
	redisGetCount  uint64
	redisMissCount uint64

	mysqlExecuteCount uint64

	kafkaErrorCount uint64
)

func IncRedisGetCount() {
	atomic.AddUint64(&redisGetCount, 1)
}

func IncRedisMissCount() {
	atomic.AddUint64(&redisMissCount, 1)
}

func GetRedisHitPercent() (float64, bool) {
	getCount, missCount := atomic.LoadUint64(&redisGetCount), atomic.LoadUint64(&redisMissCount)
	if getCount == 0 {
		return -1, false
	}
	// if redisGetCount > 1<<32-1 overflow and reset to 1
	// than we should clean miss clout to 0
	if getCount < missCount {
		atomic.StoreUint64(&redisMissCount, 0)
		atomic.StoreUint64(&redisGetCount, 0)
		return -1, false
	}

	value, err := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(getCount-missCount)/float64(getCount)), 64)
	if err != nil {
		return -1, false
	}
	return value * 100, true
}

func IncCountWithClear(typ CountWithClearType) {
	switch typ {
	case Mysql:
		atomic.AddUint64(&mysqlExecuteCount, 1)
	case Kafka:
		atomic.AddUint64(&kafkaErrorCount, 1)
	}
}

func resetCountWithClear(typ CountWithClearType) {
	switch typ {
	case Mysql:
		atomic.StoreUint64(&mysqlExecuteCount, 0)
	case Kafka:
		atomic.StoreUint64(&kafkaErrorCount, 0)
	}
}

func GetCountWithClear(typ CountWithClearType, interval int64) (float64, bool) {
	var c uint64
	switch typ {
	case Mysql:
		c = atomic.LoadUint64(&mysqlExecuteCount)
	case Kafka:
		c = atomic.LoadUint64(&kafkaErrorCount)
	}

	if c == 0 {
		return -1, false
	}
	defer resetCountWithClear(typ)
	value, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(c)/float64(interval)), 64)
	if err != nil {
		return -1, false
	}
	return value, true
}

func GetCpuPercent() float64 {
	x := toolutil.GetCpuPercent()
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", x), 64)
	return value
}

func GetMemPercent() float64 {
	x := toolutil.GetMemPercent()
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", x), 64)
	return value

}

func GetDiskPercent() float64 {
	x := toolutil.GetDiskPercent()
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", x), 64)
	return value
}
