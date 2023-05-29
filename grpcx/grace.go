// Author: huaxinrui@tal.com
// Time: 2022-11-24 09:53
// Git: huaxr

package grpcx

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huaxr/framework/internal/consensus"
	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
)

func (g *Grpcx) preStop(addr string) {
	key := fmt.Sprintf("%s/%s:%d", g.servicePath, addr, g.port)
	cons := consensus.GetConsensusClient()
	err := cons.Delete(key)
	if err != nil {
		logx.T(nil, define.ArchError).Error(err)
	}
	g.Srv.GracefulStop()
	logx.Flush()
}

// k8s deployment configuration should be changed.
// signal not support sub command like supervisor.conf
func (g *Grpcx) listenKill(addr string) {
	signalChan := make(chan os.Signal, 1)
	// when pod deploy on sub-process, SIGTERM from k8s will not received by this channel
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		g.preStop(addr)
		time.Sleep(1 * time.Second)
	}()
}
