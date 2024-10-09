/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func allProtocolIPCubesIPToRules(cubes *netset.IPBlock, direction ir.Direction) []*ir.SGRule {
	result := make([]*ir.SGRule, 0)
	for _, cidr := range cubes.SplitToCidrs() {
		result = append(result, ir.NewSGRule(direction, cidr, netp.AnyProtocol{}, netset.GetCidrAll(), ""))
	}
	return result
}

func tcpudpIPCubesToRules(cubes []ds.Pair[*netset.IPBlock, *netset.PortSet], allCubes *netset.IPBlock,
	direction ir.Direction, isTCP bool) []*ir.SGRule {
	activeRules := make(map[*netset.IPBlock]*interval.Interval) // start ip and ports
	result := make([]*ir.SGRule, 0)

	for i := range cubes {
		if i > 0 && !continuation(cubes[i-1], cubes[i], allCubes) {
			for ipb, ports := range activeRules {
				p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(ports.Start()), int(ports.End()))
				ipRange, _ := netset.IPBlockFromIPRange(ipb, cubes[i-1].Left.LastIPAddressObject())
				for _, cidr := range ipRange.SplitToCidrs() {
					result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
				}
				delete(activeRules, ipb)
			}
		}

		// rules whose ports are not in the current cube will not remain active
		activePorts := interval.NewCanonicalSet()
		for ipb, ports := range activeRules {
			if !ports.ToSet().IsSubset(cubes[i].Right) {
				p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(ports.Start()), int(ports.End()))
				ipRange, _ := netset.IPBlockFromIPRange(ipb, cubes[i-1].Left.LastIPAddressObject())
				for _, cidr := range ipRange.SplitToCidrs() {
					result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
				}
			} else {
				activePorts.AddInterval(*ports)
			}
		}

		// if the cube contains ports that are not contained in active rules, new rules will be created
		for _, ports := range cubes[i].Right.Intervals() {
			if !ports.ToSet().IsSubset(activePorts) {
				activeRules[cubes[i].Left.FirstIPAddressObject()] = &ports
			}
		}
	}

	// create the rest of the rules
	for ipb, ports := range activeRules {
		p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(ports.Start()), int(ports.End()))
		ipRange, _ := netset.IPBlockFromIPRange(ipb, cubes[len(cubes)-1].Left.LastIPAddressObject())
		for _, cidr := range ipRange.SplitToCidrs() {
			result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
		}
	}

	return result
}

// problem: where should I end the rule?
// func createTcpudpRules(rules []ds.Pair[*netset.IPBlock, *interval.Interval], span []ds.Pair[*netset.IPBlock, *netset.PortSet],
// 	direction ir.Direction, isTCP bool) (res []ir.SGRule) {
// 	res = make([]ir.SGRule, 0)
// 	for _, r := range rules {
// 		p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(r.Right.Start()), int(r.Right.Start()))
// 		for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[len(span)-1].Left))) {
// 			res = append(res, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
// 		}
// 	}
// 	return res
// }

func icmpIPCubesToRules(cubes []ds.Pair[*netset.IPBlock, *netset.ICMPSet], allCubes *netset.IPBlock, direction ir.Direction) []*ir.SGRule {
	activeRules := make(map[*netset.IPBlock]*netp.ICMP)
	result := make([]*ir.SGRule, 0)

	for i := range cubes {
		if i > 0 && !continuation(cubes[i-1], cubes[i], allCubes) {
			for ipb, icmp := range activeRules {
				p, _ := netp.NewICMP(icmp.TypeCode)
				ipRange, _ := netset.IPBlockFromIPRange(ipb, cubes[i-1].Left.LastIPAddressObject())
				for _, cidr := range ipRange.SplitToCidrs() {
					result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
				}
				delete(activeRules, ipb)
			}
		}

		// rules whose icmp value is not in the current cube will not remain active
		activeICMP := netset.EmptyICMPSet()
		for ipb, icmp := range activeRules {
			ruleIcmpSet := optimize.IcmpRuleToIcmpSet(*icmp)
			if !ruleIcmpSet.IsSubset(cubes[i].Right) { // create rules
				p, _ := netp.NewICMPWithoutRFCValidation(icmp.TypeCode)
				ipRange, _ := netset.IPBlockFromIPRange(ipb, cubes[i-1].Left.LastIPAddressObject())
				for _, cidr := range ipRange.SplitToCidrs() {
					result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
				}
			} else {
				activeICMP.Union(ruleIcmpSet)
			}
		}

		// new rules
		for _, p := range optimize.IcmpsetPartitions(cubes[i].Right) {
			if !optimize.IcmpRuleToIcmpSet(p).IsSubset(activeICMP) {
				activeRules[cubes[i].Left.FirstIPAddressObject()] = &p
			}
		}
	}

	// create the rest of the rules
	for ipb, icmp := range activeRules {
		p, _ := netp.NewICMP(icmp.TypeCode)
		ipRange, _ := netset.IPBlockFromIPRange(ipb, cubes[len(cubes)-1].Left.LastIPAddressObject())
		for _, cidr := range ipRange.SplitToCidrs() {
			result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
		}
	}

	return result
}

func continuation[T ds.Set[T]](prevPair, currPair ds.Pair[*netset.IPBlock, T], allProtocolCubes *netset.IPBlock) bool {
	prevIPBlock := prevPair.Left
	currIPBlock := currPair.Left
	touching, _ := prevIPBlock.TouchingIPRanges(currIPBlock)
	if touching {
		return true
	}
	startH, _ := prevIPBlock.NextIP()
	endH, _ := currIPBlock.PreviousIP()
	hole, _ := netset.IPBlockFromIPRange(startH, endH)
	return hole.IsSubset(allProtocolCubes)
}
