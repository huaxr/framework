// Author: huaxinrui@tal.com
// Time: 2022-12-04 12:13
// Git: huaxr

package promethu

import (
	"context"
	"fmt"
	"time"

	"github.com/huaxr/framework/pkg/toolutil/ip"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"
	"github.com/spf13/cast"

	"github.com/huaxr/framework/pkg/httputil"

	"github.com/huaxr/framework/component/ticker"

	"github.com/huaxr/framework/internal/metric"
	"github.com/huaxr/framework/pkg/confutil"
	"github.com/huaxr/framework/pkg/toolutil"
)

var url = func() string {
	if confutil.GetDefaultConfig().Monitor != nil && len(confutil.GetDefaultConfig().Monitor.Url) > 0 {
		return fmt.Sprintf("%s/collector", confutil.GetDefaultConfig().Monitor.Url)
	}
	return "http://app.xxx.com/monitor/collector"
}()

var interval = func() time.Duration {
	if confutil.GetDefaultConfig().Monitor != nil && confutil.GetDefaultConfig().Monitor.Interval > 0 {
		return time.Duration(confutil.GetDefaultConfig().Monitor.Interval)
	}
	return 10
}()

type Generator struct {
	Cli *httputil.HttpCli
}

func (s *Generator) TickHeartbeat() *ticker.T {
	return ticker.NewT(s.GenerateV3, "collector_metric", time.NewTicker(interval*time.Second))
}

type Elem struct {
	Tag   string  `json:"tag"`
	Value float64 `json:"value"`
}

type Body struct {
	Psm  string `json:"psm"`
	Pod  string `json:"pod"`
	Data []Elem `json:"data"`
}

func (s *Generator) GenerateV1() {
	var body = Body{
		Psm: confutil.GetDefaultConfig().PSM.String(),
		Pod: ip.GetIp(),
		Data: []Elem{
			{Tag: "memPercent", Value: metric.GetMemPercent()},
			{Tag: "diskPercent", Value: metric.GetDiskPercent()},
			{Tag: "cpuPercent", Value: metric.GetCpuPercent()},
		}}

	if v, ok := metric.GetRedisHitPercent(); ok {
		body.Data = append(body.Data, Elem{Tag: "redisHitPercent", Value: v})
	}
	_, _ = s.Cli.Post(context.Background(), url, body)
}

func (s *Generator) GenerateV2() {
	if !confutil.GetDefaultConfig().Env.IsOnline() {
		return
	}
	psm := confutil.GetDefaultConfig().PSM.String()
	pod := ip.GetIp()
	body := map[string]float64{
		fmt.Sprintf("%s-%s-%s", psm, pod, "memPercent"):  metric.GetMemPercent(),
		fmt.Sprintf("%s-%s-%s", psm, pod, "diskPercent"): metric.GetDiskPercent(),
		fmt.Sprintf("%s-%s-%s", psm, pod, "cpuPercent"):  metric.GetCpuPercent(),
	}
	if v, ok := metric.GetRedisHitPercent(); ok {
		body[fmt.Sprintf("%s-%s-%s", psm, pod, "redisHitPercent")] = v
	}
	_, _ = s.Cli.Post(context.Background(), url, body)
}

func (s *Generator) GenerateV3() {
	if !confutil.GetDefaultConfig().Env.IsOnline() {
		return
	}
	c, m, g, err := toolutil.PodProcess()
	if err != nil {
		logx.T(nil, define.ArchError).Error(err)
		return
	}
	psm := confutil.GetDefaultConfig().PSM.String()
	pod := ip.GetIp()
	body := map[string]string{
		fmt.Sprintf("%s-%s-%v", psm, pod, define.CpuPercent):     cast.ToString(c),
		fmt.Sprintf("%s-%s-%v", psm, pod, define.MemPercent):     cast.ToString(m),
		fmt.Sprintf("%s-%s-%v", psm, pod, define.GoroutineCount): cast.ToString(g),
	}
	if v, ok := metric.GetRedisHitPercent(); ok {
		body[fmt.Sprintf("%s-%s-%s", psm, pod, define.RedisHitPercent)] = cast.ToString(v)
	}
	if v, ok := metric.GetCountWithClear(metric.Mysql, int64(interval)); ok {
		body[fmt.Sprintf("%s-%s-%s", psm, pod, define.MysqlAvgQps)] = cast.ToString(v)
	}
	if v, ok := metric.GetCountWithClear(metric.Kafka, int64(interval)); ok {
		body[fmt.Sprintf("%s-%s-%s", psm, pod, define.KafkaErrEps)] = cast.ToString(v)
	}
	_, err = s.Cli.Post(context.Background(), url, body)
	if err != nil {
		logx.L().Errorf("err while post monitor gateway:%v", err)
		return
	}
}
