// Author: huaxr
// Time: 2022-11-07 13:47
// Git: huaxr

package confutil

import (
	"os"
	"runtime"
)

type Env string

const (
	Local  Env = "local"
	Test   Env = "test"
	Gray   Env = "gray"
	Online Env = "online"
)

func (e Env) IsOnline() bool {
	sysType := runtime.GOOS
	return !(sysType == "darwin") && e == Online
}

func e() {
	s := os.Getenv("KUBERNETES_ENV")
	switch s {
	case "online":
	case "gray":

	}
}
