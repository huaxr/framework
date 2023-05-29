// Author: huaxr
// Time: 2022-12-14 09:13
// Git: huaxr

package logx

import (
	"context"
	"fmt"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/pkg/confutil"
	"go.uber.org/zap"
)

// without tag
func L(ctx ...context.Context) *zap.SugaredLogger {
	if confutil.GetDefaultConfig().Log.Disabletags {
		return logger.levelLogger
	}
	if len(ctx) > 0 {
		return genSug(ctx[0]).With(zap.String(define.Tag.String(), ""))
	}
	return logWithoutCtx.With(zap.String(define.Tag.String(), ""))
}

// with tag
// adopting the way pod's collector to log-center kafka
// we can use tags to support metric here
func T(ctx context.Context, tag string) *zap.SugaredLogger {
	if confutil.GetDefaultConfig().Log.Disabletags {
		return logger.levelLogger
	}
	return genSug(ctx).With(zap.String(define.Tag.String(), tag))
}

// with extra
func Ext(ctx context.Context, extra map[string]string) *zap.SugaredLogger {
	if confutil.GetDefaultConfig().Log.Disabletags {
		return logger.levelLogger
	}
	sug := genSug(ctx)

	if v, ok := extra[define.Tag.String()]; ok {
		sug = sug.With(zap.String(define.Tag.String(), v))
		delete(extra, define.Tag.String())
	} else {
		sug = sug.With(zap.String(define.Tag.String(), ""))
	}

	for k, v := range extra {
		sug = sug.With(zap.String(fmt.Sprintf("x_%s", k), v))
	}
	return sug
}
