package misc

import (
	"fmt"
	"strconv"
	"strings"
)

func Ip2long(ipstr string) (ip uint32) {

	ips := strings.Split(ipstr, ".")
	if len(ips) != 4 {
		return
	}
	ip1, e := strconv.ParseUint(ips[0], 0, 32)
	if e != nil {
		return
	}
	ip2, _ := strconv.ParseUint(ips[1], 0, 32)
	if e != nil {
		return
	}
	ip3, _ := strconv.ParseUint(ips[2], 0, 32)
	if e != nil {
		return
	}
	ip4, _ := strconv.ParseUint(ips[3], 0, 32)
	if e != nil {
		return
	}

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return ip
}

func Long2ip(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24, ip<<8>>24, ip<<16>>24, ip<<24>>24)
}
