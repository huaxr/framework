// Author: huaxinrui@tal.com
// Time:   2021/6/25 上午10:26
// Git:    huaxr

package consensus

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/confutil"
	"github.com/coreos/etcd/clientv3"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	once      sync.Once
	consensus *Consensus
)

func LaunchIdGenerate(maxNodeIDPrefix, nodeIDPrefix string) {
	cli := GetConsensusClient()
	InitIdGenerate(cli.ctx, maxNodeIDPrefix, nodeIDPrefix)
}

func LaunchCampaign() {
	cli := GetConsensusClient()
	go wait.UntilWithContext(cli.ctx, cli.Elect, time.Second*5)
}

func newConsensus(cx context.Context) error {
	hosts := confutil.GetDefaultConfig().Grpc.GetHosts()
	ctx, _ := context.WithCancel(cx)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(hosts, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	// Test connection timeout
	ctxTmp, _ := context.WithTimeout(context.Background(), 3*time.Second)
	if _, err := cli.Get(ctxTmp, "/i_am_not_exist"); err != nil {
		// not key not exist, res is empty but err is nil
		logx.T(nil, define.ArchError).Infof("etcd try connect err, local environment? forgot vpn?")
		return errors.New("connect etcd error")
	}

	consensus = new(Consensus)
	consensus.client = cli
	consensus.ctx = ctx

	return nil
}

func GetConsensusClient() *Consensus {
	if consensus != nil {
		return consensus
	}
	once.Do(func() {
		if err := newConsensus(context.Background()); err != nil {
			logx.T(nil, define.ArchError).Infof("Err: %v", err)
			os.Exit(1)
		}
	})
	return consensus
}
