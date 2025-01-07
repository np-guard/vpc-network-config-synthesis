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
	allTCP, anyTCP := divideCubes(aclCubes.tcpAllow)
	allUDP, anyUDP := divideCubes(aclCubes.udpAllow)
	allICMP, anyICMP := divideCubes(aclCubes.icmpAllow)

	allTCPUDP := allTCP.Intersect(allUDP)
	allTCPICMP := allTCP.Intersect(allICMP)
	allUDPICMP := allUDP.Intersect(allICMP)

	aclCubes.anyProtocolAllow = allTCPUDP.Intersect(allICMP)

	// replace TCP and UDP rules to deny ICMP and allow any
	allTCPUDPnoICMP := allTCPUDP.Subtract(anyICMP)
	aclCubes.icmpDeny = addSrcDstCubesToProtocolCubes(aclCubes.icmpDeny, allTCPUDPnoICMP, netset.AllICMPTransport())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPUDPnoICMP)

	// replace TCP and ICMP rules to deny UDP and allow any
	allTCPICMPnoUDP := allTCPICMP.Subtract(anyUDP)
	aclCubes.udpDeny = addSrcDstCubesToProtocolCubes(aclCubes.udpDeny, allTCPICMPnoUDP, netset.AllUDPTransport())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPICMPnoUDP)

	// replace ICMP and UDP rules to deny TCP and allow any
	allUDPICMPnoTCP := allUDPICMP.Subtract(anyTCP)
	aclCubes.tcpDeny = addSrcDstCubesToProtocolCubes(aclCubes.tcpDeny, allUDPICMPnoTCP, netset.AllTCPTransport())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allUDPICMPnoTCP)

	// if there are two protocols that are complete (all ports in the case of tcp/udp or all types&codes in icmp), and the complement of the
	// third protocol can be written in one rule we will replace it with two rules: deny for the complement protocol and allow any. e.g.
	//
	// ALLOW src --> dst (tcp, icmp, udp ports 1-100)
	// will be replaced with:
	// DENY src --> dst (udp ports 101-65535)
	// ALLOW src --> dst (any protocol)
	aclCubes.udpAllow, aclCubes.udpDeny = excludeProtocol(aclCubes.udpAllow, aclCubes.udpDeny, aclCubes, allTCPICMP)
	aclCubes.tcpAllow, aclCubes.tcpDeny = excludeProtocol(aclCubes.tcpAllow, aclCubes.tcpDeny, aclCubes, allUDPICMP)
	aclCubes.icmpAllow, aclCubes.icmpDeny = excludeProtocol(aclCubes.icmpAllow, aclCubes.icmpDeny, aclCubes, allTCPUDP)

	subtractAnyProtocolCubes(aclCubes)
}

// divideCubes returns two <src X dst> products; one that represents all connections that allow all protocol combinations,
// and one that represents all connections that allow at least one combination
func divideCubes(protocolCubes protocolTripleSet) (allCombinations, anyCombination srcDstProduct) {
	allTCPSet := netset.AllTCPTransport()
	allUDPSet := netset.AllUDPTransport()
	allICMPSet := netset.AllICMPTransport()

	allCombinations = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	anyCombination = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()

	for _, p := range protocolCubes.Partitions() {
		r := ds.CartesianPairLeft(p.S1, p.S2)
		anyCombination = anyCombination.Union(r)
		if p.S3.Equal(allTCPSet) || p.S3.Equal(allUDPSet) || p.S3.Equal(allICMPSet) {
			allCombinations = allCombinations.Union(r)
		}
	}
	return
}

// Updates the cubes so that if for <src, dst> the other two protocols are complete (for tcp/udp all
// ports and for icmp all types&codes) and the complement of the third protocol can be written by one rule
// we convert it to two rules, deny complement and allow any
func excludeProtocol(allowCubes, denyCubes protocolTripleSet, cubes *aclCubesPerProtocol,
	allOtherProtocols srcDstProduct) (allow, deny protocolTripleSet) {
	var complementSet *netset.TransportSet
	var single bool

	for _, p := range allowCubes.Partitions() {
		if complementSet, single = singleComplementValue(p.S3); !single {
			continue
		}
		intersectedCube := ds.CartesianPairLeft(p.S1, p.S2).Intersect(allOtherProtocols)
		allowCubes = subtractSrcDstCubesFromProtocolCubes(allowCubes, intersectedCube, p.S3)
		denyCubes = addSrcDstCubesToProtocolCubes(denyCubes, intersectedCube, complementSet)
		cubes.anyProtocolAllow = cubes.anyProtocolAllow.Union(intersectedCube)
	}
	return allowCubes, denyCubes
}

// subtractAllowAnyProtocolCubes from each protocol
func subtractAnyProtocolCubes(aclCubes *aclCubesPerProtocol) {
	aclCubes.tcpAllow = subtractSrcDstCubesFromProtocolCubes(aclCubes.tcpAllow, aclCubes.anyProtocolAllow, netset.AllTCPTransport())
	aclCubes.udpAllow = subtractSrcDstCubesFromProtocolCubes(aclCubes.udpAllow, aclCubes.anyProtocolAllow, netset.AllUDPTransport())
	aclCubes.icmpAllow = subtractSrcDstCubesFromProtocolCubes(aclCubes.icmpAllow, aclCubes.anyProtocolAllow, netset.AllICMPTransport())
}

// Union <src X dst> cube with pr TransportSet to protocol cube
func addSrcDstCubesToProtocolCubes(protocolCubes protocolTripleSet, srcDstCube srcDstProduct, pr *netset.TransportSet) protocolTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		protocolCubes = protocolCubes.Union(t)
	}
	return protocolCubes
}

// Subtract <src X dst> cube with pr protocolSet from protocol cube
func subtractSrcDstCubesFromProtocolCubes(cubes protocolTripleSet, srcDstCube srcDstProduct, pr *netset.TransportSet) protocolTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		cubes = cubes.Subtract(t)
	}
	return cubes
}

// Note: only one of tcp, udp, icmp
// returns whether the complement of the given set can be written in one rule (and the complement set)
func singleComplementValue(protocolSet *netset.TransportSet) (*netset.TransportSet, bool) {
	if icmpSet := protocolSet.ICMPSet(); !icmpSet.IsEmpty() {
		return singleICMPComplementValue(icmpSet)
	}
	return singleTCPUDPComplementValue(protocolSet.TCPUDPSet())
}

// Note: only tcp or udp set
// returns whether the complement of the given set can be written in one rule (and the complement set)
func singleTCPUDPComplementValue(tcpudpSet *netset.TCPUDPSet) (*netset.TransportSet, bool) {
	set := netset.NewAllTCPOnlySet()
	if set.Intersect(tcpudpSet).IsEmpty() {
		set = netset.NewAllUDPOnlySet()
	}
	complementSet := set.Subtract(tcpudpSet)
	return netset.NewTCPUDPTransportFromTCPUDPSet(complementSet), len(complementSet.Partitions()) == 1
}

// returns whether the complement of the given set can be written in one rule (and the complement set)
func singleICMPComplementValue(icmpSet *netset.ICMPSet) (*netset.TransportSet, bool) {
	complementSet := netset.AllICMPSet().Subtract(icmpSet)
	return netset.NewICMPTransportFromICMPSet(complementSet), len(optimize.IcmpsetPartitions(complementSet)) == 1
}
