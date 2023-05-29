// Author: XinRui Hua
// Time:   2022/3/21 下午6:22
// Git:    huaxr

package client

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/confutil"

	"github.com/huaxr/framework/internal/metric"

	"github.com/huaxr/framework/internal/consensus"
	"github.com/huaxr/framework/internal/define"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

var serviceLock sync.Mutex
var servicePool = make(map[string]*service)

func watcher(service *service, watchChan clientv3.WatchChan) {
	for i := range watchChan {
		// metric here
		for _, e := range i.Events {
			k := string(e.Kv.Key)
			v := string(e.Kv.Value)

			switch e.Type {
			case mvccpb.PUT:
				metric.Metric(define.EtcdWatchPut, fmt.Sprintf("found put %v %v", k, v))
				_ = service.register(v)
			case mvccpb.DELETE:
				// delete option dose not contains v.
				metric.Metric(define.EtcdWatchDel, fmt.Sprintf("found del %v %v", k, v))
				vv := strings.Split(k, "/")
				service.unregister(vv[len(vv)-1])
			}
		}
	}
}

// Deprecated: ServiceFound with demon thread is not satisfy to
// set calling Config, we should use Run() after NewService create srv instance.
// if ServiceFound restart while grpc restart simultaneously,
// it would be some tricks there, because we could not guarantee currency safety.
// for the sake of situation above, ticker checker should be used to
// avoiding connections contaminate.
func ServiceFound(psm string) ServiceImpl {
	s := NewService(psm)
	//s.SetRetryConfig("")
	if err := s.Run(); err != nil {
		panic(err)
	}
	return s
}

// call NewService then call Run() method to start instance(underlying thread).
// you can SetRetryConfig before Run to manager your calling options.
func NewService(psm string) ServiceImpl {
	psm = strings.TrimLeft(psm, "/")
	if err := confutil.PSM(psm).Validate(); err != nil {
		logx.T(nil, define.ArchError).Errorf("err ServiceFound:%v", err)
		os.Exit(1)
	}
	// cut in slash suffix
	psm = fmt.Sprintf("/%s/", psm)
	serviceLock.Lock()
	defer serviceLock.Unlock()
	if s, ok := servicePool[psm]; ok {
		return s
	}
	return newService(psm, consensus.GetConsensusClient())
}
