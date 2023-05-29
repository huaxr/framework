// Author: huaxinrui@tal.com
// Time:   2021/8/27 上午11:41
// Git:    huaxr

package limiter

type TokenLimiter struct {
	app   string
	quota int32
}

func (rl *TokenLimiter) GetQuota() int32 {
	return rl.quota
}

func (rl *TokenLimiter) Request() bool {
	return false
}

func (rl *TokenLimiter) GetType() RateType {
	return LeakyBucket
}
