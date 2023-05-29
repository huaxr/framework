// Author: huaxinrui@tal.com
// Time: 2022-10-21 11:32
// Git: huaxr

package lockutil

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/huaxr/framework/logx"

	"github.com/spf13/cast"

	"github.com/go-redis/redis/v8"
)

var locked int32 = 0

// it's a blocking method.
// - avoiding multiply process acquire the lock simultaneously
//   when process died by accident, which setnx but not unlock it.
// - get rid of the possibility of dead lock when a key never be deleted
func TryLock(ctx context.Context, engine *redis.Client, key string, expire time.Duration) {
	deadLockCheck(ctx, engine, key)

	for atomic.LoadInt32(&locked) != 1 {
		now := time.Now()
		expireTime := now.Add(expire).UnixNano()
		ok, _ := engine.SetNX(ctx, key, expireTime, 0).Result()

		if ok {
			logx.L(ctx).Infof("got the lock!")
			r, _ := engine.Get(ctx, key).Result()
			avoidDeadLock, _ := engine.GetSet(ctx, key, expireTime).Result()
			if ok || (now.UnixNano() > cast.ToInt64(r) && now.UnixNano() > cast.ToInt64(avoidDeadLock)) {
				// locked here
				atomic.CompareAndSwapInt32(&locked, 0, 1)
				break
			}
		}
		logx.L(ctx).Infof("miss the lock, try again")
		time.Sleep(3 * time.Second)
		deadLockCheck(ctx, engine, key)
	}
}

func UnLock(ctx context.Context, engine *redis.Client, key string) {
	logx.L().Infof("release the lock!")
	atomic.CompareAndSwapInt32(&locked, 1, 0)
	_, _ = engine.Del(ctx, key).Result()
}

func deadLockCheck(ctx context.Context, engine *redis.Client, key string) {
	now := time.Now()
	r, _ := engine.Get(ctx, key).Result()
	if now.UnixNano() > cast.ToInt64(r) {
		logx.L().Infof("deadLockCheck now del the lock")
		UnLock(ctx, engine, key)
	}
}
