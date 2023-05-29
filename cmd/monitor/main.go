package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/huaxr/framework/pkg/toolutil/ip"

	"github.com/huaxr/framework/cmd/monitor/model"

	"github.com/huaxr/framework/component/dao/orm"

	"github.com/huaxr/framework/internal/define"

	"github.com/spf13/cast"

	"github.com/huaxr/framework/pkg/confutil"

	"github.com/huaxr/framework/ginx"
	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/toolutil"
	"github.com/gin-gonic/gin"
)

const expire = 33 * time.Second

type value struct {
	Value      string    `json:"value"`
	UpdateTime time.Time `json:"update_time"`
}

type monitorController struct {
	cache map[string]value
	lock  sync.Mutex
	queue chan map[string]string
}

type object struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func (ctl *monitorController) Router(router *gin.Engine) {
	router.GET("/metrics", ctl.MetricV2)
	router.POST("/collector", ctl.Collector)
	router.GET("/export", ctl.Export)
	router.GET("/hosts", ctl.PrometheusFound)
}

type elem struct {
	v float64
	c int
}

func (ctl *monitorController) MetricV2(c *gin.Context) {
	prometheusFormat := "%s{psm=\"%s\"} %v\n"

	var buf = strings.Builder{}
	defer buf.Reset()

	ctl.lock.Lock()
	defer ctl.lock.Unlock()

	var psmKeyMap = make(map[string]*elem)

	var cleanKeys = make([]string, 0)

	for k, v := range ctl.cache {
		if time.Now().Sub(v.UpdateTime) > expire {
			cleanKeys = append(cleanKeys, k)
			continue
		}
		val := strings.Split(k, "-")

		// PSM-HOST-KEY
		if len(val) != 3 {
			logx.L(c).Errorf("err format")
			continue
		}
		psmKey := fmt.Sprintf("%s-%s", val[0], val[2])
		if fv, ok := psmKeyMap[psmKey]; ok {
			fv.c++
			fv.v += cast.ToFloat64(v.Value)
		} else {
			psmKeyMap[psmKey] = &elem{
				v: cast.ToFloat64(v.Value),
				c: 1,
			}
		}
	}

	// clean keys
	for _, i := range cleanKeys {
		delete(ctl.cache, i)
	}

	var psmMap = make(map[string]struct{})
	for k, v := range psmKeyMap {
		ks := strings.Split(k, "-")
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", v.v/float64(v.c)), 64)
		s := fmt.Sprintf(prometheusFormat, ks[1], ks[0], value)
		buf.WriteString(s)

		if _, ok := psmMap[ks[0]]; !ok {
			psmMap[ks[0]] = struct{}{}
			sExtra := fmt.Sprintf(prometheusFormat, define.PodCount.String(), ks[0], v.c)
			buf.WriteString(sExtra)
		}
	}

	c.Data(200, "text/plain", toolutil.String2Byte(buf.String()))
}

func (ctl *monitorController) Collector(c *gin.Context) {
	var body map[string]string
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(200, "fail")
		return
	}
	ctl.queue <- body
	c.JSON(200, "success")
}

func (ctl *monitorController) Export(c *gin.Context) {
	c.JSON(200, ctl.cache)
}

func (ctl *monitorController) PrometheusFound(c *gin.Context) {
	var objs = make([]*object, 0)
	var obj = new(object)
	obj.Targets = []string{fmt.Sprintf("%s:%d", ip.GetIp(), confutil.GetDefaultConfig().Gin.Port)} // 单点监控、后期可以采用rpc分布式缓存，当前业务能抗住
	obj.Labels = map[string]string{"label": "montage"}

	objs = append(objs, obj)
	c.JSON(200, objs)
}

func (ctl *monitorController) job() {
	for {
		body := <-ctl.queue
		ctl.lock.Lock()

		for k, v := range body {
			ctl.cache[k] = value{
				Value:      v,
				UpdateTime: time.Now(),
			}
		}
		ctl.lock.Unlock()
	}
}

func (ctl *monitorController) getSnCount() {
	err := orm.InitDbInstances()
	if err != nil {
		panic(err)
	}
	mysqlCli, err := orm.GetEngine()
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			stu := new(model.Student)
			res, err := mysqlCli.Slave().Distinct("sn").Count(stu)
			if err != nil {
				logx.L().Errorf("err when count: %v", err)
				ctl.queue <- map[string]string{
					fmt.Sprintf("%s-%s-%v", confutil.GetDefaultConfig().PSM, "127.0.0.1", define.SnCount): cast.ToString(0),
				}
				continue
			}
			ctl.queue <- map[string]string{
				fmt.Sprintf("%s-%s-%v", confutil.GetDefaultConfig().PSM, "127.0.0.1", define.SnCount): cast.ToString(res),
			}
		}
	}

}

func main() {
	g := ginx.NewGinx(context.Background())

	r := &monitorController{
		cache: make(map[string]value),
		lock:  sync.Mutex{},
		queue: make(chan map[string]string, 10000),
	}

	go r.job()
	go r.getSnCount()
	g.RegisterRouter(r)
	//g.DisableMonitor()
	if err := g.Run(); err != nil {
		panic(err)
	}
}
