/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netset"
)

// reduceACLCubes unifies a (src ip x dst ip) cube, separately allowed for tcp, udp and icmp, into one "any" cube
// (assuming all ports, codes, types)
func reduceACLCubes(aclCubes *aclCubesPerProtocol) {
	allTCP, anyTCP, _ := tcpudpCubes(aclCubes.tcpAllow)
	allUDP, anyUDP, _ := tcpudpCubes(aclCubes.udpAllow)
	allICMP, anyICMP, _ := icmpCubes(aclCubes.icmpAllow)

	allTCPUDP := allTCP.Intersect(allUDP)
	allTCPICMP := allTCP.Intersect(allICMP)
	allUDPICMP := allUDP.Intersect(allICMP)

	aclCubes.anyProtocolAllow = allTCPUDP.Intersect(allICMP)

	// deny icmp, allow any
	allTCPUDPnoICMP := allTCPUDP.Subtract(anyICMP)
	aclCubes.icmpDeny = addSrcDstCubeToICMP(aclCubes.icmpDeny, allTCPUDPnoICMP)
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPUDPnoICMP)

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

func tcpudpCubes(tcpudpAllow tcpudpTripleSet) (allPorts, anyPorts, oneRule *srcDstProductLeft) {
	allTCPSet := netset.NewAllTCPOnlySet()
	allUDPSet := netset.NewAllUDPOnlySet()

	allPorts = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	anyPorts = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	oneRule = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()

	for _, p := range tcpudpAllow.Partitions() {
		r := ds.CartesianPairLeft(p.S1, p.S2)
		anyPorts = anyPorts.Union(r).(*srcDstProductLeft)
		if p.S3.Equal(allTCPSet) || p.S3.Equal(allUDPSet) { // all tcp or udp ports
			allPorts = allPorts.Union(r).(*srcDstProductLeft)
		}
		if oneExcludedTCPUDP(p.S3) {
			oneRule = oneRule.Union(r).(*srcDstProductLeft)
		}
	}
	return
}

func oneExcludedTCPUDP(_ *netset.TCPUDPSet) bool {
	return true
}

func icmpCubes(icmpAllow icmpTripleSet) (allICMP, anyICMP, oneValue *srcDstProductLeft) {
	allICMP = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	anyICMP = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	oneValue = ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()

	for _, p := range icmpAllow.Partitions() {
		r := ds.CartesianPairLeft(p.S1, p.S2)
		anyICMP = anyICMP.Union(r).(*srcDstProductLeft)
		if p.S3.IsAll() { // all icmp types and codes
			allICMP = allICMP.Union(r).(*srcDstProductLeft)
		}
		if oneExcludedICMP(p.S3) {
			oneValue = oneValue.Union(r).(*srcDstProductLeft)
		}
	}
	return
}

func oneExcludedICMP(_ *netset.ICMPSet) bool {
	return true
}

func subtractAnyProtocolCubes(aclCubes *aclCubesPerProtocol) {
	allTcpudp := ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet]()
	allIcmp := ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet]()
	for _, p := range aclCubes.anyProtocolAllow.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, netset.AllTCPUDPSet())
		allTcpudp = allTcpudp.Union(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet])
		i := ds.CartesianLeftTriple(p.Left, p.Right, netset.AllICMPSet())
		allIcmp = allIcmp.Union(i).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet])
	}

	aclCubes.tcpAllow = aclCubes.tcpAllow.Subtract(allTcpudp)
	aclCubes.udpAllow = aclCubes.udpAllow.Subtract(allTcpudp)
	aclCubes.icmpAllow = aclCubes.icmpAllow.Subtract(allIcmp)
}

func addSrcDstCubeToTCPUDP(tcpudpCube tcpudpTripleSet, srcDstCube srcDstProduct, pr *netset.TCPUDPSet) tcpudpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, pr)
		tcpudpCube = tcpudpCube.Union(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet])
	}
	return tcpudpCube
}

func addSrcDstCubeToICMP(icmpCube icmpTripleSet, srcDstCube srcDstProduct) icmpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, netset.AllICMPSet())
		icmpCube = icmpCube.Union(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet])
	}
	return icmpCube
}
