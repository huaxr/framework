package confutil

import "time"

type DCS interface {
	UpdateFromTcm(configs *DynamicConfig)
	Size() int
}

// Circuit name: GetCache #熔断器名称
//timeout: 10 #ms #熔断器执行函数超时控制
//maxConcurrentRequests: 100 # 最大并发量
//requestVolumeThreshold: 2 # 一个统计窗口,达到这个请求数量后才去判断是否要开启熔断
//sleepWindow: 800 #ms 熔断器被打开后sleepWindow 的时间就是控制过多久后去尝试服务是否可用了
//errorPercentThreshold: 50 # 错误百分比 请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
type Circuit struct {
	Name                   string `yaml:"name" json:"name"`
	Timeout                int    `yaml:"timeout" json:"timeout"`
	MaxConcurrentRequests  int    `yaml:"maxConcurrentRequests" json:"maxConcurrentRequests"`
	RequestVolumeThreshold int    `yaml:"requestVolumeThreshold" json:"requestVolumeThreshold"`
	SleepWindow            int    `yaml:"sleepWindow" json:"sleepWindow"`
	ErrorPercentThreshold  int    `yaml:"errorPercentThreshold" json:"errorPercentThreshold"`

	timeRange [2]*time.Time
}

type Limit struct {
	Path      string `yaml:"path" json:"path"`
	Eps       int    `yaml:"eps" json:"eps"`
	timeRange [2]*time.Time
}

// Cache api cache
type Cache struct {
	Path string `yaml:"path" json:"path"`
	// cache lru expire seconds
	Duration  int `yaml:"duration" json:"duration"`
	timeRange [2]*time.Time
}

type DynamicConfig struct {
	Circuits  []*Circuit     `json:"circuits"`
	Limiters  []*Limit       `json:"limits"`
	Caches    []*Cache       `json:"caches"`
	Switchers map[string]int `json:"switchers"` // 0,1,2,3
}

/*{
    "circuits":[
        {
            "name":"GetCache",
            "timeout":10,
            "maxConcurrentRequests":100,
            "requestVolumeThreshold":2,
            "sleepWindow":300,
            "errorPercentThreshold":50
        }
    ],
    "limits":[
        {
            "path":"/test/11",
            "eps":1
        }
    ]
}*/

// todo: 自动降级&放火 实现观测，不完全依赖手动降级
// 不考虑在假期步数 agent, 使各个模块具备放火掉件，如通过redis、mysql hook构造演练效果
// 暂时无法更新，望后人承接代码补全相关功能
type Degrade struct {
	// 服务发现降级，关闭服务污染检查
	ServiceFoundDegrade bool `json:"service_found_degrade"`
	// 是否关闭 pprof
	OpenPprof bool `json:"open_pprof"`
	// 是否关闭监控上报
	OpenMonitor bool `json:"open_monitor"`
}
