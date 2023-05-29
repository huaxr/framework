// Author: huaxinrui@tal.com
// Time: 2022-10-29 14:00
// Git: huaxr

package metric

import (
	"context"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
	"go.uber.org/zap"
)

var defaultMetricM = logx.WithArch().
	With(zap.Float64(define.Duration.String(), -1)).
	With(zap.String(define.TraceId.String(), "")).
	With(zap.String(define.StartTime.String(), "")).
	With(zap.String(define.CallFrom.String(), "")).
	With(zap.String(define.HandlerPath.String(), "")).
	With(zap.Int64(define.HandlerExecutePeriod.String(), -1))

func Metric(tag define.ArchTag, msg ...string) {
	if len(msg) > 0 {
		defaultMetricM.With(zap.String(define.Tag.String(), string(tag))).Info(msg[0])
	} else {
		defaultMetricM.With(zap.String(define.Tag.String(), string(tag))).Info()
	}
}

func wrapCtx(ctx context.Context) *zap.SugaredLogger {
	return logx.WithArch().With(logx.TraceId(ctx)).With(logx.DurationPair(ctx)).With(logx.CallFrom(ctx))
}

// metric two tags at same time
func MetRpc(ctx context.Context, method string, ms int64) {
	wrapCtx(ctx).With(zap.String(define.Tag.String(), string(define.GrpcQps))).
		With(zap.String(define.HandlerPath.String(), method)).
		With(zap.Int64(define.HandlerExecutePeriod.String(), ms)).
		Info()
}

func MetRpcFusing(ctx context.Context, method string) {
	wrapCtx(ctx).With(zap.String(define.Tag.String(), string(define.GrpcCircuitOpen))).
		With(zap.String(define.HandlerPath.String(), method)).
		With(zap.Int64(define.HandlerExecutePeriod.String(), -1)).
		Info()
}

func MetGin(ctx context.Context, uri string, ms int64) {
	wrapCtx(ctx).With(zap.String(define.Tag.String(), string(define.GinxQps))).
		With(zap.String(define.HandlerPath.String(), uri)).
		With(zap.Int64(define.HandlerExecutePeriod.String(), ms)).
		Info()
}

func MetGinLimit(ctx context.Context, uri string) {
	wrapCtx(ctx).With(zap.String(define.Tag.String(), string(define.GinxLimit))).
		With(zap.String(define.HandlerPath.String(), uri)).
		With(zap.Int64(define.HandlerExecutePeriod.String(), -1)).
		Info()
}
