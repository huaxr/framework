package circuit

import (
	"testing"
	"time"

	"github.com/huaxr/framework/pkg/confutil"
	"github.com/huaxr/framework/pkg/toolutil"
)

func TestCircuitDemo(t *testing.T) {
	GetBreakerSet().UpdateFromTcm(&confutil.DynamicConfig{
		Circuits: []*confutil.Circuit{
			&confutil.Circuit{
				Name:                   "monitor",
				Timeout:                10,
				MaxConcurrentRequests:  10,
				RequestVolumeThreshold: 10,
				SleepWindow:            1000,
				ErrorPercentThreshold:  10,
			},
		},
		Limiters:  nil,
		Caches:    nil,
		Switchers: nil,
	})

	for i := 0; i < 20; i++ {
		err := GetBreakerSet().Monitor("monitor", func() error {
			if toolutil.Ab() {
				time.Sleep(21 * time.Millisecond)
			}
			return nil
		}, func(err error) error {
			t.Log(err)
			return err
		})
		time.Sleep(100 * time.Millisecond)

		t.Log(err)
	}
}
