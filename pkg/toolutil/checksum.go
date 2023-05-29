package toolutil

import (
	"crypto/md5"
	"encoding/hex"
)

// CheckSum function is calculate checksum to md5
func CheckSum(b []byte) string {
	sum := md5.Sum(b)
	src := sum[:]
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return Bytes2string(dst)
}
