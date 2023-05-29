// Author: huaxinrui@tal.com
// Time: 2022-11-21 14:52
// Git: huaxr

package ginx

import (
	"github.com/huaxr/framework/pkg/toolutil"
)

var recoverHook func(string)

type logger struct {
	b chan string
}

func setRecoverHook(f func(string)) {
	recoverHook = f
}

// gin middleware io.W interface
func (l *logger) Write(p []byte) (n int, err error) {
	s := toolutil.Bytes2string(p)
	if recoverHook != nil {
		recoverHook(s)
	}
	l.b <- s + "..."
	return
}
