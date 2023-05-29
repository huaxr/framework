// Author: XinRui Hua
// Time:   2022/3/21 下午5:41
// Git:    huaxr

package consensus

import (
	"context"
	"fmt"
	"strings"

	"github.com/huaxr/framework/internal/metric"

	"github.com/huaxr/framework/internal/define"

	"github.com/huaxr/framework/logx"
	"github.com/coreos/etcd/clientv3"
)

type KeepAlive struct {
	Key string
	Ttl int64
}

func (k KeepAlive) Alive(ctx context.Context) {
	logx.L().Debugf("keepalive for:%v, ttl:%v", k.Key, k.Ttl)
	ctx, cancel := context.WithCancel(ctx)
	sps := strings.Split(k.Key, "/")
	if len(sps) < 2 {
		logx.T(nil, define.ArchError).Infof("key format err")
		return
	}
	consensus := GetConsensusClient()
	client := consensus.GetClient()

	lease := clientv3.NewLease(client)
	leaseResp, err := lease.Grant(ctx, k.Ttl)
	if err != nil {
		logx.T(nil, define.ArchError).Infof("set lease fail:%v", err)
		return
	}
	leaseID := leaseResp.ID

	// lease automatically
	leaseRespChan, err := lease.KeepAlive(ctx, leaseID)
	if err != nil {
		logx.T(nil, define.ArchError).Infof("lease keepalive fail:%v", err)
		return
	}
	err = consensus.Put(k.Key, sps[len(sps)-1], clientv3.WithLease(leaseID))
	if err != nil {
		logx.T(nil, define.ArchError).Infof("put key fail:%v", err)
		return
	}

	for {
		select {
		case leaseKeepResp := <-leaseRespChan:
			if leaseKeepResp == nil {
				metric.Metric(define.EtcdLeaseExpire, fmt.Sprintf("lease expired for :%v", k.Key))
				// represent that the lease expired
				cancel()
				return
			}
		case <-ctx.Done():

			return
		}
	}
}
