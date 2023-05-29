// Author: huaxr
// Time: 2022-11-10 11:52
// Git: huaxr

package switcher

import (
	"sync"

	"github.com/huaxr/framework/pkg/confutil"
)

var sw = &switcher{
	set: make(map[string]int),
}

type switcher struct {
	sync.RWMutex
	set map[string]int
}

func GetSwitchSet() *switcher {
	return sw
}

func (s *switcher) GetSwitch(switchKey string) int {
	s.RLock()
	s.RUnlock()
	if v, ok := s.set[switchKey]; ok {
		return v
	}
	return -888
}

func (s *switcher) UpdateFromTcm(configs *confutil.DynamicConfig) {
	s.Lock()
	s.Unlock()
	s.set = configs.Switchers
}

func (s switcher) Size() int {
	return len(s.set)
}
