// Author: huaxinrui@tal.com
// Time:   2021/12/22 下午5:26
// Git:    huaxr

package toolutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func readUint(path string) (uint64, error) {
	v, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return parseUint(strings.TrimSpace(string(v)), 10, 64)
}

func parseUint(s string, base, bitSize int) (uint64, error) {
	v, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, err
	}
	return v, nil
}

func Cpu() string {
	percent, _ := cpu.Percent(time.Second, false)
	x := fmt.Sprintf("%.1f", percent[0])
	return fmt.Sprintf("%v%%", x)
}

func Mem() runtime.MemStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

var (
	cpuPeriod, cpuQuota, memLimit uint64
)

func getCpuPercent() (uint64, error) {
	if cpuPeriod != 0 {
		return cpuPeriod, nil
	}
	return readUint("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
}

// k8s 资源限制量，而非所需资源
func getCpuQuota() (uint64, error) {
	if cpuQuota != 0 {
		return cpuQuota, nil
	}
	return readUint("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
}

// 例如： 所需0.1核，内存128MB, 则资源限制为0.2核，内存256MB，为 256*1024*1024=268435456
// 因此，无论是核数还是内存均要除以2
func getMemLimit() (uint64, error) {
	if memLimit != 0 {
		return memLimit, nil
	}
	return readUint("/sys/fs/cgroup/memory/memory.limit_in_bytes")
}

func PodProcess() (float64, float64, int, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, 0, 0, err
	}
	// non-blocking
	cpuPercent, err := p.Percent(2 * time.Second)
	if err != nil {
		return 0, 0, 0, err
	}
	cpuPeriod, err := getCpuPercent()
	if err != nil {
		return 0, 0, 0, err
	}
	cpuQuota, err := getCpuQuota()
	if err != nil {
		return 0, 0, 0, err
	}
	cpuNum := float64(cpuQuota) / float64(cpuPeriod)

	// 通过p.Percent获取到的进程占用机器所有CPU时间的比例除以计算出的核心数即可算出Go进程在容器里对CPU的占比。
	cpValue, err := strconv.ParseFloat(fmt.Sprintf("%.3f", cpuPercent/cpuNum), 64)
	if err != nil {
		return 0, 0, 0, err
	}

	memLimit, err := getMemLimit()
	if err != nil {
		return 0, 0, 0, err
	}
	memInfo, err := p.MemoryInfo()
	if err != nil {
		return 0, 0, 0, err
	}
	// RSS叫常驻内存，是在RAM里分配给进程，允许进程访问的内存量
	mpValue, err := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(memInfo.RSS*100/memLimit)), 64)
	if err != nil {
		return 0, 0, 0, err
	}
	//threadCount := pprof.Lookup("threadcreate").Count()
	gNum := runtime.NumGoroutine()
	return cpValue, mpValue, gNum, nil
}

func GetGoroutine() int {
	gNum := runtime.NumGoroutine()
	return gNum
}
