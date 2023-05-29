package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func DebugCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// dynamic using the given origin . when using "*" which will disable cookie by chrome save reasons
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", strings.Join([]string{"content-type", "JWT"}, ","))
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, gin.H{"error_code": 0, "err_msg": nil, "data": "Options Request Success!"})
			c.Abort()
			return
		}
	}
}
