// Author: XinRui Hua
// Time:   2023/02/22 10:16
// Git:    huaxr

package concurrent

import "testing"

func TestGroup_Do(t *testing.T) {
	type args struct {
		key string
		fn  func() (interface{}, error)
	}
	//tests := []struct {
	//	name    string
	//	args    args
	//	want    interface{}
	//	wantErr bool
	//}{}
}
