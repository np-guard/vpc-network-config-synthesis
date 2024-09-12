/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"log"
	"strings"

	"github.com/np-guard/models/pkg/netset"
)

// temporary file, should be implemented in models repo
const commaSeparator = ", "

// given a and b single disjoint cidrs, a<b, return true if a and b are touching
func touching(a, b *netset.IPBlock) bool {
	u := a.Copy().Union(b)
	IPRanges := u.ToIPRanges()
	ranges := strings.Split(IPRanges, commaSeparator)
	return len(ranges) == 1
}

func LastIPAddress(i *netset.IPBlock) *netset.IPBlock {
	IPRanges := i.ToIPRanges()
	ranges := strings.Split(IPRanges, commaSeparator)
	lastRange := ranges[len(ranges)-1]
	sliceLastRange := strings.Split(lastRange, "-")
	endIP := sliceLastRange[len(sliceLastRange)-1]
	ipblock, err := netset.IPBlockFromIPAddress(endIP)
	if err != nil {
		log.Fatal(err)
	}
	return ipblock
}

// both start and end are single IP addresses
func IPBlockFromRange(start, end *netset.IPBlock) *netset.IPBlock {
	s := start.ToIPAddressString() + "-" + end.ToIPAddressString()
	ipblock, err := netset.IPBlockFromIPRangeStr(s)
	if err != nil {
		log.Fatal(err)
	}
	return ipblock
}

func ToCidrs(i *netset.IPBlock) []*netset.IPBlock {
	cidrList := i.ToCidrList()
	res := make([]*netset.IPBlock, len(cidrList))
	for i := range cidrList {
		ipblock, err := netset.IPBlockFromCidrOrAddress(cidrList[i])
		if err != nil {
			log.Fatal(err)
		}
		res[i] = ipblock
	}
	return res
}

func FirstIPAddress(i *netset.IPBlock) *netset.IPBlock {
	ipblock, err := netset.IPBlockFromIPAddress(i.FirstIPAddress())
	if err != nil {
		log.Fatal(err)
	}
	return ipblock
}

// ip is not 255.255.255.255
func NextIP(ip *netset.IPBlock) *netset.IPBlock {
	other := netset.GetCidrAll().Subtract(ip)
	return FirstIPAddress(other.Split()[1])
}

// ip is not 0.0.0.0
func BeforeIP(ip *netset.IPBlock) *netset.IPBlock {
	other := netset.GetCidrAll().Subtract(ip)
	return LastIPAddress(other.Split()[0])
}
