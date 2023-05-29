package circuit

import (
	"sync"

	"github.com/huaxr/framework/pkg/confutil"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
	"github.com/afex/hystrix-go/hystrix"
)

type breakerSet struct {
	sync.RWMutex
	breakerSet map[string]hystrix.CommandConfig
}

var bs = &breakerSet{
	breakerSet: make(map[string]hystrix.CommandConfig),
}

func GetBreakerSet() *breakerSet {
	return bs
}

// Deprecated: circuit & limiter should define in tcm.
func InitCircuitInstance() {

}

func (bs *breakerSet) UpdateFromTcm(configs *confutil.DynamicConfig) {
	bs.Lock()
	defer bs.Unlock()
	for _, circuit := range configs.Circuits {
		i := circuit
		if len(i.Name) == 0 || i.Timeout == 0 || i.MaxConcurrentRequests == 0 || i.SleepWindow == 0 || i.ErrorPercentThreshold == 0 {
			logx.T(nil, define.ArchError).Infof("err circuit configuration:%+v", i)
			continue
		}
		bs.breakerSet[i.Name] = hystrix.CommandConfig{
			Timeout:                i.Timeout,
			MaxConcurrentRequests:  i.MaxConcurrentRequests,
			RequestVolumeThreshold: i.RequestVolumeThreshold,
			SleepWindow:            i.SleepWindow,
			ErrorPercentThreshold:  i.ErrorPercentThreshold,
		}
	}
	hystrix.Configure(bs.breakerSet)
}

func (bs *breakerSet) Size() int {
	return len(bs.breakerSet)
}

// when f1 return err then f2 will be called
// the f1 will never return err by default except timeout、circuit open raised
// custom usage should take attention to the name which should not match rpc method name.
func (bs *breakerSet) Monitor(name string, f1 func() error, f2 func(error) error) error {
	bs.RLock()
	defer bs.RUnlock()
	if _, ok := bs.breakerSet[name]; !ok {
		return f1()
	}
	return hystrix.Do(name, f1, f2)
}
