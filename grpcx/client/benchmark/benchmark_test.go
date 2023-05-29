// Author: huaxr
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"testing"

	"github.com/huaxr/framework/grpcx/client"
)

var srv client.ServiceImpl

func BenchmarkClient(b *testing.B) {
	b.ResetTimer()
	b.Run("a", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			srv = client.NewService("test.test.test")
			srv.Run()
		}
	})
	b.ReportAllocs()
	b.StopTimer()
}

func BenchmarkClientGetConn(b *testing.B) {
	b.ResetTimer()
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			srv.GetConn()
		}
	})
	b.ReportAllocs()
	b.StopTimer()
}
