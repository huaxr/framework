// Author: huaxr
// Time:   2021/12/16 上午10:55
// Git:    huaxr

package toolutil

import "encoding/base64"

func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base64Decode(b string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b)

}
