// Author: huaxinrui@tal.com
// Time:   2021/6/12 下午12:36
// Git:    huaxr

package consensus

import "sync/atomic"

const leader = "__hitler"

// only leader can process
var leaderFlag int32

func masterFunc(f func()) {
	if atomic.LoadInt32(&leaderFlag) == 1 {
		f()
		return
	}
}
