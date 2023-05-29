// Author: huaxinrui@tal.com
// Time:   2021/8/27 上午11:00
// Git:    huaxr

package selector

import (
	"fmt"
	"net/url"
	"strconv"
)

// Weighted is a wrapped server with  weight
type Weighted struct {
	Server          string
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

// weightedRoundRobinSelector selects servers with weighted.
type weightedRoundRobinSelector struct {
	servers []*Weighted
}

func (s weightedRoundRobinSelector) Select() string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	w := nextWeighted(ss)
	if w == nil {
		return ""
	}
	return w.Server
}

func createWeighted(servers map[string]string) []*Weighted {
	var ss = make([]*Weighted, 0, len(servers))
	for k, metadata := range servers {
		w := &Weighted{Server: k, Weight: 1, EffectiveWeight: 1}

		if v, err := url.ParseQuery(metadata); err == nil {
			ww := v.Get("weight")
			if ww != "" {
				if weight, err := strconv.Atoi(ww); err == nil {
					w.Weight = weight
					w.EffectiveWeight = weight
				}
			}
		}

		ss = append(ss, w)
	}
	return ss
}

func newWeightedRoundRobinSelector(servers []string) Selector {
	var srv = make(map[string]string)
	for d, i := range servers {
		srv[i] = fmt.Sprintf( "weight=%d", d)
	}
	ss := createWeighted(srv)
	return &weightedRoundRobinSelector{servers: ss}
}

func nextWeighted(servers []*Weighted) (best *Weighted) {
	total := 0

	for i := 0; i < len(servers); i++ {
		w := servers[i]

		if w == nil {
			continue
		}
		//if w is down, continue

		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}

	}

	if best == nil {
		return nil
	}

	best.CurrentWeight -= total
	return best
}
