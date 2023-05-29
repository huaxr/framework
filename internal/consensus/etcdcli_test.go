// Author: huaxr
// Time:   2021/6/8 下午7:06
// Git:    huaxr

package consensus

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	NewLeaseErr  = 101
	LeasTtlErr   = 102
	KeepAliveErr = 103
	PutErr       = 104
	GetErr       = 105
	RevokeErr    = 106
)

func TestLease(t *testing.T) {
	ctx := context.Background()
	client := GetConsensusClient().client
	lease := clientv3.NewLease(client)
	leaseResp, err := lease.Grant(ctx, 1)
	if err != nil {
		fmt.Printf("设置租约时间失败:%s\n", err.Error())
		os.Exit(LeasTtlErr)
	}

	//设置续租
	leaseID := leaseResp.ID
	ctx, _ = context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseID)
	if err != nil {
		fmt.Printf("续租失败:%s\n", err.Error())
		os.Exit(KeepAliveErr)
	}
	//监听租约
	go func() {
		for {
			select {
			case leaseKeepResp := <-leaseRespChan:
				if leaseKeepResp == nil {
					fmt.Printf("已经关闭续租功能\n")
					return
				} else {
					fmt.Printf("续租成功\n")
					goto END
				}
			}
		END:
			//
			time.Sleep(1500 * time.Millisecond)
		}

	}()

	//监听某个key的变化
	//ctx1, _ := context.WithTimeout(context.TODO(),20)
	go func() {
		wc := client.Watch(context.TODO(), "/job/v3/1", clientv3.WithPrevKV())
		for v := range wc {
			for _, e := range v.Events {
				fmt.Printf("type:%v kv:%v  prevKey:%v val:%v \n ", e.Type, string(e.Kv.Key), e.PrevKv, e.Kv.Value)
			}
		}
	}()

	kv := clientv3.NewKV(client)
	//通过租约put
	putResp, err := kv.Put(context.TODO(), "/job/v3/1", "koock", clientv3.WithLease(leaseID))
	if err != nil {
		fmt.Printf("put 失败：%s", err.Error())
		os.Exit(PutErr)
	}
	fmt.Printf("%v\n", putResp.Header)

	//cancelFunc()

	time.Sleep(2 * time.Second)
	_, err = lease.Revoke(context.TODO(), leaseID)
	if err != nil {
		fmt.Printf("撤销租约失败:%s\n", err.Error())
		os.Exit(RevokeErr)
	}
	fmt.Printf("撤销租约成功")
	getResp, err := kv.Get(context.TODO(), "/job/v3/1")
	if err != nil {
		fmt.Printf("get 失败：%s", err.Error())
		os.Exit(GetErr)
	}
	fmt.Printf("%v", getResp.Kvs)
	time.Sleep(20 * time.Second)

}

func TestGet(t *testing.T) {
	ctx := context.Background()
	client := GetConsensusClient().client
	res, _ := client.Get(ctx, "/services/", clientv3.WithPrefix())

	t.Log(res)
}

func TestWatch(t *testing.T) {
	ctx := context.Background()
	client := GetConsensusClient().client

	var count int
	for i := 0; i < 1000; i++ {
		go func() {
			c := client.Watch(ctx, "/tmp")
			for {
				<-c
				count++
			}
		}()
	}

	for {
		client.Put(ctx, "/tmp", "aaa")
		t.Log(count)
		time.Sleep(5 * time.Second)
	}
}
