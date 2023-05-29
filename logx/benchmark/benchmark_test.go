// Author: huaxr
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"context"
	"testing"

	"github.com/huaxr/framework/pkg/toolutil"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
)

func BenchmarkLogL(b *testing.B) {
	b.Run("l", func(b *testing.B) {
		b.ResetTimer()
		ctx := context.Background()
		ctx = context.WithValue(ctx, define.TraceId.String(), "1111")
		for i := 0; i < b.N; i++ {
			logx.L(ctx).Errorf("grpcx_test:%v", toolutil.GetRandomString(100))
		}
	})

	b.StopTimer()
}

func BenchmarkLogT(b *testing.B) {
	type a struct {
		a string
	}
	aa := a{
		a: toolutil.GetRandomString(100000),
	}
	b.Run("t", func(b *testing.B) {
		b.ResetTimer()
		ctx := context.Background()
		ctx = context.WithValue(ctx, define.TraceId.String(), "1111")
		for i := 0; i < b.N; i++ {
			logx.T(ctx, "BenchMarkTest").Debugf("grpcx_test:%v", aa)
		}
	})
	b.StopTimer()
}

func BenchmarkLogExt(b *testing.B) {
	b.Run("ext", func(b *testing.B) {
		b.ResetTimer()
		ctx := context.Background()
		ctx = context.WithValue(ctx, define.TraceId.String(), "1111")
		for i := 0; i < b.N; i++ {
			logx.Ext(nil, map[string]string{"a": "a", "b": "b"}).Errorf("grpcx_test:%v", toolutil.GetRandomString(100))
		}
	})

	b.StopTimer()
}
