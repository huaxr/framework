// Author: huaxr
// Time:   2021/9/24 下午6:57
// Git:    huaxr

package kv

import (
	"context"
	"fmt"

	"github.com/huaxr/framework/internal/define"

	"github.com/huaxr/framework/logx"

	"time"

	"github.com/huaxr/framework/pkg/confutil"
	"github.com/go-redis/redis/v8"
)

var (
	redisKvs = make([]*redisCache, 0)
)

type redisCache struct {
	cli *redis.Client
}

// init redis by arch.yml file
func InitRedisInstances() error {
	for _, red := range confutil.GetDefaultConfig().Redis {
		v := red
		client := redis.NewClient(&redis.Options{
			Addr:        v.Host,
			Password:    v.Password, // no password set
			DB:          v.Db,       // use default DB
			PoolSize:    v.Poolsize,
			IdleTimeout: time.Second * time.Duration(v.Idletimeout),
			ReadTimeout: time.Second * time.Duration(v.Readtimeout),
			MaxRetries:  v.MaxRetry,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if _, err := client.Ping(ctx).Result(); err != nil {
			err := fmt.Errorf("init redis err:%v", err)
			logx.T(nil, define.ArchError).Infof("redis init err:%v", err)
			cancel()
			return err
		}

		hook := &redisHook{}
		if v.ShowQuery {
			hook.showQuery = true
		}
		client.AddHook(hook)
		redisKvs = append(redisKvs, &redisCache{
			cli: client,
		})
	}
	return nil
}

func GetEngine() (*redis.Client, error) {
	if len(redisKvs) == 0 {
		return nil, fmt.Errorf("no available redis client")
	}
	return redisKvs[0].cli, nil
}

func GetEngineByIndex(index int) (*redis.Client, error) {
	if len(redisKvs)-1 < index {
		return nil, fmt.Errorf("index out of range")
	}
	return redisKvs[index].cli, nil
}

func (r redisCache) KVGet(key string) (val interface{}, exist bool) {
	res := r.cli.Get(context.Background(), key)
	if res == nil {
		return nil, false
	}
	return res, true
}

func (r redisCache) KVSize() int {
	logx.L().Warn("redis kv has no Size option")
	return -1
}

func (r redisCache) KVSet(key string, val interface{}, duration time.Duration) {
	r.cli.Set(context.Background(), key, val, duration)
}

func (r redisCache) use() {
	//redis.SAdd(ctx, key, n.NodeCode)
	//redis.Expire(ctx, key, duration)
	//mem, err := redis.SMembers(ctx, key).Result()
}
