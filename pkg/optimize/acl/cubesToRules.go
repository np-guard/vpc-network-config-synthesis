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

type (
	dstProtocolProduct = ds.Product[*netset.IPBlock, *netset.TransportSet]
	activeRule         = ds.Pair[*netset.IPBlock, dstProtocolProduct]
)

func aclCubesToRules(cubes *aclCubesPerProtocol, direction ir.Direction) []*ir.ACLRule {
	// we will calculate the optimized deny cubes in `reduceACLCubes` func
	cubes.tcpDeny = ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]()
	cubes.udpDeny = ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]()
	cubes.icmpDeny = ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]()
	cubes.anyProtocolDeny = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()

	reduceACLCubes(cubes)

	denyTCPRules := minRulesPartitions(cubes.tcpDeny, cubes.anyProtocolDeny, direction, ir.Deny, netset.AllTCPTransport())
	allowTCPRules := minRulesPartitions(cubes.tcpAllow, cubes.anyProtocolAllow, direction, ir.Allow, netset.AllTCPTransport())

	denyUDPRules := minRulesPartitions(cubes.udpDeny, cubes.anyProtocolDeny, direction, ir.Deny, netset.AllUDPTransport())
	allowUDPRules := minRulesPartitions(cubes.udpAllow, cubes.anyProtocolAllow, direction, ir.Allow, netset.AllUDPTransport())

	denyICMPRules := minRulesPartitions(cubes.icmpDeny, cubes.anyProtocolAllow, direction, ir.Deny, netset.AllICMPTransport())
	allowICMPRules := minRulesPartitions(cubes.icmpAllow, cubes.anyProtocolAllow, direction, ir.Allow, netset.AllICMPTransport())

	denyAnyProtocolRules := cubesToRules(convertAnyCubesToTripleSet(cubes.anyProtocolDeny), cubes.anyProtocolDeny, direction, ir.Deny)
	allowAnyProtocolRules := cubesToRules(convertAnyCubesToTripleSet(cubes.anyProtocolAllow), cubes.anyProtocolAllow, direction, ir.Allow)
	return slices.Concat(denyTCPRules, allowTCPRules, denyUDPRules, allowUDPRules, denyICMPRules,
		allowICMPRules, denyAnyProtocolRules, allowAnyProtocolRules)
}

// Creates two sets of rules: one with only protocol cubes, and the other protocol cubes combined
// with any protocol cubes. It returns the minimal set
func minRulesPartitions(tripleSet protocolTripleSet, anyCubes srcDstProduct, direction ir.Direction, action ir.Action,
	pr *netset.TransportSet) []*ir.ACLRule {
	res := minimalCubesPartitions(tripleSet, anyCubes, direction, action)
	resWithAny := minimalCubesPartitions(addSrcDstCubesToProtocolCubes(tripleSet, anyCubes, pr), anyCubes, direction, action)
	if len(resWithAny) < len(res) {
		return resWithAny
	}
	return res
}

// minimalCubesPartitions returns the minimal set of cubes partitions based on the triple set type
func minimalCubesPartitions(tripleSet protocolTripleSet, anyCubes srcDstProduct, direction ir.Direction, action ir.Action) []*ir.ACLRule {
	// only in LeftTripleSet and in OuterTripleSet S1 is src IP
	leftPartitions := cubesToRules(ds.AsLeftTripleSet(tripleSet), anyCubes, direction, action)
	outerPartitions := cubesToRules(ds.AsOuterTripleSet(tripleSet), anyCubes, direction, action)

	if len(leftPartitions) <= len(outerPartitions) {
		return leftPartitions
	}
	return outerPartitions
}

// based on sg optimization algorithm, but in this case activeRules map a srcIP (block start)
// to a pair of a single dstIP CIDR and a single protocol partition
func cubesToRules(cubes protocolTripleSet, anyProtocolCubes srcDstProduct, direction ir.Direction, action ir.Action) []*ir.ACLRule {
	partitions := convertCubesType(cubes.Partitions())
	if len(partitions) == 0 {
		return []*ir.ACLRule{}
	}
	anyProtocolSrcIPs := anyProtocolCubes.(*ds.ProductLeft[*netset.IPBlock, *netset.IPBlock]).Left(netset.NewIPBlock())

	res := make([]*ir.ACLRule, 0)
	activeRules := make([]activeRule, 0) // Left = first src's IP, Right = dst cidr & protocol details

	for i := range partitions {
		// if it is not possible to continue the rule between the cubes, generate all existing rules
		if i > 0 && optimize.UncoveredHole(partitions[i-1].Left, partitions[i].Left, anyProtocolSrcIPs) {
			res = slices.Concat(res, createActiveRules(activeRules, partitions[i-1].Left.LastIPAddressObject(), direction, action))
			activeRules = make([]activeRule, 0)
		}

		// if there are active rules whose cubeDetails are not fully included in the current cube, they will be created
		// also activeCubes will be calculated, which is the activeCubess that are still included in the active rules
		activeCubes := ds.NewProductLeft[*netset.IPBlock, *netset.TransportSet]()
		for j, rule := range slices.Backward(activeRules) {
			if rule.Right.IsSubset(partitions[i].Right) {
				activeCubes = activeCubes.Union(rule.Right).(*ds.ProductLeft[*netset.IPBlock, *netset.TransportSet])
			} else {
				res = slices.Concat(res,
					createNewRules(rule.Left, partitions[i-1].Left.LastIPAddressObject(), rule.Right.Partitions()[0], direction, action))
				activeRules = slices.Delete(activeRules, j, j+1)
			}
		}

		// if the current cube contains values that are not contained in active rules, new rules will be created
		for _, currCube := range partitions[i].Right.Partitions() {
			dstPortCidrs := currCube.Left.SplitToCidrs()
			for _, p := range transportSetToProtocols(currCube.Right) {
				for _, dstCidr := range dstPortCidrs {
					cubeDetails := ds.CartesianPairLeft(dstCidr, p)
					if !cubeDetails.IsSubset(activeCubes) {
						rule := activeRule{Left: partitions[i].Left.FirstIPAddressObject(), Right: cubeDetails}
						activeRules = append(activeRules, rule)
					}
				}
			}
		}
	}
	// generate all existing rules
	return slices.Concat(res, createActiveRules(activeRules, partitions[len(partitions)-1].Left.LastIPAddressObject(), direction, action))
}

