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

// tcpudpIPCubesToMinRules calls tcpudpIPCubesToRules func twice -- once when the any protocol cubes are converted
// to tcp/udp cubes and onces when they are not. it returns the best result (less sg rules)
func tcpudpIPCubesToMinRules(cubes []ds.Pair[*netset.IPBlock, *netset.PortSet], anyProtocolCubes *netset.IPBlock, direction ir.Direction,
	isTCP bool, l *netset.IPBlock) []*ir.SGRule {
	res := tcpudpIPCubesToRules(cubes, anyProtocolCubes, direction, isTCP, l)
	anyAsTCPUDP := partitionsToProduct(cubes).Union(ds.CartesianPairLeft(anyProtocolCubes, netset.AllPorts()))
	resWithAny := tcpudpIPCubesToRules(optimize.SortPartitionsByIPAddrs(anyAsTCPUDP.Partitions()),
		netset.NewIPBlock(), direction, isTCP, l) // pass an empty ipblock instead of anyProtocol iblock
	if len(resWithAny) < len(res) {
		return resWithAny
	}
	return res
}

// tcpudpIPCubesToRules converts cubes representing tcp or udp protocol rules to SG rules
func tcpudpIPCubesToRules(cubes []ds.Pair[*netset.IPBlock, *netset.PortSet], anyProtocolCubes *netset.IPBlock, direction ir.Direction,
	isTCP bool, l *netset.IPBlock) []*ir.SGRule {
	if len(cubes) == 0 {
		return []*ir.SGRule{}
	}

	res := make([]*ir.SGRule, 0)
	activeRules := make([]ds.Pair[*netset.IPBlock, netp.Protocol], 0) // first IP of the rule; protocol

	for i := range cubes {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && optimize.UncoveredHole(cubes[i-1].Left, cubes[i].Left, anyProtocolCubes) {
			res = slices.Concat(res, createActiveRules(activeRules, cubes[i-1].Left.LastIPAddressObject(), direction, l))
			activeRules = make([]ds.Pair[*netset.IPBlock, netp.Protocol], 0)
		}

		// if there are active rules whose ports are not fully included in the current cube, they will be created
		// also activePorts will be calculated, which is the ports that are still included in the active rules
		activePorts := interval.NewCanonicalSet()
		for j, rule := range slices.Backward(activeRules) {
			tcpudpPorts := rule.Right.(netp.TCPUDP).DstPorts() // already checked
			if !tcpudpPorts.ToSet().IsSubset(cubes[i].Right) {
				res = slices.Concat(res, createNewRules(rule.Right, rule.Left, cubes[i-1].Left.LastIPAddressObject(), direction, l))
				activeRules = slices.Delete(activeRules, j, j+1)
			} else {
				activePorts.AddInterval(tcpudpPorts)
			}
		}

		// if the current cube contains ports that are not contained in active rules, new rules will be created
		for _, ports := range cubes[i].Right.Intervals() {
			if !ports.ToSet().IsSubset(activePorts) {
				p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(ports.Start()), int(ports.End()))
				rule := ds.Pair[*netset.IPBlock, netp.Protocol]{Left: cubes[i].Left.FirstIPAddressObject(), Right: p}
				activeRules = append(activeRules, rule)
			}
		}
	}
	// generate all existing rules
	return slices.Concat(res, createActiveRules(activeRules, cubes[len(cubes)-1].Left.LastIPAddressObject(), direction, l))
}

// icmpIPCubesToMinRules calls icmpIPCubesToRules func twice -- once when the any protocol cubes are converted
// to icmp cubes and onces when they are not. it returns the best result (less sg rules)
func icmpIPCubesToMinRules(cubes []ds.Pair[*netset.IPBlock, *netset.ICMPSet], anyProtocolCubes *netset.IPBlock, direction ir.Direction,
	l *netset.IPBlock) []*ir.SGRule {
	res := icmpIPCubesToRules(cubes, anyProtocolCubes, direction, l)
	anyAsICMP := partitionsToProduct(cubes).Union(ds.CartesianPairLeft(anyProtocolCubes, netset.AllICMPSet()))
	resWithAny := icmpIPCubesToRules(optimize.SortPartitionsByIPAddrs(anyAsICMP.Partitions()),
		netset.NewIPBlock(), direction, l) // pass an empty ipblock instead of anyProtocol iblock
	if len(resWithAny) < len(res) {
		return resWithAny
	}
	return res
}

