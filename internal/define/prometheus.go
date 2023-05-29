// Author: huaxr
// Time: 2022-12-07 14:07
// Git: huaxr

package define

type PrometheusKey string

const (
	PodCount        PrometheusKey = "podCount"
	CpuPercent      PrometheusKey = "cpuPercent"
	MemPercent      PrometheusKey = "memPercent"
	GoroutineCount  PrometheusKey = "goroutineCount"
	RedisHitPercent PrometheusKey = "redisHitPercent"
	MysqlAvgQps     PrometheusKey = "mysqlAvgQps"
	KafkaErrEps     PrometheusKey = "kafkaErrAvgEps"

	// collector service
	SnCount PrometheusKey = "snCount"
)

func (p PrometheusKey) String() string {
	return string(p)
}
