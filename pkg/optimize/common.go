/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"slices"

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
func SortPartitionsByIPAddrs[T ds.Set[T]](p []ds.Pair[*netset.IPBlock, T]) []ds.Pair[*netset.IPBlock, T] {
	cmp := func(i, j ds.Pair[*netset.IPBlock, T]) int {
		return i.Left.Compare(j.Left)
	}
	slices.SortFunc(p, cmp)
	return p
}

// IcmpsetPartitions breaks the set into ICMP slice, where each element defined as legal in nACL, SG rules
func IcmpsetPartitions(icmpset *netset.ICMPSet) []netp.ICMP {
	if icmpset.IsAll() {
		icmp, _ := netp.ICMPFromTypeAndCode64WithoutRFCValidation(nil, nil)
		return []netp.ICMP{icmp}
	}

	result := make([]netp.ICMP, 0)
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

func IsAllPorts(tcpudpPorts *netset.PortSet) bool {
	return tcpudpPorts.Equal(netset.AllPorts())
}

// UncoveredHole returns true if the rules can not be continued between the two cubes
// i.e there is a hole between two ipblocks that is not a subset of anyProtocol cubes
func UncoveredHole(prevIPBlock, currIPBlock, anyProtocolCubes *netset.IPBlock) bool {
	touching, _ := prevIPBlock.TouchingIPRanges(currIPBlock)
	if touching {
		return false
	}
	holeFirstIP, _ := prevIPBlock.NextIP()
	holeEndIP, _ := currIPBlock.PreviousIP()
	hole, _ := netset.IPBlockFromIPRange(holeFirstIP, holeEndIP)
	return !hole.IsSubset(anyProtocolCubes)
}
