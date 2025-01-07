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

// same algorithm as sg cubes to rules
func anyProtocolCubesToRules(cubes srcDstProduct, direction ir.Direction, action ir.Action) []*ir.ACLRule {
	partitions := optimize.SortPartitionsByIPAddrs(cubes.Partitions())
	if len(partitions) == 0 {
		return []*ir.ACLRule{}
	}

	res := make([]*ir.ACLRule, 0)
	activeRules := make([]ds.Pair[*netset.IPBlock, *netset.IPBlock], 0) // Left = first src's IP, Right = dst cidr

	for i := range partitions {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && uncoveredHoleAny(partitions[i-1].Left, partitions[i].Left) {
			res = slices.Concat(res, createAnyActiveRules(activeRules, partitions[i-1].Left.LastIPAddressObject(), direction, action))
			activeRules = make([]ds.Pair[*netset.IPBlock, *netset.IPBlock], 0)
		}

		// if there are active rules whose dsts are not fully included in the current cube, they will be created
		// also activeDstIPs will be calculated, which is the dstIPs that are still included in the active rules
		activeDstIPs := netset.NewIPBlock()
		for j, rule := range slices.Backward(activeRules) {
			if rule.Right.IsSubset(partitions[i].Right) {
				activeDstIPs = activeDstIPs.Union(rule.Right)
			} else {
				res = slices.Concat(res, createNewAnyRules(rule.Left, partitions[i-1].Left.LastIPAddressObject(), rule.Right, direction, action))
				activeRules = slices.Delete(activeRules, j, j+1)
			}
		}

		// if the current cube contains dstIPs that are not contained in active rules, new rules will be created
		for _, dstCidr := range partitions[i].Right.SplitToCidrs() {
			if !dstCidr.IsSubset(activeDstIPs) {
				rule := ds.Pair[*netset.IPBlock, *netset.IPBlock]{Left: partitions[i].Left.FirstIPAddressObject(), Right: dstCidr}
				activeRules = append(activeRules, rule)
			}
		}
	}
	// generate all existing rules
	return slices.Concat(res, createAnyActiveRules(activeRules, partitions[len(partitions)-1].Left.LastIPAddressObject(), direction, action))
}

func createAnyActiveRules(activeRules []ds.Pair[*netset.IPBlock, *netset.IPBlock], srcLastIP *netset.IPBlock,
	direction ir.Direction, action ir.Action) []*ir.ACLRule {
	res := make([]*ir.ACLRule, 0)
	for _, rule := range activeRules {
		res = slices.Concat(res, createNewAnyRules(rule.Left, srcLastIP, rule.Right, direction, action))
	}
	return res
}

func createNewAnyRules(srcStartIP, srcEndIP, dstCidr *netset.IPBlock, direction ir.Direction, action ir.Action) []*ir.ACLRule {
	src, _ := netset.IPBlockFromIPRange(srcStartIP, srcEndIP)
	srcCidrs := src.SplitToCidrs()

	res := make([]*ir.ACLRule, len(srcCidrs))
	for i, srcCidr := range srcCidrs {
		res[i] = ir.NewACLRule(action, direction, srcCidr, dstCidr, netp.AnyProtocol{}, "")
	}
	return res
}

func uncoveredHoleAny(prevSrcIP, currSrcIP *netset.IPBlock) bool {
	touching, _ := prevSrcIP.TouchingIPRanges(currSrcIP)
	return !touching
}
