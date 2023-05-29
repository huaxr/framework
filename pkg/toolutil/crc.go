// Author: huaxr
// Time:   2021/6/7 下午1:04
// Git:    huaxr

package toolutil

import (
	"hash/crc32"
)

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func CRC(s []byte) int {
	v := int(crc32.ChecksumIEEE(s))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	return 0
}
