// Author: huaxr
// Time:   2021/8/18 下午6:49
// Git:    huaxr

package toolutil

import (
	"fmt"

	"testing"
)

// key size must 16
var key string = "2/18208s2*_820+q"

func TestAesDecrypt(t *testing.T) {
	str := "hello world fafsaaaaaaaats ohs;da;           sdasfadgfaf"
	es, err := AesEncrypt(str, []byte(key))
	if err != nil {
		t.Log("xx", err)
	}
	fmt.Println(es)
	ds, err := AesDecrypt(es, []byte(key))
	if err != nil {
		t.Log("xx", err)
	}
	fmt.Println(string(ds))

	select {}
}
