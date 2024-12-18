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
	activeRules := make(map[*netset.IPBlock]netp.Protocol) // the key is the first IP

	for i := range cubes {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && uncoveredHole(cubes[i-1], cubes[i], anyProtocolCubes) {
			res = slices.Concat(res, createActiveRules(activeRules, cubes[i-1].Left.LastIPAddressObject(), direction, l))
			activeRules = make(map[*netset.IPBlock]netp.Protocol)
		}

		// if there are active rules whose ports are not fully included in the current cube, they will be created
		// also activePorts will be calculated, which is the ports that are still included in the active rules
		activePorts := interval.NewCanonicalSet()
		for startIP, protocol := range activeRules {
			tcpudp, _ := protocol.(netp.TCPUDP) // already checked
			if !tcpudp.DstPorts().ToSet().IsSubset(cubes[i].Right) {
				res = slices.Concat(res, createNewRules(protocol, startIP, cubes[i-1].Left.LastIPAddressObject(), direction, l))
				delete(activeRules, startIP)
			} else {
				activePorts.AddInterval(tcpudp.DstPorts())
			}
		}

		// if the current cube contains ports that are not contained in active rules, new rules will be created
		for _, ports := range cubes[i].Right.Intervals() {
			if !ports.ToSet().IsSubset(activePorts) {
				p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(ports.Start()), int(ports.End()))
				activeRules[cubes[i].Left.FirstIPAddressObject()] = p
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
	activeRules := make(map[*netset.IPBlock]netp.Protocol) // the key is the first IP

	for i := range cubes {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && uncoveredHole(cubes[i-1], cubes[i], anyProtocolCubes) {
			res = slices.Concat(res, createActiveRules(activeRules, cubes[i-1].Left.LastIPAddressObject(), direction, l))
			activeRules = make(map[*netset.IPBlock]netp.Protocol)
		}

		// if there are active rules whose icmp values are not fully included in the current cube, they will be created
		// also activeICMP will be calculated, which is the icmp values that are still included in the active rules
		activeICMP := netset.EmptyICMPSet()
		for startIP, protocol := range activeRules {
			icmp, _ := protocol.(netp.ICMP)
			ruleIcmpSet := optimize.IcmpRuleToIcmpSet(icmp)
			if !ruleIcmpSet.IsSubset(cubes[i].Right) {
				res = slices.Concat(res, createNewRules(protocol, startIP, cubes[i-1].Left.LastIPAddressObject(), direction, l))
				delete(activeRules, startIP)
			} else {
				activeICMP.Union(ruleIcmpSet)
			}
		}

		// if the cube contains icmp values that are not contained in  active rules, new rules will be created
		for _, p := range optimize.IcmpsetPartitions(cubes[i].Right) {
			if !optimize.IcmpRuleToIcmpSet(p).IsSubset(activeICMP) {
				activeRules[cubes[i].Left.FirstIPAddressObject()] = p
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
func createActiveRules(activeRules map[*netset.IPBlock]netp.Protocol, lastIP *netset.IPBlock,
	direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	res := make([]*ir.SGRule, 0)
	for firstIP, protocol := range activeRules {
		res = slices.Concat(res, createNewRules(protocol, firstIP, lastIP, direction, l))
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
