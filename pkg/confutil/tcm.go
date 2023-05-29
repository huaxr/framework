// Author: huaxinrui@tal.com
// Time: 2022-10-26 13:50
// Git: huaxr

package confutil

type Tcm struct {
	NamespaceId string   `yaml:"namespaceId"`
	GroupName   string   `yaml:"groupName"`
	AccessKey   string   `yaml:"accessKey"`
	SecretKey   string   `yaml:"secretKey"`
	Endpoint    string   `yaml:"endpoint"`
	ContextPath string   `yaml:"contextPath"`
	Port        int      `yaml:"port"`
	LogDir      string   `yaml:"logDir"`
	CacheDir    string   `yaml:"cacheDir"`
	Files       []string `yaml:"files"`
}
