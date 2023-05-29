package middleware

import (
	"github.com/huaxr/framework/internal/define"
	"github.com/gin-gonic/gin"
)

func prepare(ctx *gin.Context, pt define.PT, f func() string) {
	if v := ctx.GetHeader(pt.String()); v != "" {
		ctx.Set(pt.String(), v)
		return
	}

	if v := ctx.Value(pt.String()); v == nil {
		// pass through fields
		ctx.Set(pt.String(), f())
	}

	// if header contains is_test,
	// for the purpose of pressure-test ctx set flag to us shadow-library
	// ctx.Value("x_is_test")
}

// curl -H "traceId: dayu_123" 127.0.0.1:8888/test/11
func gateWayTrace(ctx *gin.Context) bool {
	if v := ctx.GetHeader("traceId"); v != "" {
		ctx.Set(define.TraceId.String(), v)
		return true
	}

	return false
}

func MarkPassThrough() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prepare(ctx, define.StartTime, define.Nano)
		prepare(ctx, define.CallFrom, define.Chrome)
		if !gateWayTrace(ctx) {
			prepare(ctx, define.TraceId, define.Uid)
		}
		ctx.Next()
	}
}
