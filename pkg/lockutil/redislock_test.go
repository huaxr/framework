// Author: huaxr
// Time: 2022-10-25 16:18
// Git: huaxr

package lockutil

import (
	"context"
	"testing"
	"time"

	"github.com/huaxr/framework/component/kv"
)

//2022-10-25 16:35:43	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:44	miss the lock, try again
//2022-10-25 16:35:44	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:44	miss the lock, try again
//2022-10-25 16:35:45	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:44	miss the lock, try again
//2022-10-25 16:35:46	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:44	miss the lock, try again
//2022-10-25 16:35:47	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:44	miss the lock, try again
//2022-10-25 16:35:48	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:44	miss the lock, try again
//2022-10-25 16:35:49	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:59	deadLockCheck now del the lock
//2022-10-25 16:35:49	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:50	release the lock!
//2022-10-25 16:35:49	INFO	/Users/huaxinrui/go/src/github.com/huaxr/framework/pkg/lockutil/redislock.go:34	got the lock!
func TestRedisLock1(t *testing.T) {
	kv.InitRedisInstances()
	ctx := context.Background()
	redisCli, _ := kv.GetEngineByIndex(1)

	TryLock(ctx, redisCli, "lock", 10*time.Second)
}

func TestRedisLock2(t *testing.T) {
	kv.InitRedisInstances()
	ctx := context.Background()
	redisCli, _ := kv.GetEngine()

	TryLock(ctx, redisCli, "lock", 10*time.Second)
}
