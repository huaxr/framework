package response

import (
	"net/http"
	"time"

	"github.com/huaxr/framework/component/plugin/apicache"
	"github.com/huaxr/framework/internal/define"

	"github.com/gin-gonic/gin"
)

type (
	respCode int
)

// for to simplify transaction between backend and client,
// to that respCode only support
// 0 success
// 1 normal err
// 2 auth err
// considering the backward compatibility, ErrorWithCode is fine to
// dealing our work.
const (
	SuccessCode respCode = iota
	ErrorCode
	AuthError

	_
	//InternalError
	//UnknownError
)

type Response struct {
	Code respCode    `json:"code"` // 0 1 2
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`

	Code2   int    `json:"code2,omitempty"`
	TraceId string `json:"trace_id,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	respJSON(c, SuccessCode, "", data)
}

func Error(c *gin.Context, error error) {
	respJSON(c, ErrorCode, error.Error(), nil)
}

func ErrorWithCode(c *gin.Context, code respCode, error error) {
	respJSON(c, code, error.Error(), nil)
}

// 业务自定义返回code2 使用此方法，由于向前兼容，命名有些奇怪
func RespJsonWithCode(c *gin.Context, error error, code2 int) {
	traceId := c.GetString(define.TraceId.String())
	response := Response{Code: ErrorCode, Msg: error.Error(), Data: nil, TraceId: traceId, Code2: code2}
	c.JSON(http.StatusOK, response)
}

func respJSON(c *gin.Context, code respCode, msg string, data interface{}) {
	traceId := c.GetString(define.TraceId.String())
	response := Response{Code: code, Msg: msg, Data: data, TraceId: traceId}

	if v, ok := c.Get(define.GinWrapCache); ok && v.(bool) {
		if vv, ok := apicache.GetApiCache().GetExpSec(c.FullPath()); ok {
			apicache.GetApiCache().GetCache().KVSet(c.FullPath(), response, time.Duration(vv)*time.Second)
		}
	}

	c.JSON(http.StatusOK, response)
}
