// Author: XinRui Hua
// Time:   2022/4/27 下午9:59
// Git:    huaxr

package toolutil

import "math/rand"

// Random buried point
func Ab() bool {
	return rand.Float64()*100 >= 50
}
