// Author: huaxinrui@tal.com
// Time:   2021/12/22 下午5:30
// Git:    huaxr

package toolutil

import (
	"os"
	"testing"
	"time"

	"github.com/shirou/gopsutil/process"
)

func TestCpu(t *testing.T) {
	x := GetCpuPercent()
	t.Log(x)

	y := GetMemPercent()
	t.Log(y)

	z := GetDiskPercent()
	t.Log(z)

	go func() {
		select {}
	}()
	t.Log(PodProcess())

}

func TestCp(t *testing.T) {
	var i int
	go func() {
		for {
			i++
		}
	}()
	p, _ := process.NewProcess(int32(os.Getpid()))
	cpuPercent, err := p.Percent(1 * time.Second)

	t.Log(cpuPercent, err)
}
