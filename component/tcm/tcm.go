package tcm

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/huaxr/framework/component/plugin/apicache"

	"github.com/huaxr/framework/component/plugin/circuit"
	"github.com/huaxr/framework/component/plugin/limiter"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/internal/metric"

	"github.com/huaxr/framework/pkg/confutil"

	"github.com/huaxr/framework/logx"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	configClient *ConfigClient
	once         sync.Once
)

type clientOption func(option *clientOptions)

type clientOptions struct {
	ctx         context.Context
	namespaceId string
	groupName   string
	endpoint    string
	contextPath string
	accessKey   string
	secretKey   string
	port        int
	logDir      string
	cacheDir    string
}

type ConfigClient struct {
	opt    clientOptions
	mu     sync.RWMutex
	data   map[string]string
	client config_client.IConfigClient
	signal chan struct{}
}

func contains(s []string, str string) bool {
	for _, i := range s {
		if i == str {
			return true
		}
	}
	return false
}

func InitTcmInstance(cc ...context.Context) {
	var ctx = context.Background()
	if len(cc) == 0 {
		ctx = cc[0]
	}
	once.Do(func() {
		if confutil.GetDefaultConfig().Tcm == nil {
			logx.L().Debugf("current service not use tcm, the circuit and limiter tools disabled")
			return
		}
		if len(confutil.GetDefaultConfig().Tcm.NamespaceId) == 0 || len(confutil.GetDefaultConfig().Tcm.GroupName) == 0 {
			return
		}
		//logx.L().Debugf("executing InitTcmInstance...")
		configClient := newConfigClient(ctx,
			withNamespaceId(confutil.GetDefaultConfig().Tcm.NamespaceId),
			withGroupName(confutil.GetDefaultConfig().Tcm.GroupName),
			withAccessKey(confutil.GetDefaultConfig().Tcm.AccessKey),
			withSecretKey(confutil.GetDefaultConfig().Tcm.SecretKey),
			withEndpoint(confutil.GetDefaultConfig().Tcm.Endpoint),
			withContextPath(confutil.GetDefaultConfig().Tcm.ContextPath),
			withPort(confutil.GetDefaultConfig().Tcm.Port),
			withLogDir(confutil.GetDefaultConfig().Tcm.LogDir),
			withCacheDir(confutil.GetDefaultConfig().Tcm.CacheDir),
		)

		configClient.signal = make(chan struct{}, 10)
		files := confutil.GetDefaultConfig().Tcm.Files
		if !contains(files, confutil.GetDefaultConfig().PSM.String()) {
			// system config
			files = append(files, confutil.GetDefaultConfig().PSM.String())
		}
		for _, file := range files {
			configClient.loadConfig(file)
		}

		setClient(configClient)
		go configClient.cheer()
	})

}

func newConfigClient(ctx context.Context, opts ...clientOption) *ConfigClient {
	option := clientOptions{
		ctx: ctx,
	}
	for _, o := range opts {
		o(&option)
	}

	configClient := &ConfigClient{
		opt:  option,
		mu:   sync.RWMutex{},
		data: make(map[string]string),
	}
	configManagerClient := configClient.initClient()
	configClient.client = configManagerClient

	return configClient
}

func setClient(client *ConfigClient) {
	configClient = client
}

func withNamespaceId(namespaceId string) clientOption {
	return func(o *clientOptions) {
		o.namespaceId = namespaceId
	}
}

func withGroupName(groupName string) clientOption {
	return func(o *clientOptions) {
		o.groupName = groupName
	}
}

func withEndpoint(endpoint string) clientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func withContextPath(contextPath string) clientOption {
	return func(o *clientOptions) {
		o.contextPath = contextPath
	}
}

func withAccessKey(accessKey string) clientOption {
	return func(o *clientOptions) {
		o.accessKey = accessKey
	}
}

