// Author: XinRui Hua
// Time:   2022/3/22 下午3:48
// Git:    huaxr

package client

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/huaxr/framework/cmd/grpcx_test/proto"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestServicesV2(t *testing.T) {
	s := NewService("montage.arch.metric")
	//MaxAttempts : 最大尝试次数
	//InitialBackoff : 默认退避时间
	//MaxBackoff : 最大退避时间
	//BackoffMultiplier : 退避时间增加倍率
	//RetryableStatusCodes : 服务端返回什么错误代码才重试
	// todo: 重试机制适配到配置文件中 建议 retry 参数对读接口可以设置一下，但是对写接口最好是不要设置。
	s.SetRetryConfig(`{
		"methodConfig": [{
		  "name": [{"service": "proto.TestRpc","method":"GetCache"}],
	     "max_request_message_bytes": 1024,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".1s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNKNOWN" ]
		  }
		}]}`)

	for {
		if err := s.Run(); err != nil {
			t.Log(err)
			break
		} else {
			t.Log("init service success")
		}
	}

	assertions := require.New(t)
	cli, err := s.GetConn()
	assertions.NotNil(t, err, "GetConn is nil")
	client := proto.NewTestRpcClient(cli)
	var res interface{}
	// timeout 超时参数设置，通常是这么设置的，对于你要调用的系统你要看看他平时调用要多久能返回，然后比正常的耗时设置的多个 50% 就可以了。
	// 比如平时一般正常在 100~200ms，偶尔高峰会在 500ms，那你设置个 timeout=800ms 或者 1s 其实都可以。
	res, err = client.GetCache(context.Background(), &proto.TestReq{Data: 888})
	assertions.NotNil(t, err, "GetCache is nil, res:%v", res)
	select {}
}

func TestServices(t *testing.T) {
	s := ServiceFound("montage.framework.test")
	for {
		time.Sleep(10 * time.Millisecond)
		go func() {
			c, err := s.GetConn()
			if err != nil {
				t.Log(err)
				return
			}
			//c.Close()
			m := c.GetState()
			t.Log(m)
		}()
	}
}

func BenchmarkBalance(b *testing.B) {
	b.Run("balance", func(b *testing.B) {
		b.ResetTimer()
		rpcServices := ServiceFound("aaa.aaa.aaa")
		for i := 0; i < b.N; i++ {
			rpcServices.register("192.168.1.4:9999")
			b.Log(rpcServices)
			rpcServices.unregister("192.168.1.4:9999")
			b.Log(rpcServices)
		}
	})

	b.StopTimer()
}

type serviceTestifyMock struct {
	mock.Mock
}

func (s *serviceTestifyMock) register(addr string) error {
	args := s.Called(addr)
	fmt.Println(addr, args.Get(0))
	return nil
}
func (*serviceTestifyMock) unregister(addr string)             {}
func (*serviceTestifyMock) GetConn() (*grpc.ClientConn, error) { return nil, errors.New("hello") }
func (*serviceTestifyMock) GetConnByIp(ip string) (*grpc.ClientConn, error) {
	return nil, errors.New("hello")
}
func (*serviceTestifyMock) GetNewDialConn() (*grpc.ClientConn, error)              { return nil, nil }
func (*serviceTestifyMock) GetNewDialConnByIp(ip string) (*grpc.ClientConn, error) { return nil, nil }
func (*serviceTestifyMock) GetIpByModKey(key string) (string, error)               { return "", nil }
func (*serviceTestifyMock) Run() error                                             { return errors.New("no") }
func (*serviceTestifyMock) SetRetryConfig(string)                                  {}
func (*serviceTestifyMock) String() string                                         { return "" }
func TestProductServiceImpl_IsProductReservable(t *testing.T) {
	assertions := require.New(t)
	// Register test mocks
	serviceMock := &serviceTestifyMock{}
	serviceMock.On("Run").Return(errors.New("ok"))

	err := serviceMock.Run()
	t.Log(err)
	assertions.Nil(err, "mock success")
}

func TestServiceIps(t *testing.T) {
	s := NewService("montage.arch.metric")
	if err := s.Run(); err != nil {
		t.Log(err)
	}
	t.Log(s.String())
}
