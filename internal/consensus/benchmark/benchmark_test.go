// Author: huaxinrui@tal.com
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"testing"
	"time"

	"github.com/huaxr/framework/internal/consensus"
)

func BenchmarkGenID(b *testing.B) {
	consensus.LaunchCampaign()
	time.Sleep(3 * time.Second)
	b.ResetTimer()
	b.Run("a", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := consensus.GetID()
			b.Logf("get id: %v", id)
		}
	})
	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = consensus.GetID()
		}
	})
	b.ReportAllocs()
	b.StopTimer()
}
