package consensus

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/huaxr/framework/internal/metric"

	"github.com/huaxr/framework/internal/define"

	"github.com/huaxr/framework/logx"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

type Consensus struct {
	ctx         context.Context
	client      *clientv3.Client
	electPrefix string
}

func (c *Consensus) Put(key, val string, option clientv3.OpOption) (err error) {
	metric.Metric(define.EtcdPut)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = c.client.Put(ctx, key, val, option)
	if err != nil {
		logx.T(nil, define.ArchFatal).Infof("consensus put err :%v", err)
		panic(err)
	}
	return
}

func (c *Consensus) Get(key string) (res []string) {
	metric.Metric(define.EtcdGet)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := c.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		logx.T(nil, define.ArchFatal).Infof("consensus get err :%v", err)
		return
	}
	for _, ev := range resp.Kvs {
		res = append(res, string(ev.Value))
	}
	return
}

func (c *Consensus) Delete(key string) (err error) {
	metric.Metric(define.EtcdDel)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = c.client.Delete(ctx, key)
	if err != nil {
		logx.T(nil, define.ArchFatal).Infof("consensus del err :%v", err)
		return
	}
	return
}

func (c *Consensus) GetClient() *clientv3.Client {
	return c.client
}

func (c *Consensus) WatchKey(key string) clientv3.WatchChan {
	return c.client.Watch(c.ctx, key, clientv3.WithPrefix())
}

func (c *Consensus) Elect(ctx context.Context) {
	session, err := concurrency.NewSession(c.client, concurrency.WithTTL(8))
	if err != nil {
		// k8s alarm here. why?
		logx.L().Warnf("elect concurrency.NewSession err: %v", err.Error())
		return
	}
	e := concurrency.NewElection(session, c.electPrefix)

	// Campaign puts a value as eligible for the election. It blocks until
	// it is elected, an error occurs, or the context is cancelled.
	if err = e.Campaign(context.TODO(), leader); err != nil {
		// k8s alarm here. why?
		logx.L().Warnf("elect campaign err:%v", err)
		return
	}

	logx.L().Debugf("Elect success")

	atomic.CompareAndSwapInt32(&leaderFlag, 0, 1)

	select {
	case <-session.Done():
		atomic.CompareAndSwapInt32(&leaderFlag, 1, 0)

		logx.L().Warnf("elect expired restart elect")
	}
}
