// Author: huaxr
// Time: 2022-12-03 14:22
// Git: huaxr

package kv

import (
	"context"
	"fmt"

	"github.com/huaxr/framework/logx"

	"github.com/spf13/cast"

	"github.com/huaxr/framework/internal/metric"
	"github.com/go-redis/redis/v8"
)

type redisHook struct {
	showQuery bool
}

var _ redis.Hook = redisHook{}

func (redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	// todo: debug query time
	//ctx = context.WithValue(ctx, "_time", time.Now().UnixNano())
	return ctx, nil
}

func (r redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	args := cmd.Args()
	if len(args) > 0 {
		if cast.ToString(args[0]) == "get" {
			metric.IncRedisGetCount()
			// redis: nil error
			if err := cmd.Err(); err != nil {
				metric.IncRedisMissCount()
			}
		}

		if r.showQuery {
			logx.L(ctx).Infof("redis query:" + fmt.Sprint(args))
		}
	}
	return nil
}

func (redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
