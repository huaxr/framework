// Author: huaxr
// Time: 2022-12-03 14:23
// Git: huaxr

package grpcx

import "context"

type InterfereHook interface {
	Before(ctx context.Context) context.Context
	After(ctx context.Context)
	Panic(string)
}

type defaultHook struct{}

var _ InterfereHook = defaultHook{}

func (h defaultHook) Before(ctx context.Context) context.Context {
	return ctx
}

func (h defaultHook) After(ctx context.Context) {
	return
}

func (h defaultHook) Panic(string) {
	return
}
