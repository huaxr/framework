// Author: huaxinrui@tal.com
// Time:   2021/8/27 上午11:00
// Git:    huaxr

package selector

import "sync/atomic"

// roundRobinSelector selects servers with roundrobin.
type roundRobinSelector struct {
	servers []string
	r       *int32
}

func (s roundRobinSelector) Select() string {
	if len(s.servers) == 0 {
		return ""
	}
	i := *s.r
	i = i % int32(len(s.servers))

	atomic.AddInt32(s.r, 1)
	if *s.r >= int32(len(s.servers)) {
		atomic.StoreInt32(s.r, 0)
	}
	return s.servers[i]
}

func newRoundRobinSelector(servers []string) Selector {
	r := int32(0)
	return &roundRobinSelector{servers: servers, r: &r}
}
