// Author: huaxinrui@tal.com
// Time: 2022-11-15 17:22
// Git: huaxr

package toolutil

import "runtime"

func GetStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return Bytes2string(buf[:n])
}
