/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"slices"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func aclCubesToRules(cubes *aclCubesPerProtocol, direction ir.Direction) []*ir.ACLRule {
	reduceACLCubes(cubes)
	anyProtocolRules := anyProtocolCubesToRules(cubes.anyProtocolAllow, direction)
	tcpudpRules := tcpudpTriplesToRules(cubes.tcpudpAllow, direction)
	icmpRules := icmpTriplesToRules(cubes.icmpAllow, direction)
	return slices.Concat(anyProtocolRules, tcpudpRules, icmpRules)
}

// same algorithm as sg cubes to rules
func anyProtocolCubesToRules(cubes srcDstProduct, direction ir.Direction) []*ir.ACLRule {
	partitions := optimize.SortPartitionsByIPAddrs(cubes.Partitions())
	if len(partitions) == 0 {
		return []*ir.ACLRule{}
	}

	res := make([]*ir.ACLRule, 0)
	activeRules := make(map[*netset.IPBlock]*netset.IPBlock) // key = first src's IP, value = dst IP

	for i := range partitions {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && uncoveredHole(partitions[i].Left, partitions[i].Left) {
			res = slices.Concat(res, createActiveRules(activeRules, partitions[i-1].Left.LastIPAddressObject(), direction))
			activeRules = make(map[*netset.IPBlock]*netset.IPBlock)
		}

		// if there are active rules whose dsts are not fully included in the current cube, they will be created
		// also activeDstIPs will be calculated, which is the dstIPs that are still included in the active rules
		activeDstIPs := netset.NewIPBlock()
		for srcFirstIP, dstCidr := range activeRules {
			if dstCidr.IsSubset(partitions[i].Right) {
				activeDstIPs = activeDstIPs.Union(dstCidr)
			} else {
				res = createNewRules(srcFirstIP, partitions[i-1].Left.LastIPAddressObject(), dstCidr, direction) // create active rule
				delete(activeRules, srcFirstIP)
			}
		}

		// if the current cube contains dstIPs that are not contained in active rules, new rules will be created
		for _, dstCidr := range partitions[i].Right.SplitToCidrs() {
			if !dstCidr.IsSubset(activeDstIPs) {
				activeRules[partitions[i].Left.FirstIPAddressObject()] = dstCidr
			}
		}
	}
	// generate all existing rules
	return slices.Concat(res, createActiveRules(activeRules, partitions[len(partitions)-1].Left.LastIPAddressObject(), direction))
}

func createActiveRules(activeRules map[*netset.IPBlock]*netset.IPBlock, srcLastIP *netset.IPBlock, direction ir.Direction) []*ir.ACLRule {
	res := make([]*ir.ACLRule, 0)
	for srcFirstIP, dstCidr := range activeRules {
		res = slices.Concat(res, createNewRules(srcFirstIP, srcLastIP, dstCidr, direction))
	}
	return res
}

func createNewRules(srcStartIP, srcEndIP, dstCidr *netset.IPBlock, direction ir.Direction) []*ir.ACLRule {
	src, _ := netset.IPBlockFromIPRange(srcStartIP, srcEndIP)
	srcCidrs := src.SplitToCidrs()

	res := make([]*ir.ACLRule, len(srcCidrs))
	for i, srcCidr := range srcCidrs {
		res[i] = ir.NewACLRule(ir.Allow, direction, srcCidr, dstCidr, netp.AnyProtocol{}, "")
	}
	return res
}

func uncoveredHole(prevSrcIP, currSrcIP *netset.IPBlock) bool {
	touching, _ := prevSrcIP.TouchingIPRanges(currSrcIP)
	return touching
}
