package iputil

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

// https://blog.ip2location.com/knowledge-base/how-to-convert-ip-address-range-into-cidr/

func Ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}

	return binary.BigEndian.Uint32(ip)
}

func Int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)

	return ip
}

func iMask(s int) uint32 {
	rem := math.Pow(2, 32) - math.Pow(2, (32-float64(s)))
	ret := math.Round(rem)

	return uint32(ret)
}

func Iprange2Cidr(Start, End string) []string {
	start := Ip2int(net.ParseIP(Start))
	end := Ip2int(net.ParseIP(End))

	iplist := make([]string, 0)

	for end >= start {
		var maxSize uint8 = 32
		for maxSize > 0 {
			mask := iMask(int(maxSize - 1))
			maskBase := start & mask

			if maskBase != start {
				break
			}

			maxSize--
		}

		x := math.Log(float64(end-start+1)) / math.Log(2)
		maxDiff := uint8(32 - math.Floor(x))

		if maxSize < maxDiff {
			maxSize = maxDiff
		}

		ip := Int2ip(start)
		cidr := fmt.Sprintf("%s/%d", ip, maxSize)
		iplist = append(iplist, cidr)

		rem := math.Pow(2, float64(32-maxSize))
		start += uint32(rem)
	}

	return iplist
}
