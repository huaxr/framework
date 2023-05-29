// Author: huaxinrui@tal.com
// Time:   2021/8/27 上午10:30
// Git:    huaxr

package selector

// SelectMode defines the algorithm of selecting a services from candidates.
type SelectMode int

const (
	// RandomSelect is selecting randomly
	RandomSelect SelectMode = iota
	// RoundRobin is selecting by round robin
	RoundRobin
	// WeightedRoundRobin is selecting by weighted round robin
	WeightedRoundRobin
	// ConsistentHash is selecting by hashing
	ConsistentHash
)

// Selector defines selector that selects one service from candidates.
// stereotype for each implements. including RandomSelect,RoundRobin..
type Selector interface {
	Select() string
}

func NewSelector(mode SelectMode, servers []string) Selector {
	switch mode {
	case RandomSelect:
		return newRandomSelector(servers)
	case RoundRobin:
		return newRoundRobinSelector(servers)
	case WeightedRoundRobin:
		return newWeightedRoundRobinSelector(servers)
	case ConsistentHash:
		return newConsistentHashSelector(servers)
	default:
		return newRandomSelector(servers)
	}
}