func withSecretKey(secretKey string) clientOption {
	return func(o *clientOptions) {
		o.secretKey = secretKey
	}
}

func withPort(port int) clientOption {
	return func(o *clientOptions) {
		o.port = port
	}
}

func withLogDir(logDir string) clientOption {
	return func(o *clientOptions) {
		o.logDir = logDir
	}
}

func withCacheDir(cacheDir string) clientOption {
	return func(o *clientOptions) {
		o.cacheDir = cacheDir
	}
}

func (c *ConfigClient) initClient() config_client.IConfigClient {
	//服务地址配置
	sc := []constant.ServerConfig{
		{
			IpAddr:      c.opt.endpoint,
			Port:        uint64(c.opt.port),
			ContextPath: c.opt.contextPath,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         c.opt.namespaceId, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		AccessKey:           c.opt.accessKey,
		SecretKey:           c.opt.secretKey,
		LogDir:              c.opt.logDir,
		CacheDir:            c.opt.cacheDir,
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "info",
	}
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(fmt.Errorf("tcm config set but connect failed:%v", err))
	}

	return client
}

func (c *ConfigClient) loadConfig(name string) {
	if len(name) == 0 {
		return
	}
	content, err := c.client.GetConfig(vo.ConfigParam{
		DataId: name,
		Group:  c.opt.groupName,
	})
	if err != nil {
		logx.T(nil, define.ArchError).Infof("tcm loadConfig err name:%v, err:%v", name, err)
		return
	}
	c.saveData(name, content)
	//logx.L().Debugf("tcm config name:%v, value:%v", name, content)
	if confutil.GetDefaultConfig().PSM.Equals(name) && len(content) > 0 {
		logx.L().Infof("tcm content:%v", content)
		metric.Metric(define.TcmUpdate, name)
		c.signal <- struct{}{}
	}

	err = c.client.ListenConfig(vo.ConfigParam{
		DataId: name,
		Group:  c.opt.groupName,
		OnChange: func(namespace, group, dataId, data string) {
			logx.L().Debugf("tcm content:%v", content)
			metric.Metric(define.TcmUpdate, name)
			// only psm change signal should be send.
			if confutil.GetDefaultConfig().PSM.Equals(name) {
				c.signal <- struct{}{}
			}
			logx.L().Infof("tcm config reload for key:%v, val:%v", name, data)
			c.saveData(name, data)
		},
	})

	if err != nil {
		logx.T(nil, define.ArchError).Infof("tcm ListenConfig err name:%v, err:%v", name, err)
		return
	}
}

func (c *ConfigClient) saveData(name, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[name] = content
}

// when your config is string
func (c *ConfigClient) GetConfig(name string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	configData, ok := c.data[name]
	if !ok {
		return ""
	}
	return configData
}

func (c *ConfigClient) cheer() {
	for {
		select {
		case <-c.signal:
			var dc confutil.DynamicConfig
			err := GetJsonConfig(confutil.GetDefaultConfig().PSM.String(), &dc)
			if err != nil {
				logx.T(nil, define.ArchError).Infof("err load dc:%v, raw:%v", err, GetTextConfig(confutil.GetDefaultConfig().PSM.String()))
				continue
			}
			circuit.GetBreakerSet().UpdateFromTcm(&dc)
			limiter.GetLimiterSet().UpdateFromTcm(&dc)
			apicache.GetApiCache().UpdateFromTcm(&dc)
		case <-c.opt.ctx.Done():
			logx.L().Info("tcm ctx done")
			return
		}
	}
}

// notice: you can't use v as a persistence variable cause
// naco onChange is dynamic
func GetJsonConfig(name string, v interface{}) error {
	configData := configClient.GetConfig(name)
	return json.Unmarshal([]byte(configData), v)
}

func GetTextConfig(name string) string {
	configData := configClient.GetConfig(name)
	return configData
}

func GetClient() *ConfigClient {
	return configClient
}
