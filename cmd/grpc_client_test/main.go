// Author: huaxr
// Time: 2022-11-04 21:14
// Git: huaxr

package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/huaxr/framework/cmd/grpcx_test/proto"
	"github.com/huaxr/framework/pkg/confutil"

	"github.com/huaxr/framework/grpcx/client"
	"github.com/huaxr/framework/logx"
)

var SRV1 client.ServiceImpl

func tt() {
	cli, err := SRV1.GetConn()
	if err != nil {
		logx.L().Error(err)
		return
	}
	//defer cli.Close()
	client := proto.NewTestRpcClient(cli)

	rand.Seed(time.Now().UnixNano())

	req := rand.Int31n(1 << 30)
	r, err := client.GetCache(context.Background(), &proto.TestReq{Data: req})

	logx.L().Info(req, r.Data)
	if r.Data != req {
		panic("xxxxxx")
	}
}

func main() {
	SRV1 = client.NewService(confutil.GetDefaultConfig().PSM.String())
	SRV1.Run()

	var i = 100
	for i > 0 {
		go tt()
		i--
	}

	select {}
}
