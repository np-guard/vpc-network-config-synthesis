/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"slices"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func aclCubesToRules(cubes *aclCubesPerProtocol, direction ir.Direction) []*ir.ACLRule {
	tcpRules := triplesToRulesTCPUDP(cubes.tcpAllow, direction, true)
	udpRules := triplesToRulesTCPUDP(cubes.udpAllow, direction, false)
	icmpRules := triplesToRulesICMP(cubes.icmpAllow, direction)
	return slices.Concat(tcpRules, udpRules, icmpRules)
}

func triplesToRulesTCPUDP(tripleSet tcpudpTripleSet, direction ir.Direction, isTCP bool) []*ir.ACLRule {
	partitions := minimalPartitionsTCPUDP(tripleSet, isTCP)
	res := make([]*ir.ACLRule, len(partitions))
	for i, t := range partitions {
		res[i] = ir.NewACLRule(ir.Allow, direction, t.S1, t.S2, t.S3, "")
	}
	return res
}

func minimalPartitionsTCPUDP(t tcpudpTripleSet, isTCP bool) []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP] {
	leftPartitions := actualPartitionsTCPUDP(ds.AsLeftTripleSet(t), isTCP)
	outerPartitions := actualPartitionsTCPUDP(ds.AsOuterTripleSet(t), isTCP)
	rightPartitions := actualPartitionsTCPUDP(ds.AsRightTripleSet(t), isTCP)

	switch {
	case len(leftPartitions) <= len(outerPartitions) && len(leftPartitions) <= len(rightPartitions):
		return leftPartitions
	case len(outerPartitions) <= len(leftPartitions) && len(outerPartitions) <= len(rightPartitions):
		return outerPartitions
	default:
		return rightPartitions
	}
}

func actualPartitionsTCPUDP(t tcpudpTripleSet, isTCP bool) []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP], 0)
	for _, p := range t.Partitions() {
		res = append(res, breakTCPUDPTriple(p, isTCP)...)
	}
	return res
}

// break multi-cube to regular cube
func breakTCPUDPTriple(t ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.PortSet], isTCP bool) []ds.Triple[*netset.IPBlock,
	*netset.IPBlock, netp.TCPUDP] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP], 0)

	dstCidrs := t.S2.SplitToCidrs()
	portIntervals := t.S3.Intervals()

	for _, src := range t.S1.SplitToCidrs() {
		for _, dst := range dstCidrs {
			for _, ports := range portIntervals {
				protocol, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(ports.Start()), int(ports.End()))
				res = append(res, ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP]{S1: src, S2: dst, S3: protocol})
			}
		}
	}
	return res
}

func triplesToRulesICMP(tripleSet icmpTripleSet, direction ir.Direction) []*ir.ACLRule {
	partitions := minimalPartitionsICMP(tripleSet)
	res := make([]*ir.ACLRule, len(partitions))
	for i, t := range partitions {
		res[i] = ir.NewACLRule(ir.Allow, direction, t.S1, t.S2, t.S3, "")
	}
	return res
}

func minimalPartitionsICMP(t icmpTripleSet) []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.ICMP] {
	leftPartitions := actualPartitionsICMP(ds.AsLeftTripleSet(t))
	outerPartitions := actualPartitionsICMP(ds.AsOuterTripleSet(t))
	rightPartitions := actualPartitionsICMP(ds.AsRightTripleSet(t))

	switch {
	case len(leftPartitions) <= len(outerPartitions) && len(leftPartitions) <= len(rightPartitions):
		return leftPartitions
	case len(outerPartitions) <= len(leftPartitions) && len(outerPartitions) <= len(rightPartitions):
		return outerPartitions
	default:
		return rightPartitions
	}
}

func actualPartitionsICMP(t icmpTripleSet) []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.ICMP] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.ICMP], 0)
	for _, p := range t.Partitions() {
		res = append(res, breakICMPTriple(p)...)
	}
	return res
}

// break multi-cube to regular cube
func breakICMPTriple(t ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet]) []ds.Triple[*netset.IPBlock,
	*netset.IPBlock, netp.ICMP] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.ICMP], 0)

	dstCidrs := t.S2.SplitToCidrs()
	icmpPartitions := optimize.IcmpsetPartitions(t.S3)

	for _, src := range t.S1.SplitToCidrs() {
		for _, dst := range dstCidrs {
			for _, icmp := range icmpPartitions {
				a := ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.ICMP]{S1: src, S2: dst, S3: icmp}
				res = append(res, a)
			}
		}
	}
	return res
}
