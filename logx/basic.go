// Author: huaxr
// Time: 2022-12-14 09:12
// Git: huaxr

package logx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/huaxr/framework/internal/define"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func TraceId(ctx context.Context) zap.Field {
	if v := ctx.Value(define.TraceId.String()); v != nil {
		return zap.String(define.TraceId.String(), cast.ToString(v))
	} else {
		return zap.String(define.TraceId.String(), "")
	}
}

func DurationPair(ctx context.Context) (zap.Field, zap.Field) {
	if v := ctx.Value(define.StartTime.String()); v != nil {
		vv := cast.ToString(v)
		t1 := cast.ToInt64(vv)
		t2 := time.Now().UnixNano()
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float32(t2-t1)/float32(1e9)), 64)
		if value < 0 {
			// timezone different may cause nuance.
			value = -value
		}
		return zap.Float64(define.Duration.String(), value), zap.String(define.StartTime.String(), vv)
	} else {
		return zap.Float64(define.Duration.String(), -1), zap.String(define.StartTime.String(), "")
	}
}

func CallFrom(ctx context.Context) zap.Field {
	if v := ctx.Value(define.CallFrom.String()); v != nil {
		return zap.String(define.CallFrom.String(), cast.ToString(v))
	} else {
		return zap.String(define.CallFrom.String(), "")
	}
}
