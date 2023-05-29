// Author: XinRui Hua
// Time:   2022/3/17 下午4:49
// Git:    huaxr

package confutil

import (
	"fmt"

	"github.com/huaxr/framework/pkg/confutil/manager"
	"github.com/ghodss/yaml"
)

// Deprecated: LoadYmlFile should rename to InitYmlFile
// it should be called as init instead of each time
func LoadYmlFile(file string, obj interface{}) {
	InitYmlFile(file, obj)
}

// InitYmlFile do not call twice
func InitYmlFile(file string, obj interface{}) {
	if !manager.Has(file) {
		panic(fmt.Sprintf("file %s not exist", file))
	}
	initYml(file, obj)
}

func initYml(file string, obj interface{}) {
	yamlText := manager.GetConfigByKey(file)
	if len(yamlText) == 0 {
		panic(fmt.Sprintf("missing %s", file))
	}
	if err := yaml.Unmarshal(yamlText, obj); err != nil {
		panic(err)
	}
}
