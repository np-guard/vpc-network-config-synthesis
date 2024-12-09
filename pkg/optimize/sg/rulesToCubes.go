/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"slices"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

// SG remote
func rulesToSGCubes(rules *rulesPerProtocol) *sgCubesPerProtocol {
	return &sgCubesPerProtocol{tcp: tcpudpRulesSGCubes(rules.tcp),
		udp:  tcpudpRulesSGCubes(rules.udp),
		icmp: icmpRulesSGCubes(rules.icmp),
		anyP: anyProtocolRulesToSGCubes(rules.anyP),
	}
}

// all protocol rules to cubes
func anyProtocolRulesToSGCubes(rules []*ir.SGRule) []ir.SGName {
	res := make([]ir.SGName, len(rules))
	for i := range rules {
		remote := rules[i].Remote.(ir.SGName) // already checked
		res[i] = remote
	}
	return slices.Compact(slices.Sorted(slices.Values(res)))
}

// tcp/udp rules to cubes -- map where the key is the SG name and the value is the protocol ports
func tcpudpRulesSGCubes(rules []*ir.SGRule) map[ir.SGName]*netset.PortSet {
	res := make(map[ir.SGName]*netset.PortSet)
	for _, rule := range rules {
		p := rule.Protocol.(netp.TCPUDP)  // already checked
		remote := rule.Remote.(ir.SGName) // already checked
		if res[remote] == nil {
			res[remote] = interval.NewCanonicalSet()
		}
		res[remote].AddInterval(p.DstPorts())
	}
	return res
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
	return &ipCubesPerProtocol{tcp: tcpudpRulesToIPCubes(rules.tcp),
		udp:  tcpudpRulesToIPCubes(rules.udp),
		icmp: icmpRulesToIPCubes(rules.icmp),
		anyP: anyProtocolRulesToIPCubes(rules.anyP),
	}
}

// any protocol rules to cubes
func anyProtocolRulesToIPCubes(rules []*ir.SGRule) *netset.IPBlock {
	res := netset.NewIPBlock()
	for i := range rules {
		res = res.Union(rules[i].Remote.(*netset.IPBlock))
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
		r := ds.CartesianPairLeft(ipb, optimize.IcmpRuleToIcmpSet(p))
		cubes = cubes.Union(r).(*ds.ProductLeft[*netset.IPBlock, *netset.ICMPSet])
	}
	return optimize.SortPartitionsByIPAddrs(cubes.Partitions())
}
