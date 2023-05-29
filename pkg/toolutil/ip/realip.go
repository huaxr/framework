// Author: XinRui Hua
// Time:   2022/3/18 下午6:28
// Git:    huaxr

package ip

import (
	"fmt"
	"net"
)

var (
	privateBlocks []*net.IPNet
	realIp        string
)

func init() {
	for _, b := range []string{"10.0.0.0/8", "100.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
	ip, err := getIp()
	if err != nil {
		panic(err)
	} else {
		realIp = ip
	}
}

func isPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

// GetRealIp returns a real ip
func getIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get interface addresses! Err: %v", err)
	}

	var ipAddr []byte

	for _, rawAddr := range addrs {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		if ip.To4() == nil {
			continue
		}

		if !isPrivateIP(ip.String()) {
			continue
		}

		ipAddr = ip
		break
	}

	if ipAddr == nil {
		return "", fmt.Errorf("no private IP address found, and explicit IP not provided")
	}

	return net.IP(ipAddr).String(), nil
}

func GetIp() string {
	//if len(realIp) == 0 {
	//	var err error
	//	realIp, err = getIp()
	//	if err != nil {
	//		panic(err)
	//	}
	//	return realIp
	//}
	return realIp
}
