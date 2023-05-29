// Author: huaxr
// Time: 2022-12-04 11:47
// Git: huaxr

package promethu

import (
	"testing"

	"github.com/huaxr/framework/internal/metric"
	"github.com/huaxr/framework/pkg/toolutil"
)

func BenchmarkColl(b *testing.B) {
	b.Run("a", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			metric.GetRedisHitPercent()
			toolutil.Cpu()
			metric.GetMemPercent()
			metric.GetDiskPercent()
		}
	})
	b.StopTimer()
}
