// Author: huaxinrui@tal.com
// Time:   2021/8/31 下午5:11
// Git:    huaxr

package limiter

type WindowLimiter struct {
	app   string
	quota int32
	// circle queue
	window     []struct{}
	head, tail int
}

func (rl *WindowLimiter) GetQuota() int32 {
	return rl.quota
}

func (rl *WindowLimiter) Request() bool {
	return false
}

func (rl *WindowLimiter) GetType() RateType {
	return SlidingWindow
}