func createActiveRules(activeRules []activeRule, srcLastIP *netset.IPBlock, direction ir.Direction, action ir.Action) []*ir.ACLRule {
	res := make([]*ir.ACLRule, 0)
	for _, rule := range activeRules {
		res = slices.Concat(res, createNewRules(rule.Left, srcLastIP, rule.Right.Partitions()[0], direction, action))
	}
	return res
}

func createNewRules(srcStartIP, srcEndIP *netset.IPBlock, cubeDetails ds.Pair[*netset.IPBlock, *netset.TransportSet],
	direction ir.Direction, action ir.Action) []*ir.ACLRule {
	src, _ := netset.IPBlockFromIPRange(srcStartIP, srcEndIP)
	srcCidrs := src.SplitToCidrs()

	res := make([]*ir.ACLRule, len(srcCidrs))
	for i, srcCidr := range srcCidrs {
		res[i] = ir.NewACLRule(action, direction, srcCidr, cubeDetails.Left, transportSetToProtocol(cubeDetails.Right), "")
	}
	return res
}

// transportSetToProtocols returns a slice of TransportSets, each one is a valid nACL rule protocol
func transportSetToProtocols(t *netset.TransportSet) []*netset.TransportSet {
	if t.IsAll() {
		return []*netset.TransportSet{t}
	}
	res := make([]*netset.TransportSet, 0)
	for _, icmp := range optimize.IcmpsetPartitions(t.ICMPSet()) {
		res = append(res, netset.NewICMPTransportFromICMPSet(netset.ICMPSetFromICMP(icmp)))
	}
	tcpudpPartitions := t.TCPUDPSet().Partitions()
	if len(tcpudpPartitions) == 0 {
		return res
	}
	protocolString := netp.ProtocolStringUDP
	if tcpudpPartitions[0].S1.Elements()[0] == netset.TCPCode { // tcp
		protocolString = netp.ProtocolStringTCP
	}
	for _, tcpudp := range tcpudpPartitions {
		tcpudpDstPorts := tcpudp.S3.Intervals()
		for _, srcPorts := range tcpudp.S2.Intervals() {
			for _, dstPorts := range tcpudpDstPorts {
				p := netset.NewTCPorUDPTransport(protocolString, srcPorts.Start(), srcPorts.End(), dstPorts.Start(), dstPorts.End())
				res = append(res, p)
			}
		}
	}
	return res
}

// // assuming the transport set contains a single protocol cube that can be used in a single nACL rule
func transportSetToProtocol(t *netset.TransportSet) netp.Protocol {
	icmpSet := t.ICMPSet()
	tcpudpSet := t.TCPUDPSet()

	switch {
	case t.IsAll():
		return netp.AnyProtocol{}
	case !icmpSet.IsEmpty():
		return optimize.IcmpsetPartitions(icmpSet)[0]
	}
	p := tcpudpSet.Partitions()[0]
	srcPorts := p.S2.Intervals()[0]
	dstPorts := p.S3.Intervals()[0]
	res, _ := netp.NewTCPUDP(p.S1.Elements()[0] == netset.TCPCode, int(srcPorts.Start()), int(srcPorts.End()),
		int(dstPorts.Start()), int(dstPorts.End()))
	return res
}

// converts cubes from a slices of triples to a slice of `activeRule` type
func convertCubesType(cubes []ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]) []activeRule {
	res := make([]activeRule, len(cubes))
	for i := range cubes {
		res[i] = activeRule{Left: cubes[i].S1, Right: ds.CartesianPairLeft(cubes[i].S2, cubes[i].S3)}
	}
	cmp := func(i, j activeRule) int { return i.Left.Compare(j.Left) }
	slices.SortFunc(res, cmp)
	return res
}

func convertAnyCubesToTripleSet(cubes srcDstProduct) protocolTripleSet {
	res := ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]()
	for _, p := range cubes.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, netset.AllTransports())
		res = res.Union(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet])
	}

	return res
}
