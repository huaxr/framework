// Author: huaxr
// Time: 2022-12-05 12:27
// Git: huaxr

package ginx

import (
	"github.com/huaxr/framework/ginx/middleware"
)

func (r *Gin) DisableMonitor() {
	r.openMonitor = false
}

func (r *Gin) DisableTcm() {
	r.openTcm = false
}

func (r *Gin) DisablePprof() {
	r.openPprof = false
}

func (r *Gin) DisableApiCache() {
	middleware.OpenCache(false)
}
