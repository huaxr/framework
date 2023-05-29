// Author: huaxr
// Time:   2021/8/26 下午6:38
// Git:    huaxr

package limiter

import (
	"fmt"
	"sync"
	"time"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/confutil"
	"golang.org/x/time/rate"
)

type rateLimitSet struct {
	sync.RWMutex
	rateLimiterSet map[string]Limiter
}

var rs = &rateLimitSet{
	rateLimiterSet: make(map[string]Limiter),
}

func GetLimiterSet() *rateLimitSet {
	return rs
}

type RateLimiter struct {
	limit *rate.Limiter
	app   string
	quota int32
}

func (rl *RateLimiter) GetQuota() int32 {
	return rl.quota
}

func (rl *RateLimiter) GetType() RateType {
	return GoRate
}

func (rl *RateLimiter) Request() bool {
	return rl.limit.Allow()
}

func (rs *rateLimitSet) newRateLimiter(app string, eps int32) Limiter {
	if len(app) == 0 || eps <= 0 {
		logx.T(nil, define.ArchError).Infof(fmt.Sprintf("error limiter parameter, app:%v, eps:%v", app, eps))
		return nil
	}
	rs.Lock()
	defer rs.Unlock()
	if r, ok := rs.rateLimiterSet[app]; ok {
		if r.GetQuota() == eps {
			return r
		}
	}
	rl := NewRateLimiter(app, eps)
	rs.rateLimiterSet[app] = rl
	return rl
}

func (rs *rateLimitSet) GetRateLimiter(app string) (Limiter, bool) {
	rs.RLock()
	defer rs.RUnlock()
	if r, ok := rs.rateLimiterSet[app]; ok {
		return r, true
	}
	return nil, false
}

func (rs *rateLimitSet) UpdateFromTcm(configs *confutil.DynamicConfig) {
	for _, limit := range configs.Limiters {
		i := limit
		rs.newRateLimiter(i.Path, int32(i.Eps))
	}
}

func (rs *rateLimitSet) Size() int {
	return len(rs.rateLimiterSet)
}

func NewRateLimiter(app string, eps int32) Limiter {
	// every 1 millisecond put a token in this bucket
	limit := rate.Every(time.Duration(1e6/eps) * time.Microsecond)
	// the second param is the bucket size.
	limiter := rate.NewLimiter(limit, int(eps)*2)

	rl := new(RateLimiter)
	rl.limit = limiter
	rl.app = app
	rl.quota = eps
	return rl
}
