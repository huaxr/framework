// Author: XinRui Hua
// Time:   2022/12/29 13:53
// Git:    huaxr

package pprofutil

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/confutil"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var port = func() int {
	if confutil.GetDefaultConfig().Pprof != nil && confutil.GetDefaultConfig().Pprof.Port > 0 {
		return confutil.GetDefaultConfig().Pprof.Port
	}
	return 10086
}()

// go tool pprof http://localhost:10086/debug/pprof/profile
// go tool pprof http://localhost:10086/debug/pprof/heap
// open graph
// go tool pprof -http 127.0.0.1:9999 /Users/huaxinrui/pprof/pprof.samples.cpu.001.pb.gz
func Pprof() {
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
		if err != nil {
			logx.T(nil, define.ArchFatal).Errorf("start pprof fail on:%d, err:%v", port, err)
		}
	}()
}

// gin server can register handler in the same srv port.
func GinPprof(engine *gin.Engine) {
	pprof.Register(engine, "/debug/pprof")
}
