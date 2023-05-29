// Author: huaxinrui@tal.com
// Time:   2021/9/8 下午4:58
// Git:    huaxr

package benchmark

import (
	"testing"

	"github.com/huaxr/framework/pkg/toolutil"
)

func BenchmarkString2Byte(b *testing.B) {
	b.Run("zerocopy", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			elems := "abcdefghigklmn"
			_ = toolutil.String2Byte(elems)
		}

	})
	b.Run("common", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			elems := "abcdefghigklmn"
			_ = []byte(elems)
		}
	})

	b.StopTimer()
}
