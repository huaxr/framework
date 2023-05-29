// Author: XinRui Hua
// Time:   2023/01/05 10:35
// Git:    huaxr

package ip

import "testing"

func TestSort(t *testing.T) {
	var ips = []string{"10.10.10.13:8000", "10.10.10.11:8000", "10.10.10.12:8000"}
	Sort(ips)
	t.Log(ips)

	ips = []string{"10.10.10.13:", "10.10.10.11", "10.10.10.12:8000", "111"}
	Sort(ips)
	t.Log(ips)

}
