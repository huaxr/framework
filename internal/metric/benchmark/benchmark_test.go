// Author: huaxinrui@tal.com
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"testing"

	"github.com/huaxr/framework/internal/metric"
)

func BenchmarkCollector(b *testing.B) {
	b.Run("BenchmarkCollector", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			metric.IncCountWithClear(metric.Mysql)
			metric.GetCountWithClear(metric.Mysql, 10)
		}
		b.StopTimer()
	})
}
