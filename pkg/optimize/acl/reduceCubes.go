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
	allTCP := allTCPUDP(aclCubes.tcpAllow)
	allUDP := allTCPUDP(aclCubes.udpAllow)
	allicmp := allICMP(aclCubes.icmpAllow)
	aclCubes.anyProtocolAllow = allTCP.Intersect(allUDP).Intersect(allicmp)

	allTCPUDPnoICMP := allTCP.Intersect(allUDP).Subtract(anyICMP(aclCubes.icmpAllow))
	allTCPICMPnoUDP := allTCP.Intersect(allicmp).Subtract(anyTCPUDP(aclCubes.udpAllow))
	allUDPICMPnoTCP := allUDP.Intersect(allicmp).Subtract(anyTCPUDP(aclCubes.tcpAllow))

	// deny icmp, allow any
	aclCubes.icmpDeny = addSrcDstCubeToICMP(aclCubes.icmpDeny, allTCPUDPnoICMP)
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPUDPnoICMP)

	// deny udp, allow any
	aclCubes.udpDeny = addSrcDstCubeToTCPUDP(aclCubes.udpDeny, allTCPICMPnoUDP)
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allTCPICMPnoUDP)

	// deny tcp, allow any
	aclCubes.tcpDeny = addSrcDstCubeToTCPUDP(aclCubes.tcpDeny, allUDPICMPnoTCP)
	aclCubes.anyProtocolAllow = aclCubes.anyProtocolAllow.Union(allUDPICMPnoTCP)

	subtractAnyProtocolCubes(aclCubes)
}

func allTCPUDP(tcpudpAllow ds.TripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet]) *srcDstProductLeft {
	res := ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	allTCPSet := netset.NewAllTCPOnlySet()
	allUDPSet := netset.NewAllUDPOnlySet()
	for _, p := range tcpudpAllow.Partitions() {
		if p.S3.Equal(allTCPSet) || p.S3.Equal(allUDPSet) { // all tcp or udp ports
			r := ds.CartesianPairLeft(p.S1, p.S2)
			res = res.Union(r).(*srcDstProductLeft)
		}
	}
	return res
}

func anyTCPUDP(tcpudpAllow ds.TripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet]) *srcDstProductLeft {
	res := ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	for _, p := range tcpudpAllow.Partitions() {
		r := ds.CartesianPairLeft(p.S1, p.S2)
		res = res.Union(r).(*srcDstProductLeft)
	}
	return res
}

func allICMP(icmpAllow ds.TripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet]) *srcDstProductLeft {
	res := ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	for _, p := range icmpAllow.Partitions() {
		if p.S3.IsAll() { // all icmp types and codes
			r := ds.CartesianPairLeft(p.S1, p.S2)
			res = res.Union(r).(*srcDstProductLeft)
		}
	}
	return res
}

func anyICMP(icmpAllow ds.TripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet]) *srcDstProductLeft {
	res := ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	for _, p := range icmpAllow.Partitions() {
		r := ds.CartesianPairLeft(p.S1, p.S2)
		res = res.Union(r).(*srcDstProductLeft)
	}
	return res
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

// func subtractSrcDstCubeFromTCPUDP(tcpudpCube tcpudpTripleSet, srcDstCube srcDstProduct) tcpudpTripleSet {
// 	for _, p := range srcDstCube.Partitions() {
// 		t := ds.CartesianLeftTriple(p.Left, p.Right, netset.AllTCPUDPSet())
// 		tcpudpCube = tcpudpCube.Subtract(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet])
// 	}
// 	return tcpudpCube
// }

// func subtractSrcDstCubeFromICMP(icmpCube icmpTripleSet, srcDstCube srcDstProduct) icmpTripleSet {
// 	for _, p := range srcDstCube.Partitions() {
// 		t := ds.CartesianLeftTriple(p.Left, p.Right, netset.AllICMPSet())
// 		icmpCube = icmpCube.Subtract(t).(*ds.LeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet])
// 	}
// 	return icmpCube
// }

func addSrcDstCubeToTCPUDP(tcpudpCube tcpudpTripleSet, srcDstCube srcDstProduct) tcpudpTripleSet {
	for _, p := range srcDstCube.Partitions() {
		t := ds.CartesianLeftTriple(p.Left, p.Right, netset.AllTCPUDPSet())
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
