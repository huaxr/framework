package confutil

import (
	"sync"

	"github.com/huaxr/framework/pkg/confutil/manager"
)

var arch *ArchConfig
var onceDefault sync.Once

type ArchConfig struct {
	// basic stuff should be filled in yml
	Env     Env      `yaml:"env"`
	PSM     PSM      `yaml:"psm"`
	Log     *Log     `yaml:"log"`
	Grpc    *Grpc    `yaml:"grpc"`
	Gin     *Gin     `yaml:"gin"`
	Mysql   []*Mysql `yaml:"mysql"`
	Redis   []*Redis `yaml:"redis"`
	Tcm     *Tcm     `yaml:"tcm"`
	Monitor *Monitor `yaml:"monitor"` // host addr
	Pprof   *Pprof   `yaml:"pprof"`
	// queue is a utility function that users define in their own configuration
	// which will no longer use arch.yml define
	//Queue *Queue `yaml:"queue"`
	//Circuit []*Circuit `yaml:"circuit"`
}

func GetDefaultConfig() *ArchConfig {
	if arch == nil {
		onceDefault.Do(func() {
			manager.Init(*confDir, []string{"yml"})

			arch = new(ArchConfig)
			initYml("arch.yml", &arch)
			if arch.Log == nil {
				arch.Log = defaultLog
			}
			// project name can be loaded to prepare some import options
			if err := arch.PSM.Validate(); err != nil {
				panic(err)
			}

			if len(arch.Env) == 0 {
				arch.Env = Local
			}
			// extra validation not defined!
			// client usage should be limited
		})
	}

	return arch
}

func SetConfig(c *ArchConfig) {
	arch = c
}
