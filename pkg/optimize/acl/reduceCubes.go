/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netset"
)

func reduceACLCubes(aclCubes *aclCubesPerProtocol) {
	allUDPTCP := allTCPUDP(aclCubes.tcpAllow).Intersect(allTCPUDP(aclCubes.udpAllow))
	aclCubes.anyProtocolAllow = allUDPTCP.Intersect(allICMP(aclCubes.icmpAllow))
	subtractAnyProtocolCubes(aclCubes)
}

func allTCPUDP(tcpudpAllow ds.TripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet]) *srcDstProductLeft {
	res := ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	for _, p := range tcpudpAllow.Partitions() {
		if p.S3.Equal(netset.NewAllTCPOnlySet()) || p.S3.Equal(netset.NewAllUDPOnlySet()) { // all tcp or udp ports
			r := ds.CartesianPairLeft(p.S1, p.S2)
			res = res.Union(r).(*srcDstProductLeft)
		}
	}
	return res
}

func allICMP(icmpAllow ds.TripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet]) srcDstProduct {
	res := ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock]()
	for _, p := range icmpAllow.Partitions() {
		if p.S3.IsAll() { // all icmp types and codes
			r := ds.CartesianPairLeft(p.S1, p.S2)
			res = res.Union(r).(*srcDstProductLeft)
		}
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
