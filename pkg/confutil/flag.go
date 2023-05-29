// Author: XinRui Hua
// Time:   2022/3/17 下午4:03
// Git:    huaxr

package confutil

import (
	"flag"
	"testing"
)

var (
	confDir *string
)

func contains(s []string, val string) bool {
	for _, i := range s {
		if i == val {
			return true
		}
	}
	return false
}

// define your flag before launching the program.
// notice. on testing, dir should be absolute path
// to illustrate please refer to
func init() {
	confDir = flag.String("dir", "", "config yml dir")
	testing.Init()

	if !flag.Parsed() {
		flag.Parse()
	}

	if len(*confDir) == 0 {
		panic("flag dir not init yet")
	}
}
