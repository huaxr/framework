// Author: huaxinrui@tal.com
// Time: 2022-10-31 10:35
// Git: huaxr

package consensus

import (
	"context"
	"testing"
	"time"

	"github.com/huaxr/framework/logx"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Alive1(ctx context.Context) {
	logx.L().Info("i am alive")
	return
}

func Alive2(ctx context.Context) {
	logx.L().Info("i am alive")
	select {}
}

func Alive3(ctx context.Context) {
	ctx2, cancel := context.WithCancel(ctx)

	cancel()
	logx.L().Info("i am alive")
	select {
	case <-ctx2.Done():
		logx.L().Error("ctx done")
	default:
	}

}

func TestAlive(t *testing.T) {
	for {
		wait.UntilWithContext(context.Background(), Alive3, time.Second*2)
	}
}
