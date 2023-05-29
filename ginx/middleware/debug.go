// Author: huaxinrui@tal.com
// Time: 2022-11-21 18:46
// Git: huaxr

package middleware

import (
	"bytes"

	"github.com/huaxr/framework/logx"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ginBodyLogMiddleware(c *gin.Context) {
	logx.L(c).Infof("gin request:url%v, method:%v", c.Request.URL, c.Request.Method)
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	statusCode := c.Writer.Status()
	logx.L(c).Infof("statue: %v gin response:%v", statusCode, blw.body.String())
}

func LogRequestResponse() gin.HandlerFunc {
	return ginBodyLogMiddleware
}
