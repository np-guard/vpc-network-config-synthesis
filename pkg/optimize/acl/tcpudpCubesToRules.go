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
)

func tcpudpTriplesToRules(tripleSet tcpudpTripleSet, direction ir.Direction) []*ir.ACLRule {
	partitions := minimalPartitionsTCPUDP(tripleSet)
	res := make([]*ir.ACLRule, len(partitions))
	for i, t := range partitions {
		res[i] = ir.NewACLRule(ir.Allow, direction, t.S1, t.S2, t.S3, "")
	}
	return res
}

func minimalPartitionsTCPUDP(t tcpudpTripleSet) []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP] {
	leftPartitions := actualPartitionsTCPUDP(ds.AsLeftTripleSet(t))
	outerPartitions := actualPartitionsTCPUDP(ds.AsOuterTripleSet(t))
	rightPartitions := actualPartitionsTCPUDP(ds.AsRightTripleSet(t))

	switch {
	case len(leftPartitions) <= len(outerPartitions) && len(leftPartitions) <= len(rightPartitions):
		return leftPartitions
	case len(outerPartitions) <= len(leftPartitions) && len(outerPartitions) <= len(rightPartitions):
		return outerPartitions
	default:
		return rightPartitions
	}
}

func actualPartitionsTCPUDP(t tcpudpTripleSet) []ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP], 0)
	for _, p := range t.Partitions() {
		res = slices.Concat(res, breakTCPUDPTriple(p))
	}
	return res
}

// break multi-cube to regular cube
func breakTCPUDPTriple(t ds.Triple[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet]) []ds.Triple[*netset.IPBlock,
	*netset.IPBlock, netp.TCPUDP] {
	res := make([]ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP], 0)

	dstCidrs := t.S2.SplitToCidrs()
	tcpudpTriples := t.S3.Partitions()

	for _, src := range t.S1.SplitToCidrs() {
		for _, dst := range dstCidrs {
			for _, protocolTriple := range tcpudpTriples {
				srcPorts := protocolTriple.S2.Intervals()
				dstPorts := protocolTriple.S3.Intervals()
				for _, protocol := range protocolTriple.S1.Elements() {
					for _, srcPorts := range srcPorts {
						for _, dstPorts := range dstPorts {
							p, _ := netp.NewTCPUDP(protocol == 1, int(srcPorts.Start()), int(srcPorts.End()), int(dstPorts.Start()), int(dstPorts.End()))
							res = append(res, ds.Triple[*netset.IPBlock, *netset.IPBlock, netp.TCPUDP]{S1: src, S2: dst, S3: p})
						}
					}
				}
			}
		}
	}
	return res
}
