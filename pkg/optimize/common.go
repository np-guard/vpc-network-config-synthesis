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

func IcmpToIcmpSet(icmp netp.ICMP) *netset.ICMPSet {
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

func IsAllPorts(tcpudpPorts *netset.PortSet) bool {
	return tcpudpPorts.Equal(netset.AllPorts())
}

func MinimalPartitions[P ds.Set[P]](t ds.TripleSet[*netset.IPBlock, *netset.IPBlock, P]) []ds.Triple[*netset.IPBlock, *netset.IPBlock, P] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, P], 0)

	//leftTripleSet := AsLeftTripleSet(t)
	//leftTripleSetPartitions :=

	return res
}

// func AsLeftTripleSet[P ds.Set[P]](t ds.TripleSet[*netset.IPBlock, *netset.IPBlock, P]) ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, P]{

// }

// func ActualPartitions[P ds.Set[P]](t ds.TripleSet[*netset.IPBlock, *netset.IPBlock, P]) []ds.Triple[*netset.IPBlock, *netset.IPBlock, P] {
// 	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, P], 0)
// 	for _, p := range t.Partitions() {
// 		if tcpudp, ok := p.S3.(netp.TCPUDP); ok {
// 			res = append(res, DecomposeTCPUDPTriple(p)...)
// 		} else {
// 			res = append(res, DecomposeICMPTriple(p)...)
// 		}
// 	}
// 	return res
// }

func DecomposeTCPUDPTriple(t ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.PortSet]) []ds.Triple[*netset.IPBlock,
	*netset.IPBlock, *netset.PortSet] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.PortSet], 0)

	dstCidrs := t.S2.SplitToCidrs()
	portIntervals := t.S3.Intervals()

	for _, src := range t.S1.SplitToCidrs() {
		for _, dst := range dstCidrs {
			for _, ports := range portIntervals {
				a := ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.PortSet]{S1: src, S2: dst, S3: ports.ToSet()}
				res = append(res, a)
			}
		}
	}
	return res
}

func DecomposeICMPTriple(t ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet]) []ds.Triple[*netset.IPBlock,
	*netset.IPBlock, *netset.ICMPSet] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet], 0)

	dstCidrs := t.S2.SplitToCidrs()
	icmpPartitions := IcmpsetPartitions(t.S3)

	for _, src := range t.S1.SplitToCidrs() {
		for _, dst := range dstCidrs {
			for _, icmp := range icmpPartitions {
				a := ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet]{S1: src, S2: dst, S3: IcmpToIcmpSet(icmp)}
				res = append(res, a)
			}
		}
	}
	return res
}
