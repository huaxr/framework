// Author: huaxinrui@tal.com
// Time: 2022-11-16 15:00
// Git: huaxr

package logx

import (
	"testing"

	"github.com/huaxr/framework/pkg/toolutil"
)

func TestPath(t *testing.T) {
	Ext(nil, map[string]string{"a": "a", "b": "b"}).Errorf("test:%v", toolutil.GetRandomString(100))
}
