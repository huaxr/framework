// Author: huaxinrui@tal.com
// Time: 2022-12-05 12:26
// Git: huaxr

package grpcx

// not register but using service calling
func (g *Grpcx) DisableFound() {
	g.openFound = false
}

func (g *Grpcx) DisableTcm() {
	g.openTcm = false
}

func (g *Grpcx) DisableMonitor() {
	g.openMonitor = false
}

func (g *Grpcx) DisablePprof() {
	g.openPprof = false
}

func (g *Grpcx) DisableApiCache() {
	useCache = false
}
