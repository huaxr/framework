// Author: huaxr
// Time: 2022-12-08 18:49
// Git: huaxr

package client

import (
	"fmt"

	"google.golang.org/grpc"
)

type ServiceImpl interface {
	fmt.Stringer

	register(addr string) error
	unregister(addr string)

	// do not close!!
	GetConn() (*grpc.ClientConn, error)
	// do not close!!
	GetConnByIp(ip string) (*grpc.ClientConn, error)

	// need close, GetConn is recommended by the way
	GetNewDialConn() (*grpc.ClientConn, error)
	// need close, GetConnByIp is recommended
	GetNewDialConnByIp(ip string) (*grpc.ClientConn, error)

	// GetIpByModKey get ip then use GetConnByIp to fetch your connection.
	GetIpByModKey(key string) (string, error)
	// you can define client dial options here such as retry times.
	// for the config definition detail please refer:
	// https://www.lixueduan.com/posts/grpc/09-retry/
	// you must define `export GRPC_GO_RETRY=on` in your env first.
	SetRetryConfig(string)
	// NewService first then call Run
	Run() error
}
