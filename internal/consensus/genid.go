// Author: huaxinrui@tal.com
// Time:   2021/6/8 下午5:59
// Git:    huaxr

package consensus

import (
	"context"
	"fmt"

	"github.com/huaxr/framework/pkg/toolutil/ip"

	"strconv"

	"time"

	"github.com/huaxr/framework/logx"

	"github.com/bwmarrin/snowflake"
	v3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	NodeIDLimit = 1 << 10
)

// Generator
type Generator struct {
	nodeID  int64
	snow    *snowflake.Node
	session *concurrency.Session
	// prefix define here
	maxNodeIDPrefix, nodeIDPrefix string
}

var (
	defaultGenerator *Generator
)

// Run UntilWithContext loops until context is done, running f every period.
func (n *Generator) Run(ctx context.Context) {
	wait.UntilWithContext(ctx, n.start, time.Second*5)
}

func (n *Generator) start(ctx context.Context) {
	logx.L().Debugf("start loop id generator, %v", time.Now())
	err := n.register(ctx)
	if err != nil {
		logx.L().Errorf("register err: %v", err.Error())
		n.reset()
		return
	}
	n.snow, err = snowflake.NewNode(n.nodeID)
	if err != nil {
		logx.L().Errorf("snowflake err %v", err.Error())
		n.reset()
		return
	}

	for {
		select {
		// when the lease is orphaned, expires, or is otherwise no longer being refreshed.
		case <-n.session.Done():
			logx.L().Warnf("session done %v", "Generator session expire")
			n.reset()
			return
		case <-ctx.Done():
			logx.L().Warnf("ctx done %v", "Generator context done")
			n.reset()
			return
		}
	}
}

func (n *Generator) register(ctx context.Context) (e error) {
	session, err := concurrency.NewSession(consensus.client, concurrency.WithTTL(20))
	if err != nil {
		logx.L().Errorf("register err %v", err)
		return err
	}
	getRes, err := session.Client().Get(ctx, n.maxNodeIDPrefix)
	if err != nil {
		logx.L().Errorf("register err %v", err)
		return err
	}
	maxStr := "0"
	if getRes.Count > 0 {
		maxStr = string(getRes.Kvs[0].Value)
	}

	maxID, err := strconv.ParseInt(maxStr, 10, 64)
	if err != nil {
		logx.L().Errorf("register err %v", err)
		return err
	}

	nodeID := (maxID + 1) % NodeIDLimit
	k := fmt.Sprintf("%s/%d", n.nodeIDPrefix, nodeID)
	addr := ip.GetIp()
	txRes, err := session.Client().Txn(ctx).
		If(v3.Compare(v3.CreateRevision(k), "=", 0)).
		Then(v3.OpPut(k, fmt.Sprintf("%s", addr),
			v3.WithLease(session.Lease())),
			v3.OpPut(n.maxNodeIDPrefix, strconv.FormatInt(nodeID, 10))).
		Commit()
	if err != nil {
		logx.L().Errorf("register err %v", err)
		return err
	}
	if !txRes.Succeeded {
		logx.L().Errorf("register not success %v", err)
		return fmt.Errorf("failed to register nodeID: %d", nodeID)
	}
	n.nodeID = nodeID
	n.session = session
	return nil
}

func (n *Generator) reset() {
	if n.session != nil {
		n.session.Close()
		n.session = nil
	}
	n.nodeID = 0
}

func InitIdGenerate(ctx context.Context, maxNodeIDPrefix, nodeIDPrefix string) {
	defaultGenerator = &Generator{
		maxNodeIDPrefix: maxNodeIDPrefix,
		nodeIDPrefix:    nodeIDPrefix,
	}
	go defaultGenerator.Run(ctx)
}

func GetID() int64 {
	id := defaultGenerator.snow.Generate()
	return int64(id)
}
