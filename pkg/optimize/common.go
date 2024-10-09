/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"sort"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

type Optimizer interface {
	// attempts to reduce number of SG/nACL rules
	Optimize() (ir.Collection, error)
}

// each IPBlock is a single CIDR. The CIDRs are disjoint.
func SortPartitionsByIPAddrs[T any](p []ds.Pair[*netset.IPBlock, T]) []ds.Pair[*netset.IPBlock, T] {
	cmp := func(i, j int) bool {
		if p[i].Left.FirstIPAddress() == p[j].Left.FirstIPAddress() {
			return p[i].Left.LastIPAddress() < p[j].Left.LastIPAddress()
		}
		return p[i].Left.FirstIPAddress() < p[j].Left.FirstIPAddress()
	}
	sort.Slice(p, cmp)
	return p
}

// returns true if this<other
func LessIPBlock(this, other *netset.IPBlock) bool {
	if this.FirstIPAddress() == this.FirstIPAddress() {
		return this.LastIPAddress() < other.LastIPAddress()
	}
	return this.FirstIPAddress() < other.FirstIPAddress()
}

func IcmpsetPartitions(icmpset *netset.ICMPSet) []netp.ICMP {
	result := make([]netp.ICMP, 0)
	if icmpset.IsAll() {
		icmp, _ := netp.ICMPFromTypeAndCode64WithoutRFCValidation(nil, nil)
		return []netp.ICMP{icmp}
	}

	for _, cube := range icmpset.Partitions() {
		for _, typeInterval := range cube.Left.Intervals() {
			for _, icmpType := range typeInterval.Elements() {
				if cube.Right.Equal(netset.AllICMPCodes()) {
					icmp, _ := netp.ICMPFromTypeAndCode64WithoutRFCValidation(&icmpType, nil)
					result = append(result, icmp)
					continue
				}
				for _, codeInterval := range cube.Right.Intervals() {
					for _, icmpCode := range codeInterval.Elements() {
						icmp, _ := netp.ICMPFromTypeAndCode64WithoutRFCValidation(&icmpType, &icmpCode)
						result = append(result, icmp)
					}
				}
			}
		}
	}
	return result
}

func IcmpRuleToIcmpSet(icmp netp.ICMP) *netset.ICMPSet {
	if icmp.TypeCode == nil {
		return netset.AllICMPSet()
	}
	icmpType := int64(icmp.TypeCode.Type)
	if icmp.TypeCode.Code == nil {
		return netset.NewICMPSet(icmpType, icmpType, int64(netp.MinICMPCode), int64(netp.MaxICMPCode))
	}
	icmpCode := int64(*icmp.TypeCode.Code)
	return netset.NewICMPSet(icmpType, icmpType, icmpCode, icmpCode)
}
