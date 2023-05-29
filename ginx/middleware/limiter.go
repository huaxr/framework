// Author: huaxr
// Time: 2022-11-09 10:41
// Git: huaxr

package middleware

import (
	"fmt"

	"github.com/huaxr/framework/component/plugin/limiter"
	"github.com/huaxr/framework/ginx/response"
	"github.com/huaxr/framework/internal/metric"
	"github.com/gin-gonic/gin"
)

func Limiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.FullPath()
		l, ok := limiter.GetLimiterSet().GetRateLimiter(path)
		if ok && !l.Request() {
			metric.MetGinLimit(ctx, path)
			response.Error(ctx, fmt.Errorf("qps limited by %v, max size:%v", path, l.GetQuota()))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
