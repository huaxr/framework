package main

import (
	"context"

	handler2 "github.com/huaxr/framework/cmd/grpcx_test/handler"
	proto2 "github.com/huaxr/framework/cmd/grpcx_test/proto"

	"github.com/huaxr/framework/logx"

	"github.com/huaxr/framework/grpcx"
)

type defaultHook struct{}

func (h defaultHook) Before(ctx context.Context) context.Context {
	logx.L().Infof("aaaaaaaaaaa")
	return ctx
}

func (h defaultHook) After(ctx context.Context) {
	return
}

func (h defaultHook) Panic(string) {
	return
}

func main() {
	rp := grpcx.NewGrpcx(context.Background())
	proto2.RegisterTestRpcServer(rp.Srv, &handler2.TestService{})
	rp.RegisterHook(defaultHook{})
	rp.Run()
}
