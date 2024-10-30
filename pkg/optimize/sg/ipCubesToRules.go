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

// all protocol cubes, represented by a single ipblock that will be decomposed
// into cidrs. Each cidr will be the remote of a SG rule
func allProtocolIPCubesIPToRules(cubes *netset.IPBlock, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	result := make([]*ir.SGRule, 0)
	for _, cidr := range cubes.SplitToCidrs() {
		result = append(result, ir.NewSGRule(direction, cidr, netp.AnyProtocol{}, l, ""))
	}
	return result
}

// tcpudpIPCubesToRules converts cubes representing tcp or udp protocol rules to SG rules
func tcpudpIPCubesToRules(cubes []ds.Pair[*netset.IPBlock, *netset.PortSet], allCubes *netset.IPBlock, direction ir.Direction,
	isTCP bool, l *netset.IPBlock) []*ir.SGRule {
	if len(cubes) == 0 {
		return []*ir.SGRule{}
	}

	activeRules := make(map[*netset.IPBlock]netp.Protocol) // the key is the first IP
	result := make([]*ir.SGRule, 0)

	for i := range cubes {
		// if it is not possible to continue the rule between the cubes, generate all the existing rules
		if i > 0 && !continuation(cubes[i-1], cubes[i], allCubes) {
			result = append(result, createActiveRules(activeRules, cubes[i-1].Left.LastIPAddressObject(), direction, l)...)
			activeRules = make(map[*netset.IPBlock]netp.Protocol)
		}

		// if the proctol is not contained in the current cube, we will generate SG rules
		// calculate active ports = active rules covers these ports
		activePorts := interval.NewCanonicalSet()
		for ipb, protocol := range activeRules {
			if tcpudp, ok := protocol.(netp.TCPUDP); ok {
				if !tcpudp.DstPorts().ToSet().IsSubset(cubes[i].Right) {
					result = append(result, createNewRules(protocol, ipb, cubes[i-1].Left.LastIPAddressObject(), direction, l)...)
				} else {
					activePorts.AddInterval(tcpudp.DstPorts())
				}
			}
		}

		// if the cube contains ports that are not contained in the active rules, new rules will be created
		for _, ports := range cubes[i].Right.Intervals() {
			if !ports.ToSet().IsSubset(activePorts) {
				p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(ports.Start()), int(ports.End()))
				activeRules[cubes[i].Left.FirstIPAddressObject()] = p
			}
		}
	}
	// generate all the existing rules
	return append(result, createActiveRules(activeRules, cubes[len(cubes)-1].Left.LastIPAddressObject(), direction, l)...)
}

// icmpIPCubesToRules converts cubes representing icmp protocol rules to SG rules
func icmpIPCubesToRules(cubes []ds.Pair[*netset.IPBlock, *netset.ICMPSet], allCubes *netset.IPBlock, direction ir.Direction,
	l *netset.IPBlock) []*ir.SGRule {
	if len(cubes) == 0 {
		return []*ir.SGRule{}
	}

	activeRules := make(map[*netset.IPBlock]netp.Protocol) // the key is the first IP
	result := make([]*ir.SGRule, 0)

	for i := range cubes {
		// if it is not possible to continue the rule between the cubes, generate all the existing rules
		if i > 0 && !continuation(cubes[i-1], cubes[i], allCubes) {
			result = append(result, createActiveRules(activeRules, cubes[i-1].Left.LastIPAddressObject(), direction, l)...)
			activeRules = make(map[*netset.IPBlock]netp.Protocol)
		}

		// if the proctol is not contained in the current cube, we will generate SG rules
		// calculate activeICMP = active rules covers these icmp values
		activeICMP := netset.EmptyICMPSet()
		for ipb, protocol := range activeRules {
			if icmp, ok := protocol.(netp.ICMP); ok {
				ruleIcmpSet := optimize.IcmpRuleToIcmpSet(icmp)
				if !ruleIcmpSet.IsSubset(cubes[i].Right) {
					result = append(result, createNewRules(protocol, ipb, cubes[i-1].Left.LastIPAddressObject(), direction, l)...)
				} else {
					activeICMP.Union(ruleIcmpSet)
				}
			}
		}

		// if the cube contains icmp values that are not contained in the active rules, new rules will be created
		for _, p := range optimize.IcmpsetPartitions(cubes[i].Right) {
			if !optimize.IcmpRuleToIcmpSet(p).IsSubset(activeICMP) {
				activeRules[cubes[i].Left.FirstIPAddressObject()] = p
			}
		}
	}

	// generate all the existing rules
	return append(result, createActiveRules(activeRules, cubes[len(cubes)-1].Left.LastIPAddressObject(), direction, l)...)
}

// continuation returns true if the rules can be continued between the two cubes
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

// creates sgRules from SG active rules
func createActiveRules(activeRules map[*netset.IPBlock]netp.Protocol, endIP *netset.IPBlock, direction ir.Direction,
	l *netset.IPBlock) []*ir.SGRule {
	res := make([]*ir.SGRule, 0)
	for ipb, protocol := range activeRules {
		res = append(res, createNewRules(protocol, ipb, endIP, direction, l)...)
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
