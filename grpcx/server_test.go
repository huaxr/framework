// Author: huaxr
// Time: 2022-10-31 10:58
// Git: huaxr

package grpcx

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/huaxr/framework/pkg/confutil"

	"github.com/huaxr/framework/cmd/grpcx_test/proto"

	"github.com/huaxr/framework/grpcx/client"
	"github.com/huaxr/framework/logx"
)

type testHook struct{}

func (h testHook) Before(ctx context.Context) context.Context {
	logx.L().Infof("test")
	return ctx
}
func (h testHook) After(ctx context.Context) { return }
func (h testHook) Panic(string)              { return }

func TestServer(t *testing.T) {
	rp := NewGrpcx(context.Background())
	rp.RegisterHook(testHook{})
	rp.Run()
}

var SRV1 client.ServiceImpl

func tt() {
	cli, err := SRV1.GetConn()
	if err != nil {
		logx.L().Error(err)
		return
	}
	client := proto.NewTestRpcClient(cli)
	rand.Seed(time.Now().UnixNano())
	req := rand.Int31n(1 << 30)
	r, err := client.GetCache(context.Background(), &proto.TestReq{Data: req})
	if r.Data != req {
		panic("not safe yoo.")
	}
}

func TestConcurrent(t *testing.T) {
	SRV1 = client.NewService(confutil.GetDefaultConfig().PSM.String())
	SRV1.Run()
	var i = 6500
	for i > 0 {
		go tt()
		i--
	}

	select {}
}
