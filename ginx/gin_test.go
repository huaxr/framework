// Author: XinRui Hua
// Time:   2023/01/11 14:39
// Git:    huaxr

package ginx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/huaxr/framework/pkg/httputil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type header struct {
	Key   string
	Value string
}

// PerformRequest for testing gin router.
func PerformRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestEngineHandleContext(t *testing.T) {

	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/v2"
		r.HandleContext(c)
	})
	v2 := r.Group("/v2")
	{
		v2.GET("/", func(c *gin.Context) {})
	}

	assert.NotPanics(t, func() {
		w := PerformRequest(r, "GET", "/")
		assert.Equal(t, 301, w.Code)
	})
}

func TestContext(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	cli := httputil.NewHttpClient(1*time.Second, 0, true)
	res, err := cli.Get(c, "http://www.baidu.com")
	t.Log(res, err)
}
