// Author: huaxinrui@tal.com
// Time: 2022-11-18 17:02
// Git: huaxr

package middleware

import (
	"time"

	"github.com/huaxr/framework/internal/metric"
	"github.com/gin-gonic/gin"
)

func TimeMetric() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t1 := time.Now()
		ctx.Next()
		// time.Since(t1)
		metric.MetGin(ctx, ctx.FullPath(), time.Now().Sub(t1).Milliseconds())
	}
}
