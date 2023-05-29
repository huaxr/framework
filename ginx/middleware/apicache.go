// Author: huaxr
// Time: 2022-12-15 11:09
// Git: huaxr

package middleware

import (
	"net/http"

	"github.com/huaxr/framework/component/plugin/apicache"
	"github.com/huaxr/framework/internal/define"
	"github.com/gin-gonic/gin"
)

var useCache = true

func OpenCache(open bool) {
	useCache = open
}

func WrapCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if useCache {
			v, ok := apicache.GetApiCache().GetCache().KVGet(ctx.FullPath())
			if ok {
				ctx.JSON(http.StatusOK, v)
				ctx.Abort()
				return
			}
			ctx.Set(define.GinWrapCache, useCache)
		}
		ctx.Next()
	}
}
