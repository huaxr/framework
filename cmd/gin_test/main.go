package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huaxr/framework/cmd/grpcx_test/proto"

	"github.com/huaxr/framework/component/plugin/apicache"

	"github.com/huaxr/framework/ginx"
	"github.com/huaxr/framework/ginx/response"
	"github.com/huaxr/framework/grpcx/client"
	"github.com/huaxr/framework/grpcx/transport"
	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/confutil"
	"github.com/gin-gonic/gin"
)

// 和grpc交互用
var SRV1 client.ServiceImpl

type testController struct {
	// 这里也可以定义一些限流器、熔断器之类的
}

func (ctl *testController) Test(c *gin.Context) {
	logx.T(c, "TestCall").Infof("Hi i am server1")
	var err error
	//// 发现一个rpc节点
	//var SRV1 = client.ServiceFound(confutil.GetDefaultConfig().PSM.String())
	cli, err := SRV1.GetConn()
	//cli, err := SRV1.GetNewDialConn()
	if err != nil {
		logx.L(c).Errorf("err %v", err)
		response.Error(c, err)
		return
	}
	client := proto.NewTestRpcClient(cli)

	var x = 111111111111111
	c.Set("zzz", &x)
	cc := transport.CtxToGRpcCtxWithAllField(c, 3*time.Second)
	// grpc 内部解析有 state 和 frame， frame 控制了参数例如上述的xxxx， state控制了状态，例如下面的 timeout
	// 操作ctx 会把所有信息wrap起来发送给 server 端， 因此实现起来非常方便
	//cc, _ = context.WithTimeout(cc, 1*time.Second)

	var res interface{}
	// max retry 用于控制连接的
	res, err = client.GetCache(cc, &proto.TestReq{Data: 1})
	if err != nil {
		logx.L(c).Errorf("resp:%v err %v", res, err)
		response.Error(c, err)
		return
	}
	response.Success(c, res)
}

func (ctl *testController) Test3(c *gin.Context) {
	panic(2222222222)
	response.Success(c, nil)
}

func (ctl *testController) Router(router *gin.Engine) {
	entry := router.Group("/test")
	entry.GET("/11", ctl.Test)
	entry.GET("/22", ctl.Test3)

}

func xx() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

type x struct {
	A string
	B string
}

func main() {
	transport.SetDefaultTimeout(1 * time.Second)

	SRV1 = client.NewService("/" + confutil.GetDefaultConfig().PSM.String())
	SRV1.Run()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
		os.Exit(1)
	}()

	go func() {
		time.Sleep(3 * time.Second)
		apicache.GetApiCache().UpdateFromTcm(&confutil.DynamicConfig{
			Circuits: nil,
			Limiters: nil,
			Caches: []*confutil.Cache{
				{
					Path:     "/test/11",
					Duration: 10,
				},
			},
			Switchers: nil,
		})
	}()

	g := ginx.NewGinx(ctx)
	g.DisableMonitor()
	g.RegisterMiddleware(xx())
	g.RegisterRouter(&testController{})
	//g.RegisterPanicHook(func(s string) {
	//	logx.L().Error("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	//})
	g.Run()
}
