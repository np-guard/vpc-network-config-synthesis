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

// reduceACLCubes unifies a (src ip x dst ip) cube, separately allowed for tcp, udp and icmp, into one "any" cube
// (assuming all ports, codes, types)
func reduceACLCubes(aclCubes *aclCubesPerProtocol) {
	allTCP, anyTCP := tcpudpCubes(aclCubes.tcpAllow)
	allUDP, anyUDP := tcpudpCubes(aclCubes.udpAllow)
	allICMP, anyICMP := icmpCubes(aclCubes.icmpAllow)

	allTCPUDP := allTCP.Intersect(allUDP)
	allTCPICMP := allTCP.Intersect(allICMP)
	allUDPICMP := allUDP.Intersect(allICMP)

	aclCubes.anyProtocolAllow = allTCPUDP.Intersect(allICMP)

	// deny icmp, allow any
	allTCPUDPnoICMP := allTCPUDP.Subtract(anyICMP)
	aclCubes.icmpDeny = addSrcDstCubeToICMP(aclCubes.icmpDeny, allTCPUDPnoICMP, netset.AllICMPSet())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPUDPnoICMP)
	excludeICMP(aclCubes, allTCPUDP)

	// deny udp, allow any
	allTCPICMPnoUDP := allTCPICMP.Subtract(anyUDP)
	aclCubes.udpDeny = addSrcDstCubeToTCPUDP(aclCubes.udpDeny, allTCPICMPnoUDP, netset.NewAllUDPOnlySet())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPICMPnoUDP)

	// deny tcp, allow any
	allUDPICMPnoTCP := allUDPICMP.Subtract(anyTCP)
	aclCubes.tcpDeny = addSrcDstCubeToTCPUDP(aclCubes.tcpDeny, allUDPICMPnoTCP, netset.NewAllTCPOnlySet())
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allUDPICMPnoTCP)

	subtractAnyProtocolCubes(aclCubes)
}

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

func excludeICMP(cubes *aclCubesPerProtocol, allTCPUDP srcDstProduct) {
	var icmpSet *netset.ICMPSet
	var single bool
	for _, p := range cubes.icmpAllow.Partitions() {
		if icmpSet, single = singleICMPValue(p.S3); !single {
			continue
		}
		currCube := ds.CartesianPairLeft(p.S1, p.S2).Intersect(allTCPUDP)
		cubes.icmpAllow = subtractSrcDstCubeFromICMP(cubes.icmpAllow, currCube, p.S3)
		cubes.icmpDeny = addSrcDstCubeToICMP(cubes.icmpDeny, currCube, icmpSet)
		cubes.anyProtocolAllow = cubes.anyProtocolAllow.Union(currCube).(*srcDstProductLeft)
	}
}

func singleICMPValue(icmpSet *netset.ICMPSet) (*netset.ICMPSet, bool) {
	complementSet := netset.AllICMPSet().Subtract(icmpSet)
	return complementSet, len(optimize.IcmpsetPartitions(complementSet)) == 1
}

func subtractAnyProtocolCubes(aclCubes *aclCubesPerProtocol) {
	aclCubes.tcpAllow = subtractSrcDstCubeFromTCPUDP(aclCubes.tcpAllow, aclCubes.anyProtocolAllow, netset.NewAllTCPOnlySet())
	aclCubes.udpAllow = subtractSrcDstCubeFromTCPUDP(aclCubes.udpAllow, aclCubes.anyProtocolAllow, netset.NewAllUDPOnlySet())
	aclCubes.icmpAllow = subtractSrcDstCubeFromICMP(aclCubes.icmpAllow, aclCubes.anyProtocolAllow, netset.AllICMPSet())
}

func addSrcDstCubeToTCPUDP(tcpudpCube tcpudpTripleSet, srcDstCube srcDstProduct, pr *netset.TCPUDPSet) tcpudpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		tcpudpCube = tcpudpCube.Union(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet])
	}
	return tcpudpCube
}

func subtractSrcDstCubeFromTCPUDP(tcpudpCube tcpudpTripleSet, srcDstCube srcDstProduct, pr *netset.TCPUDPSet) tcpudpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		tcpudpCube = tcpudpCube.Subtract(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet])
	}
	return tcpudpCube
}

func addSrcDstCubeToICMP(icmpCube icmpTripleSet, srcDstCube srcDstProduct, pr *netset.ICMPSet) icmpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		icmpCube = icmpCube.Union(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet])
	}
	return icmpCube
}

func subtractSrcDstCubeFromICMP(icmpCube icmpTripleSet, srcDstCube srcDstProduct, pr *netset.ICMPSet) icmpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		icmpCube = icmpCube.Subtract(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet])
	}
	return icmpCube
}
