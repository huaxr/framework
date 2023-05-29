package handler

import (
	"context"

	proto2 "github.com/huaxr/framework/cmd/grpcx_test/proto"
)

type TestService struct{}

// 如果函数很长，执行时间很长
// 唯一的办法就是将函数拆件为n个带contenxt的函数，每执行一步可以鉴定contxt done
// 请注意该函数没有工具自动生成！
func (t *TestService) GetCache(ctx context.Context, req *proto2.TestReq) (*proto2.TestResponse, error) {
	//a := []string{}

	return &proto2.TestResponse{
		Data: int32(req.Data),
	}, nil
}