// icmpIPCubesToRules converts cubes representing icmp protocol rules to SG rules
func icmpIPCubesToRules(cubes []ds.Pair[*netset.IPBlock, *netset.ICMPSet], anyProtocolCubes *netset.IPBlock, direction ir.Direction,
	l *netset.IPBlock) []*ir.SGRule {
	if len(cubes) == 0 {
		return []*ir.SGRule{}
	}

	res := make([]*ir.SGRule, 0)
	activeRules := make([]ds.Pair[*netset.IPBlock, netp.Protocol], 0) // first IP of the rule; protocol; in which cube we added the rule

	for i := range cubes {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && optimize.UncoveredHole(cubes[i-1].Left, cubes[i].Left, anyProtocolCubes) {
			res = slices.Concat(res, createActiveRules(activeRules, cubes[i-1].Left.LastIPAddressObject(), direction, l))
			activeRules = make([]ds.Pair[*netset.IPBlock, netp.Protocol], 0)
		}

		// if there are active rules whose icmp values are not fully included in the current cube, they will be created
		// also activeICMP will be calculated, which is the icmp values that are still included in the active rules
		activeICMP := netset.EmptyICMPSet()
		for j, rule := range slices.Backward(activeRules) {
			icmpSet := netset.ICMPSetFromICMP(rule.Right.(netp.ICMP))
			if !icmpSet.IsSubset(cubes[i].Right) {
				res = slices.Concat(res, createNewRules(rule.Right, rule.Left, cubes[i-1].Left.LastIPAddressObject(), direction, l))
				activeRules = slices.Delete(activeRules, j, j+1)
			} else {
				activeICMP.Union(icmpSet)
			}
		}

		// if the cube contains icmp values that are not contained in  active rules, new rules will be created
		for _, p := range optimize.IcmpsetPartitions(cubes[i].Right) {
			if !netset.ICMPSetFromICMP(p).IsSubset(activeICMP) {
				rule := ds.Pair[*netset.IPBlock, netp.Protocol]{Left: cubes[i].Left.FirstIPAddressObject(), Right: p}
				activeRules = append(activeRules, rule)
			}
		}
	}

	// generate all  existing rules
	return slices.Concat(res, createActiveRules(activeRules, cubes[len(cubes)-1].Left.LastIPAddressObject(), direction, l))
}

// creates sgRules from SG active rules
func createActiveRules(activeRules []ds.Pair[*netset.IPBlock, netp.Protocol], lastIP *netset.IPBlock,
	direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	res := make([]*ir.SGRule, 0)
	for _, pair := range activeRules {
		res = slices.Concat(res, createNewRules(pair.Right, pair.Left, lastIP, direction, l))
	}
	return res
}

// createNewRules breaks the startIP-endIP ip range into cidrs and creates SG rules
func createNewRules(protocol netp.Protocol, startIP, endIP *netset.IPBlock, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	ipRange, _ := netset.IPBlockFromIPRange(startIP, endIP)
	remoteCidrs := ipRange.SplitToCidrs()

	res := make([]*ir.SGRule, len(remoteCidrs))
	for i, remoteCidr := range remoteCidrs {
		res[i] = ir.NewSGRule(direction, remoteCidr, protocol, l, "")
	}
	return res
}

func partitionsToProduct[T ds.Set[T]](pairs []ds.Pair[*netset.IPBlock, T]) ds.Product[*netset.IPBlock, T] {
	res := ds.NewProductLeft[*netset.IPBlock, T]()
	for i := range pairs {
		res = res.Union(ds.CartesianPairLeft(pairs[i].Left, pairs[i].Right)).(*ds.ProductLeft[*netset.IPBlock, T])
	}
	return res
}
