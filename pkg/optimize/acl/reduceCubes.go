/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

// reduceACLCubes attempts to reduce the total number of ACL cubes used to represent allowed/denied connections
func reduceACLCubes(aclCubes *aclCubesPerProtocol) {
	allTCP, anyTCP := tcpudpCubes(aclCubes.tcpAllow)
	allUDP, anyUDP := tcpudpCubes(aclCubes.udpAllow)
	allICMP, anyICMP := icmpCubes(aclCubes.icmpAllow)

	allTCPUDP := allTCP.Intersect(allUDP)
	allTCPICMP := allTCP.Intersect(allICMP)
	allUDPICMP := allUDP.Intersect(allICMP)

	aclCubes.anyProtocolAllow = allTCPUDP.Intersect(allICMP)

	// replace TCP and UDP rules to deny ICMP and allow any
	allTCPUDPnoICMP := allTCPUDP.Subtract(anyICMP)
	aclCubes.icmpDeny = addSrcDstCubeToICMP(aclCubes.icmpDeny, allTCPUDPnoICMP, netset.AllICMPSet())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPUDP)

	// replace TCP and ICMP rules to deny UDP and allow any
	allTCPICMPnoUDP := allTCPICMP.Subtract(anyUDP)
	aclCubes.udpDeny = addSrcDstCubeToTCPUDP(aclCubes.udpDeny, allTCPICMPnoUDP, netset.NewAllUDPOnlySet())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPICMPnoUDP)

	// replace ICMP and UDP rules to deny TCP and allow any
	allUDPICMPnoTCP := allUDPICMP.Subtract(anyTCP)
	aclCubes.tcpDeny = addSrcDstCubeToTCPUDP(aclCubes.tcpDeny, allUDPICMPnoTCP, netset.NewAllTCPOnlySet())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allUDPICMPnoTCP)

	// if there are two protocols that are complete (all ports in the case of tcp/udp or all types&codes in icmp), and the complement of the
	// third protocol can be written in one rule we will replace it with two rules: deny for the complement protocol and allow any. e.g.
	//
	// ALLOW src --> dst (tcp, icmp, udp ports 1-100)
	// will be replaced with:
	// DENY src --> dst (udp ports 101-65535)
	// ALLOW src --> dst (any protocol)
	excludeICMP(aclCubes, allTCPUDP)
	aclCubes.udpAllow, aclCubes.udpDeny = excludeTCPUDP(aclCubes.udpAllow, aclCubes.udpDeny, aclCubes, allTCPICMP)
	aclCubes.tcpAllow, aclCubes.tcpDeny = excludeTCPUDP(aclCubes.tcpAllow, aclCubes.tcpDeny, aclCubes, allUDPICMP)

	subtractAnyProtocolCubes(aclCubes)
}

// tcpudpCubes returns two <src X dst> products; one that represents all connections that allow all ports,
// and one that represents all connections that allow at least one port
func tcpudpCubes(tcpudpAllow tcpudpTripleSet) (allPorts, anyPorts *srcDstProductLeft) {
	allTCPSet := netset.NewAllTCPOnlySet()
	allUDPSet := netset.NewAllUDPOnlySet()

	allPorts = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	anyPorts = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()

	for _, p := range tcpudpAllow.Partitions() {
		r := ds.CartesianPairLeft(p.S1, p.S2)
		anyPorts = anyPorts.Union(r).(*srcDstProductLeft)
		if p.S3.Equal(allTCPSet) || p.S3.Equal(allUDPSet) { // all tcp or udp ports
			allPorts = allPorts.Union(r).(*srcDstProductLeft)
		}
	}
	return
}

// Updates the cubes so that if for <src, dst> the other two protocols are complete (for tcp/udp all
// ports and for icmp all types&codes) and the complement of the third protocol can be written by one rule
// we convert it to two rules, deny complement and allow any
func excludeTCPUDP(allowCubes, denyCubes tcpudpTripleSet, cubes *aclCubesPerProtocol,
	allOtherProtocols srcDstProduct) (allow, deny tcpudpTripleSet) {
	var tcpudpSet *netset.TCPUDPSet
	var single bool

	for _, p := range allowCubes.Partitions() {
		if tcpudpSet, single = singleTCPUDPComplementValue(p.S3); !single {
			continue
		}
		intersectedCube := ds.CartesianPairLeft(p.S1, p.S2).Intersect(allOtherProtocols)
		allowCubes = subtractSrcDstCubeFromTCPUDP(allowCubes, intersectedCube, p.S3)
		denyCubes = addSrcDstCubeToTCPUDP(denyCubes, intersectedCube, tcpudpSet)
		cubes.anyProtocolAllow = cubes.anyProtocolAllow.Union(intersectedCube).(*srcDstProductLeft)
	}
	return allowCubes, denyCubes
}

