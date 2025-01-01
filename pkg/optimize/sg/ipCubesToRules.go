/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"slices"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

// any protocol IP-segments, represented by a single ipblock that will be decomposed
// into cidrs. Each cidr will be a remote of a single SG rule
func anyProtocolIPCubesToRules(cubes *netset.IPBlock, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	result := make([]*ir.SGRule, 0)
	for _, cidr := range cubes.SplitToCidrs() {
		result = append(result, ir.NewSGRule(direction, cidr, netp.AnyProtocol{}, l, ""))
	}
	return result
}

// tcpudpIPCubesToRules converts cubes representing tcp or udp protocol rules to SG rules
func tcpudpIPCubesToRules(cubes []ds.Pair[*netset.IPBlock, *netset.PortSet], anyProtocolCubes *netset.IPBlock, direction ir.Direction,
	isTCP bool, l *netset.IPBlock) []*ir.SGRule {
	if len(cubes) == 0 {
		return []*ir.SGRule{}
	}

	res := make([]*ir.SGRule, 0)
	activeRules := make([]ds.Triple[*netset.IPBlock, netp.Protocol, int], 0) // first IP of the rule; protocol; in which cube we added the rule

	for i := range cubes {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && uncoveredHole(cubes[i-1], cubes[i], anyProtocolCubes) {
			res = slices.Concat(res, createActiveRules(activeRules, cubes[i-1].Left.LastIPAddressObject(), direction, l))
			activeRules = make([]ds.Triple[*netset.IPBlock, netp.Protocol, int], 0)
		}

		// if there are active rules whose ports are not fully included in the current cube, they will be created
		// also activePorts will be calculated, which is the ports that are still included in the active rules
		activePorts := interval.NewCanonicalSet()
		for j := len(activeRules) - 1; j >= 0; j-- {
			rule := activeRules[j]
			tcpudpPorts := rule.S2.(netp.TCPUDP).DstPorts() // already checked
			if !tcpudpPorts.ToSet().IsSubset(cubes[i].Right) {
				if !redundantTCPUDPRule(cubes, activeRules, j, i) {
					res = slices.Concat(res, createNewRules(rule.S2, rule.S1, cubes[i-1].Left.LastIPAddressObject(), direction, l))
				}
				activeRules = slices.Concat(activeRules[0:j], activeRules[j+1:])
			} else {
				activePorts.AddInterval(tcpudpPorts)
			}
		}

		// if the current cube contains ports that are not contained in active rules, new rules will be created
		for _, ports := range cubes[i].Right.Intervals() {
			if !ports.ToSet().IsSubset(activePorts) {
				p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(ports.Start()), int(ports.End()))
				rule := ds.Triple[*netset.IPBlock, netp.Protocol, int]{S1: cubes[i].Left.FirstIPAddressObject(), S2: p, S3: i}
				activeRules = append(activeRules, rule)
			}
		}
	}
	// generate all existing rules
	return slices.Concat(res, createActiveRules(activeRules, cubes[len(cubes)-1].Left.LastIPAddressObject(), direction, l))
}

// icmpIPCubesToRules converts cubes representing icmp protocol rules to SG rules
func icmpIPCubesToRules(cubes []ds.Pair[*netset.IPBlock, *netset.ICMPSet], anyProtocolCubes *netset.IPBlock, direction ir.Direction,
	l *netset.IPBlock) []*ir.SGRule {
	if len(cubes) == 0 {
		return []*ir.SGRule{}
	}

	res := make([]*ir.SGRule, 0)
	activeRules := make([]ds.Triple[*netset.IPBlock, netp.Protocol, int], 0) // first IP of the rule; protocol; in which cube we added the rule

	for i := range cubes {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && uncoveredHole(cubes[i-1], cubes[i], anyProtocolCubes) {
			res = slices.Concat(res, createActiveRules(activeRules, cubes[i-1].Left.LastIPAddressObject(), direction, l))
			activeRules = make([]ds.Triple[*netset.IPBlock, netp.Protocol, int], 0)
		}

		// if there are active rules whose icmp values are not fully included in the current cube, they will be created
		// also activeICMP will be calculated, which is the icmp values that are still included in the active rules
		activeICMP := netset.EmptyICMPSet()
		for j := len(activeRules) - 1; j >= 0; j-- {
			rule := activeRules[j]
			icmpSet := optimize.IcmpToIcmpSet(rule.S2.(netp.ICMP))

			if !icmpSet.IsSubset(cubes[i].Right) {
				if !redundantICMPRule(cubes, activeRules, j, i) {
					res = slices.Concat(res, createNewRules(rule.S2, rule.S1, cubes[i-1].Left.LastIPAddressObject(), direction, l))
				}
				activeRules = slices.Concat(activeRules[0:j], activeRules[j+1:])
			} else {
				activeICMP.Union(icmpSet)
			}
		}

		// if the cube contains icmp values that are not contained in  active rules, new rules will be created
		for _, p := range optimize.IcmpsetPartitions(cubes[i].Right) {
			if !optimize.IcmpToIcmpSet(p).IsSubset(activeICMP) {
				rule := ds.Triple[*netset.IPBlock, netp.Protocol, int]{S1: cubes[i].Left.FirstIPAddressObject(), S2: p, S3: i}
				activeRules = append(activeRules, rule)
			}
		}
	}

	// generate all  existing rules
	return slices.Concat(res, createActiveRules(activeRules, cubes[len(cubes)-1].Left.LastIPAddressObject(), direction, l))
}

