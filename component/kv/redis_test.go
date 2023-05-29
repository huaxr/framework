// Author: huaxinrui@tal.com
// Time:   2021/12/10 下午2:38
// Git:    huaxr

package kv

import (
	"context"

	"github.com/huaxr/framework/internal/metric"

	"github.com/huaxr/framework/logx"

	"testing"
)

func TestC(t *testing.T) {
	InitRedisInstances()
	ctx := context.Background()
	redisKey := "aaa"
	redisCli, _ := GetEngine()

	redisCli.RPush(ctx, redisKey, "1", "@", "3")
	//c := redisCli.Incr(ctx, redisKey).Val()
	//redisCli.Expire(ctx, redisKey, 5*time.Second)
	res := redisCli.LRange(ctx, redisKey, 0, -1)
	a, err := res.Result()
	logx.L().Info("xx", a, err)

	redisKey2 := "bbb"
	redisCli.HSet(ctx, redisKey2, "1", "2", "3", "1")
	xx := redisCli.HGetAll(ctx, redisKey2)
	b, err := xx.Result()
	logx.L().Info("xx", b, err)

	redisKey3 := "ccc"
	redisCli.SAdd(ctx, redisKey3, "1", "2", "3", "1")
	c, err := redisCli.SMembers(ctx, redisKey3).Result()
	logx.L().Info("xx", c, err)

	redisCli.Get(ctx, "test")
}

func TestHook(t *testing.T) {
	_ = InitRedisInstances()
	ctx := context.Background()
	redisCli, _ := GetEngine()
	redisCli.Get(ctx, "not_exist")

	redisCli.Set(ctx, "a", "b", -1)
	redisCli.Get(ctx, "a")
	redisCli.Get(ctx, "a")
	redisCli.Get(ctx, "a")

	hit, ok := metric.GetRedisHitPercent()
	t.Log(hit, ok)
}
