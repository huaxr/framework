package promethu

import (
	"github.com/huaxr/framework/pkg/confutil"
	"github.com/prometheus/client_golang/prometheus"
)

func NewManager() *manager {
	return &manager{
		redisHitPercentDesc: prometheus.NewDesc(
			"redisHitPercentDesc",
			"redisHitPercentDesc",
			[]string{},
			prometheus.Labels{"psm": confutil.GetDefaultConfig().PSM.String()}),

		cpuPercentDesc: prometheus.NewDesc(
			"cpuPercentDesc",
			"cpuPercentDesc",
			[]string{},
			prometheus.Labels{"psm": confutil.GetDefaultConfig().PSM.String()}),

		memPercentDesc: prometheus.NewDesc(
			"memPercentDesc",
			"memPercentDesc",
			[]string{},
			prometheus.Labels{"psm": confutil.GetDefaultConfig().PSM.String()}),

		diskPercentDesc: prometheus.NewDesc(
			"diskPercentDesc",
			"diskPercentDesc",
			[]string{},
			prometheus.Labels{"psm": confutil.GetDefaultConfig().PSM.String()}),
	}
}

type manager struct {
	redisHitPercentDesc *prometheus.Desc
	cpuPercentDesc      *prometheus.Desc
	memPercentDesc      *prometheus.Desc
	diskPercentDesc     *prometheus.Desc
}

// Describe simply sends the two Desc in the struct to the channel.
func (m *manager) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.redisHitPercentDesc
	//ch <- m.cpuPercentDesc
	//ch <- m.memPercentDesc
	//ch <- m.diskPercentDesc
}

func (m *manager) Collect(ch chan<- prometheus.Metric) {
	//ch <- prometheus.MustNewConstMetric(m.redisHitPercentDesc, prometheus.GaugeValue, metric.GetRedisHitPercent())
	//ch <- prometheus.MustNewConstMetric(m.cpuPercentDesc, prometheus.GaugeValue, metric.GetCpuPercent())
	//ch <- prometheus.MustNewConstMetric(m.memPercentDesc, prometheus.GaugeValue, metric.GetMemPercent())
	//ch <- prometheus.MustNewConstMetric(m.diskPercentDesc, prometheus.GaugeValue, metric.GetDiskPercent())

}
