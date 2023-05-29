// Author: XinRui Hua
// Time:   2023/01/05 10:24
// Git:    huaxr

package ip

import (
	"bytes"
	"net"
	"sort"
	"strings"
)

func replace(ipPort string) string {
	ss := strings.Split(ipPort, ":")
	if len(ss) != 2 {
		return ipPort
	}
	return ss[0]
}

// compare return true if ip1 < ip2
func compare(ipStr1, ipStr2 string) bool {
	ipStr1, ipStr2 = replace(ipStr1), replace(ipStr2)
	ip1 := net.ParseIP(ipStr1)
	ip2 := net.ParseIP(ipStr2)
	if ip1 == nil || ip2 == nil {
		return ipStr1 < ipStr2
	}
	return iPCompare(ip1, ip2)
}

func iPCompare(ip1, ip2 net.IP) bool {
	if bytes.Compare(ip1.To16(), ip2.To16()) < 0 {
		return true
	}
	return false
}

func Sort(ips []string) {
	sort.Slice(ips, func(i, j int) bool {
		return compare(ips[i], ips[j])
	})
}

func SortIP(ips []net.IP) {
	sort.Slice(ips, func(i, j int) bool {
		return iPCompare(ips[i], ips[j])
	})
}
