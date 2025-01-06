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

// Creates two sets of rules: one with only icmp protocol, and the other icmp protocol combined
// with any protocol cubes. It returns the minimal set
func icmpRules(icmpCubes icmpTripleSet, anyCubes srcDstProduct, direction ir.Direction) []*ir.ACLRule {
	res := icmpTriplesToRules(icmpCubes, direction, ir.Allow)
	resWithAny := icmpTriplesToRules(addSrcDstCubeToICMP(icmpCubes, anyCubes, netset.AllICMPSet()), direction, ir.Allow)
	if len(resWithAny) < len(res) {
		return resWithAny
	}
	return res
}

func icmpTriplesToRules(tripleSet icmpTripleSet, direction ir.Direction, action ir.Action) []*ir.ACLRule {
	partitions := minimalPartitionsICMP(tripleSet)
	res := make([]*ir.ACLRule, len(partitions))
	for i, t := range partitions {
		res[i] = ir.NewACLRule(action, direction, t.S1, t.S2, t.S3, "")
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
		res = slices.Concat(res, breakICMPTriple(p))
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
