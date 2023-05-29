// Author: XinRui Hua
// Time:   2022/3/24 下午2:47
// Git:    huaxr

package transport

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/confutil"
	"github.com/huaxr/framework/pkg/toolutil"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

func genDefaultMD(ctx context.Context) metadata.MD {
	defaultMd := metadata.MD{}

	for _, i := range define.CtxPassThrough {
		switch i {
		case define.CallFrom:
			defaultMd.Set(i.String(), string(confutil.GetDefaultConfig().PSM))
		default:
			if v := ctx.Value(i.String()); v != nil {
				defaultMd.Set(i.String(), cast.ToString(v))
			}
		}
	}
	return defaultMd
}

func wrapTimeout(ctx context.Context, timeX ...time.Duration) context.Context {
	if len(timeX) > 0 {
		ctx, _ = context.WithTimeout(ctx, timeX[0])
	} else {
		ctx, _ = context.WithTimeout(ctx, defaultTimeout)
	}
	return ctx
}

// Deprecated: please use CtxToGRpcCtxWithOneField or CtxToGRpcCtxWithManyFields or GinCtxToGRpc instead to
// pass-through your data, otherwise use CtxToGRpcCtx to pass default infos.
func CtxToGRpcCtxWithField(ctx context.Context, field string, value []byte) context.Context {
	logx.L(ctx).Warn("CtxToGRpcCtxWithField Deprecated")
	defaultMd := genDefaultMD(ctx)
	if len(field) > 0 {
		defaultMd.Set(field, toolutil.Bytes2string(value))
	}
	return metadata.NewOutgoingContext(ctx, defaultMd)
}

// no kv defined, then you call this method to transform context.
// there is no ctx need to trans while kafka-consumer -> rpc,
// cause rpc is not only called by gin, but also some other options.
// which means it is not required before calling rpc with ctx prepared.
func CtxToGRpcCtx(ctx context.Context, timeX ...time.Duration) context.Context {
	defaultMd := genDefaultMD(ctx)
	return wrapTimeout(metadata.NewOutgoingContext(ctx, defaultMd), timeX...)
}

// one kv to pass
func CtxToGRpcCtxWithOneField(ctx context.Context, field string, value string, timeX ...time.Duration) context.Context {
	defaultMd := genDefaultMD(ctx)
	if len(field) > 0 {
		defaultMd.Set(field, value)
	}
	return wrapTimeout(metadata.NewOutgoingContext(ctx, defaultMd), timeX...)
}

// many kv to pass
func CtxToGRpcCtxWithManyFields(ctx context.Context, fields map[string]string, timeX ...time.Duration) context.Context {
	md := metadata.Join(genDefaultMD(ctx), metadata.New(fields))
	return wrapTimeout(metadata.NewOutgoingContext(ctx, md), timeX...)
}

// convert all gin ctx(when you call c.Set) to grpc metadata
// considering that there may be a lot of kv (slim the rpc frame body),
// if you only need to specify pass-through,
// please use CtxToGRpcCtxWithField or CtxToGRpcCtxWithFields

// warning: performance will be relatively poor than others
func CtxToGRpcCtxWithAllField(ctx *gin.Context, timeoutX ...time.Duration) context.Context {
	md := metadata.MD{}
	for k, v := range ctx.Keys {
		// do not pass CacheWrapper flag
		if v == nil || k == define.GinWrapCache {
			continue
		}
		if k == define.CallFrom.String() {
			md.Set(k, string(confutil.GetDefaultConfig().PSM))
			continue
		}
		// nil v will panic here
		switch reflect.TypeOf(v).Kind() {
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			md[k] = append(md[k], cast.ToString(v))
		default:
			b, err := json.Marshal(v)
			if err != nil {
				logx.T(ctx, define.ArchError).Errorf("CtxToGRpcCtxWithAllField marshal key err:%v, value:%v", k, v)
				continue
			}
			md[k] = append(md[k], toolutil.Bytes2string(b))
		}
	}
	return wrapTimeout(metadata.NewOutgoingContext(ctx, md), timeoutX...)
}