// Note: only tcp or udp set
// returns whether the complement of the given set can be written in one rule (and the complement set)
func singleTCPUDPComplementValue(tcpudpSet *netset.TCPUDPSet) (*netset.TCPUDPSet, bool) {
	set := netset.NewAllTCPOnlySet()
	if set.Intersect(tcpudpSet).IsEmpty() {
		set = netset.NewAllUDPOnlySet()
	}
	complementSet := set.Subtract(tcpudpSet)
	return complementSet, len(complementSet.Partitions()) == 1
}

// icmpCubes returns two <src X dst> products; one that represents all connections that allow
// all types and codes, and one that represents all connections with at least one ICMP value.
func icmpCubes(icmpAllow icmpTripleSet) (allICMP, anyICMP *srcDstProductLeft) {
	allICMP = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	anyICMP = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()

	for _, p := range icmpAllow.Partitions() {
		r := ds.CartesianPairLeft(p.S1, p.S2)
		anyICMP = anyICMP.Union(r).(*srcDstProductLeft)
		if p.S3.IsAll() { // all icmp types and codes
			allICMP = allICMP.Union(r).(*srcDstProductLeft)
		}
	}
	return
}

// Updates the cubes so that if for <src, dst> the other two protocols are complete (for tcp/udp all
// ports and for icmp all types&codes) and the complement of the third protocol can be written by one rule
// we convert it to two rules, deny complement and allow any
func excludeICMP(cubes *aclCubesPerProtocol, allTCPUDP srcDstProduct) {
	var icmpSet *netset.ICMPSet
	var single bool
	for _, p := range cubes.icmpAllow.Partitions() {
		if icmpSet, single = singleICMPComplementValue(p.S3); !single {
			continue
		}
		intersectedCube := ds.CartesianPairLeft(p.S1, p.S2).Intersect(allTCPUDP)
		cubes.icmpAllow = subtractSrcDstCubeFromICMP(cubes.icmpAllow, intersectedCube, p.S3)
		cubes.icmpDeny = addSrcDstCubeToICMP(cubes.icmpDeny, intersectedCube, icmpSet)
		cubes.anyProtocolAllow = cubes.anyProtocolAllow.Union(intersectedCube).(*srcDstProductLeft)
	}
}

// returns whether the complement of the given set can be written in one rule (and the complement set)
func singleICMPComplementValue(icmpSet *netset.ICMPSet) (*netset.ICMPSet, bool) {
	complementSet := netset.AllICMPSet().Subtract(icmpSet)
	return complementSet, len(optimize.IcmpsetPartitions(complementSet)) == 1
}

// subtractAnyProtocolCubes from each protocol
func subtractAnyProtocolCubes(aclCubes *aclCubesPerProtocol) {
	aclCubes.tcpAllow = subtractSrcDstCubeFromTCPUDP(aclCubes.tcpAllow, aclCubes.anyProtocolAllow, netset.NewAllTCPOnlySet())
	aclCubes.udpAllow = subtractSrcDstCubeFromTCPUDP(aclCubes.udpAllow, aclCubes.anyProtocolAllow, netset.NewAllUDPOnlySet())
	aclCubes.icmpAllow = subtractSrcDstCubeFromICMP(aclCubes.icmpAllow, aclCubes.anyProtocolAllow, netset.AllICMPSet())
}

// Union <src X dst> cube with pr protocolSet to tcpudp cube
func addSrcDstCubeToTCPUDP(tcpudpCube tcpudpTripleSet, srcDstCube srcDstProduct, pr *netset.TCPUDPSet) tcpudpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		tcpudpCube = tcpudpCube.Union(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet])
	}
	return tcpudpCube
}

// Subtract <src X dst> cube with pr protocolSet from tcpudp cube
func subtractSrcDstCubeFromTCPUDP(tcpudpCube tcpudpTripleSet, srcDstCube srcDstProduct, pr *netset.TCPUDPSet) tcpudpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		tcpudpCube = tcpudpCube.Subtract(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet])
	}
	return tcpudpCube
}

// Union <src X dst> cube with pr protocolSet to icmp cube
func addSrcDstCubeToICMP(icmpCube icmpTripleSet, srcDstCube srcDstProduct, pr *netset.ICMPSet) icmpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		icmpCube = icmpCube.Union(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet])
	}
	return icmpCube
}

// Subtract <src X dst> cube with pr protocolSet from icmp cube
func subtractSrcDstCubeFromICMP(icmpCube icmpTripleSet, srcDstCube srcDstProduct, pr *netset.ICMPSet) icmpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		icmpCube = icmpCube.Subtract(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet])
	}
	return icmpCube
}