// uncoveredHole returns true if the rules can not be continued between the two cubes
// i.e there is a hole between two ipblocks that is not a subset of anyProtocol cubes
func uncoveredHole[T ds.Set[T]](prevPair, currPair ds.Pair[*netset.IPBlock, T], anyProtocolCubes *netset.IPBlock) bool {
	prevIPBlock := prevPair.Left
	currIPBlock := currPair.Left
	touching, _ := prevIPBlock.TouchingIPRanges(currIPBlock)
	if touching {
		return false
	}
	holeFirstIP, _ := prevIPBlock.NextIP()
	holeEndIP, _ := currIPBlock.PreviousIP()
	hole, _ := netset.IPBlockFromIPRange(holeFirstIP, holeEndIP)
	return !hole.IsSubset(anyProtocolCubes)
}

// creates sgRules from SG active rules
func createActiveRules(activeRules []ds.Triple[*netset.IPBlock, netp.Protocol, int], lastIP *netset.IPBlock,
	direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	res := make([]*ir.SGRule, 0)
	for _, triple := range activeRules {
		res = slices.Concat(res, createNewRules(triple.S2, triple.S1, lastIP, direction, l))
	}
	return res
}

// createNewRules breaks the startIP-endIP ip range into cidrs and creates SG rules
func createNewRules(protocol netp.Protocol, startIP, endIP *netset.IPBlock, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	res := make([]*ir.SGRule, 0)
	ipRange, _ := netset.IPBlockFromIPRange(startIP, endIP)
	for _, cidr := range ipRange.SplitToCidrs() {
		res = append(res, ir.NewSGRule(direction, cidr, protocol, l, ""))
	}
	return res
}

func redundantTCPUDPRule(cubes []ds.Pair[*netset.IPBlock, *netset.PortSet], activeRules []ds.Triple[*netset.IPBlock, netp.Protocol, int],
	ruleIdx, currIdx int) bool {
	rule := activeRules[ruleIdx]
	for j := rule.S3; j < currIdx; j++ {
		if !cubes[currIdx].Left.IsSubset(cubes[j].Left) {
			return false
		}
	}

	uncoveredPorts := rule.S2.(netp.TCPUDP).DstPorts().ToSet()
	for i := 0; i < ruleIdx; i++ {
		tcpudpPorts := activeRules[ruleIdx].S2.(netp.TCPUDP).DstPorts().ToSet()
		uncoveredPorts = uncoveredPorts.Subtract(tcpudpPorts)
	}
	return uncoveredPorts.IsEmpty()
}

func redundantICMPRule(cubes []ds.Pair[*netset.IPBlock, *netset.ICMPSet], activeRules []ds.Triple[*netset.IPBlock, netp.Protocol, int],
	ruleIdx, currIdx int) bool {
	rule := activeRules[ruleIdx]
	for j := rule.S3; j < currIdx; j++ {
		if !cubes[currIdx].Left.IsSubset(cubes[j].Left) {
			return false
		}
	}

	uncoveredICMP := optimize.IcmpToIcmpSet(rule.S2.(netp.ICMP))
	for i := 0; i < ruleIdx; i++ {
		icmp := optimize.IcmpToIcmpSet(activeRules[ruleIdx].S2.(netp.ICMP))
		uncoveredICMP = uncoveredICMP.Subtract(icmp)
	}
	return uncoveredICMP.IsEmpty()
}
