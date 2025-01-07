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
	// we will calculate the optimized deny cubes in `reduceACLCubes` func
	cubes.tcpDeny = ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]()
	cubes.udpDeny = ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]()
	cubes.icmpDeny = ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]()
	cubes.anyProtocolDeny = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()

	reduceACLCubes(cubes)

	denyTCPRules := protocolCubesToRules(cubes.tcpDeny, cubes.anyProtocolDeny, direction, netset.AllTCPTransport(), ir.Deny)
	allowTCPRules := protocolCubesToRules(cubes.tcpAllow, cubes.anyProtocolAllow, direction, netset.AllTCPTransport(), ir.Allow)

	denyUDPRules := protocolCubesToRules(cubes.udpDeny, cubes.anyProtocolDeny, direction, netset.AllUDPTransport(), ir.Deny)
	allowUDPRules := protocolCubesToRules(cubes.udpAllow, cubes.anyProtocolAllow, direction, netset.AllUDPTransport(), ir.Allow)

	denyICMPRules := protocolCubesToRules(cubes.icmpDeny, cubes.anyProtocolAllow, direction, netset.AllICMPTransport(), ir.Deny)
	allowICMPRules := protocolCubesToRules(cubes.icmpAllow, cubes.anyProtocolAllow, direction, netset.AllICMPTransport(), ir.Allow)

	denyAnyProtocolRules := anyProtocolCubesToRules(cubes.anyProtocolDeny, direction, ir.Deny)
	allowAnyProtocolRules := anyProtocolCubesToRules(cubes.anyProtocolAllow, direction, ir.Allow)
	return slices.Concat(denyTCPRules, allowTCPRules, denyUDPRules, allowUDPRules, denyICMPRules,
		allowICMPRules, denyAnyProtocolRules, allowAnyProtocolRules)
}

// Creates two sets of rules: one with only protocol cubes, and the other protocol cubes combined
// with any protocol cubes. It returns the minimal set
func protocolCubesToRules(tripleSet protocolTripleSet, anyCubes srcDstProduct, direction ir.Direction,
	pr *netset.TransportSet, action ir.Action) []*ir.ACLRule {
	res := tripleSetToCubes(tripleSet, direction, action)
	resWithAny := tripleSetToCubes(addSrcDstCubesToProtocolCubes(tripleSet, anyCubes, pr), direction, action)
	if len(resWithAny) < len(res) {
		return resWithAny
	}
	return res
}

// tripleSetToCubes calculates the minimal cubes partition and returns the corresponded slice of nACL rules
func tripleSetToCubes(tripleSet protocolTripleSet, direction ir.Direction, action ir.Action) []*ir.ACLRule {
	partitions := minimalCubesPartitions(tripleSet)
	res := make([]*ir.ACLRule, len(partitions))
	for i, t := range partitions {
		res[i] = ir.NewACLRule(action, direction, t.S1, t.S2, t.S3, "")
	}
	return res
}

// minimalCubesPartitions returns the minimal set of cubes partitions based on the triple set type
func minimalCubesPartitions(t protocolTripleSet) []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol] {
	leftPartitions := actualPartitions(ds.AsLeftTripleSet(t))
	outerPartitions := actualPartitions(ds.AsOuterTripleSet(t))
	rightPartitions := actualPartitions(ds.AsRightTripleSet(t))

	switch {
	case len(leftPartitions) <= len(outerPartitions) && len(leftPartitions) <= len(rightPartitions):
		return leftPartitions
	case len(outerPartitions) <= len(leftPartitions) && len(outerPartitions) <= len(rightPartitions):
		return outerPartitions
	default:
		return rightPartitions
	}
}

func actualPartitions(t protocolTripleSet) []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol], 0)
	for _, p := range t.Partitions() {
		res = slices.Concat(res, breakTCPUDPTriple(p), breakICMPTriple(p)) // here, one function returns an empty slice
	}
	return res
}

// break multi-cube to regular cube
func breakTCPUDPTriple(t ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]) []ds.Triple[*netset.IPBlock,
	*netset.IPBlock, netp.Protocol] {
	if t.S3.TCPUDPSet().IsEmpty() {
		return []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol]{}
	}
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol], 0)

	dstCidrs := t.S2.SplitToCidrs()
	tcpudpTriples := t.S3.TCPUDPSet().Partitions()
	isTCP := tcpudpTriples[0].S1.Elements()[0] == netset.TCPCode
	for _, src := range t.S1.SplitToCidrs() {
		for _, dst := range dstCidrs {
			for _, protocolTriple := range tcpudpTriples {
				tcpudpSrcPorts := protocolTriple.S2.Intervals()
				tcpudpDstPorts := protocolTriple.S3.Intervals()
				for _, srcPorts := range tcpudpSrcPorts {
					for _, dstPorts := range tcpudpDstPorts {
						p, _ := netp.NewTCPUDP(isTCP, int(srcPorts.Start()), int(srcPorts.End()), int(dstPorts.Start()), int(dstPorts.End()))
						res = append(res, ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol]{S1: src, S2: dst, S3: p})
					}
				}
			}
		}
	}
	return res
}

// break multi-cube to regular cube
func breakICMPTriple(t ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]) []ds.Triple[*netset.IPBlock,
	*netset.IPBlock, netp.Protocol] {
	if t.S3.ICMPSet().IsEmpty() {
		return []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol]{}
	}
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol], 0)

	dstCidrs := t.S2.SplitToCidrs()
	icmpPartitions := optimize.IcmpsetPartitions(t.S3.ICMPSet())
	for _, src := range t.S1.SplitToCidrs() {
		for _, dst := range dstCidrs {
			for _, icmp := range icmpPartitions {
				a := ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.Protocol]{S1: src, S2: dst, S3: icmp}
				res = append(res, a)
			}
		}
	}
	return res
}
