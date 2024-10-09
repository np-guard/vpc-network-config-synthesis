/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

// SG remote
func rulesToSGCubes(rules *rulesPerProtocol) *sgCubesPerProtocol {
	tcpSpan := tcpudpRulesSGCubes(rules.tcp)
	udpSpan := tcpudpRulesSGCubes(rules.udp)
	icmpSpan := icmpRulesSGCubes(rules.icmp)
	allSpan := allProtocolRulesToSGCubes(rules.all)
	return &sgCubesPerProtocol{tcp: tcpSpan, udp: udpSpan, icmp: icmpSpan, all: allSpan}
}

// all protocol rules to cubes
func allProtocolRulesToSGCubes(rules []*ir.SGRule) []ir.SGName {
	result := make(map[ir.SGName]struct{})
	for i := range rules {
		remote := rules[i].Remote.(ir.SGName) // already checked
		result[remote] = struct{}{}
	}
	return utils.SortedMapKeys(result)
}

// tcp/udp rules to cubes -- map where the key is the SG name and the value is the protocol ports
func tcpudpRulesSGCubes(rules []*ir.SGRule) map[ir.SGName]*netset.PortSet {
	result := make(map[ir.SGName]*netset.PortSet)
	for _, rule := range rules {
		p := rule.Protocol.(netp.TCPUDP)  // already checked
		remote := rule.Remote.(ir.SGName) // already checked
		if result[remote] == nil {
			result[remote] = interval.NewCanonicalSet()
		}
		result[remote].AddInterval(p.DstPorts())
	}
	return result
}

// icmp rules to cubes -- map where the key is the SG name and the value is icmpset
func icmpRulesSGCubes(rules []*ir.SGRule) map[ir.SGName]*netset.ICMPSet {
	result := make(map[ir.SGName]*netset.ICMPSet)
	for _, rule := range rules {
		p := rule.Protocol.(netp.ICMP)    // already checked
		remote := rule.Remote.(ir.SGName) // already checked
		if result[remote] == nil {
			result[remote] = netset.EmptyICMPSet()
		}
		icmpSet := optimize.IcmpRuleToIcmpSet(p)
		result[remote] = result[remote].Union(icmpSet)
	}
	return result
}

// IP remote
func rulesToIPCubes(rules *rulesPerProtocol) *ipCubesPerProtocol {
	tcpCubes := tcpudpRulesToIPCubes(rules.tcp)
	udpCubes := tcpudpRulesToIPCubes(rules.udp)
	icmpCubes := icmpRulesToIPCubes(rules.icmp)
	allCubes := allProtocolRulesToIPCubes(rules.all)
	return &ipCubesPerProtocol{tcp: tcpCubes, udp: udpCubes, icmp: icmpCubes, all: allCubes}
}

// all protocol rules to cubes
func allProtocolRulesToIPCubes(rules []*ir.SGRule) *netset.IPBlock {
	res := netset.NewIPBlock()
	for i := range rules {
		res.Union(rules[i].Remote.(*netset.IPBlock))
	}
	return res
}

// tcp/udp rules (separately) to cubes (IPBlock X portset).
func tcpudpRulesToIPCubes(rules []*ir.SGRule) []ds.Pair[*netset.IPBlock, *netset.PortSet] {
	cubes := ds.NewProductLeft[*netset.IPBlock, *netset.PortSet]()
	for _, rule := range rules {
		ipb := rule.Remote.(*netset.IPBlock) // already checked
		p := rule.Protocol.(netp.TCPUDP)     // already checked
		r := ds.CartesianPairLeft(ipb, p.DstPorts().ToSet())
		cubes = cubes.Union(r).(*ds.ProductLeft[*netset.IPBlock, *netset.PortSet])
	}
	return optimize.SortPartitionsByIPAddrs(cubes.Partitions())
}

// icmp rules to cubes (IPBlock X icmp set).
func icmpRulesToIPCubes(rules []*ir.SGRule) []ds.Pair[*netset.IPBlock, *netset.ICMPSet] {
	cubes := ds.NewProductLeft[*netset.IPBlock, *netset.ICMPSet]()
	for _, rule := range rules {
		ipb := rule.Remote.(*netset.IPBlock) // already checked
		p := rule.Protocol.(netp.ICMP)       // already checked
		icmpSet := optimize.IcmpRuleToIcmpSet(p)
		r := ds.CartesianPairLeft(ipb, icmpSet)
		cubes = cubes.Union(r).(*ds.ProductLeft[*netset.IPBlock, *netset.ICMPSet])
	}
	return optimize.SortPartitionsByIPAddrs(cubes.Partitions())
}
