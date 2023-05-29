// Author: huaxr
// Time: 2022-12-03 20:57
// Git: huaxr

package metric

import (
	"fmt"
	"strconv"
	"testing"
)

func TestUint(t *testing.T) {
	//var a uint8 = 256 // overflow
	var a uint8 = 1<<8 - 1
	a++
	t.Log(a)

	a += 100
	t.Log(a)

	var aa uint64 = 129
	var b uint64 = 8911

	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float32(aa)/float32(b)), 64)
	t.Log(value)

}
