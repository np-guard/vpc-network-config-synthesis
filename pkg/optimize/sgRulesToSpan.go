/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

// Rules with remote SG
func sgRulesToSGToSpans(rules *sgRulesPerProtocol) *sgRulesToSGSpans {
	tcpSpan := tcpudpRulesToSGToPortsSpan(rules.tcp)
	udpSpan := tcpudpRulesToSGToPortsSpan(rules.udp)
	icmpSpan := icmpRulesToSGToSpan(rules.icmp)
	allSpan := allProtocolRulesToSGToSpan(rules.all)
	return &sgRulesToSGSpans{tcp: tcpSpan, udp: udpSpan, icmp: icmpSpan, all: allSpan}
}

func tcpudpRulesToSGToPortsSpan(rules []ir.SGRule) map[ir.SGName]*interval.CanonicalSet {
	result := make(map[ir.SGName]*interval.CanonicalSet)
	for i := range rules {
		p := rules[i].Protocol.(netp.TCPUDP)  // already checked
		remote := rules[i].Remote.(ir.SGName) // already checked
		if result[remote] == nil {
			result[remote] = interval.NewCanonicalSet()
		}
		result[remote].AddInterval(p.DstPorts())
	}
	return result
}

func icmpRulesToSGToSpan(rules []ir.SGRule) map[ir.SGName]*netset.ICMPSet {
	result := make(map[ir.SGName]*netset.ICMPSet)
	for i := range rules {
		p := rules[i].Protocol.(netp.ICMP)    // already checked
		remote := rules[i].Remote.(ir.SGName) // already checked
		if result[remote] == nil {
			result[remote] = netset.EmptyICMPSet()
		}
		result[remote].Union(netset.NewICMPSet(p))
	}
	return result
}

func allProtocolRulesToSGToSpan(rules []ir.SGRule) []*ir.SGName {
	result := make(map[ir.SGName]struct{})
	for i := range rules {
		remote := rules[i].Remote.(ir.SGName)
		result[remote] = struct{}{}
	}
	return utils.ToPtrSlice(utils.SortedMapKeys(result))
}

// Rules with IPAddrs remote
func sgRulesToIPAddrsToSpans(rules *sgRulesPerProtocol) *sgRulesToIPAddrsSpans {
	tcpSpan := tcpudpRulesToIPAddrsToPortsSpan(rules.tcp)
	udpSpan := tcpudpRulesToIPAddrsToPortsSpan(rules.udp)
	icmpSpan := icmpRulesToIPAddrsToSpan(rules.icmp)
	allSpan := allProtocolRulesToIPAddrsToSpan(rules.all)
	return &sgRulesToIPAddrsSpans{tcp: tcpSpan, udp: udpSpan, icmp: icmpSpan, all: allSpan}
}

// converts []ir.SGRule (where all rules or either TCP/UDP but not both) to a span of (IPBlock X ports).
// all IPBlocks are disjoint.
func tcpudpRulesToIPAddrsToPortsSpan(rules []ir.SGRule) []ds.Pair[*netset.IPBlock, *interval.CanonicalSet] {
	span := tcpudpMapSpan(rules)
	result := ds.NewProductLeft[*netset.IPBlock, *interval.CanonicalSet]()
	for ipblock, portsSet := range span {
		r := ds.CartesianPairLeft(ipblock, portsSet)
		result = result.Union(r).(*ds.ProductLeft[*netset.IPBlock, *interval.CanonicalSet])
	}
	return sortPartitionsByIPAddrs(result.Partitions())
}

func icmpRulesToIPAddrsToSpan(rules []ir.SGRule) []ds.Pair[*netset.IPBlock, *netset.ICMPSet] {
	span := icmpMapSpan((rules))
	result := ds.NewProductLeft[*netset.IPBlock, *netset.ICMPSet]()
	for ipblock, icmpSet := range span {
		r := ds.CartesianPairLeft(ipblock, icmpSet)
		result = result.Union(r).(*ds.ProductLeft[*netset.IPBlock, *netset.ICMPSet])
	}
	return sortPartitionsByIPAddrs(result.Partitions())
}

func allProtocolRulesToIPAddrsToSpan(rules []ir.SGRule) []*netset.IPBlock {
	res := netset.NewIPBlock()
	for i := range rules {
		res.Union(rules[i].Remote.(*netset.IPBlock))
	}
	return sortIPBlockSlice(res.Split())
}

// Help functions
func tcpudpMapSpan(rules []ir.SGRule) map[*netset.IPBlock]*interval.CanonicalSet {
	span := make(map[*netset.IPBlock]*interval.CanonicalSet, 0) // all keys are disjoint
	for i := range rules {
		p := rules[i].Protocol.(netp.TCPUDP) // already checked
		portsSet := p.DstPorts().ToSet()
		ruleIP := rules[i].Remote.(*netset.IPBlock) // already checked

		if ports, ok := span[ruleIP]; ok {
			span[ruleIP] = ports.Union(portsSet)
			continue
		}

		newMap := make(map[*netset.IPBlock]*interval.CanonicalSet, 0)
		for ipblock := range span {
			if ipblock.Overlap(ruleIP) {
				overlappingIPs := ruleIP.Subtract(ipblock)
				for _, ip := range overlappingIPs.Split() {
					newMap[ip] = span[ipblock].Copy().Union(portsSet)
				}
				notOverlappingIPs := ipblock.Subtract(overlappingIPs)
				for _, ip := range notOverlappingIPs.Split() {
					newMap[ip] = span[ipblock].Copy()
				}
			}
		}
		span = mergeTcpudpMaps(span, newMap)
	}
	return span
}

func icmpMapSpan(rules []ir.SGRule) map[*netset.IPBlock]*netset.ICMPSet {
	span := make(map[*netset.IPBlock]*netset.ICMPSet, 0) // all keys are disjoint
	for i := range rules {
		icmpSet := netset.NewICMPSet(rules[i].Protocol.(netp.ICMP)) // already checked
		ruleIP := rules[i].Remote.(*netset.IPBlock)                 // already checked

		if ports, ok := span[ruleIP]; ok {
			span[ruleIP] = ports.Union(icmpSet)
			continue
		}

		newMap := make(map[*netset.IPBlock]*netset.ICMPSet, 0)
		for ipblock := range span {
			if ipblock.Overlap(ruleIP) {
				overlappingIPs := ruleIP.Subtract(ipblock)
				for _, ip := range overlappingIPs.Split() {
					newMap[ip] = span[ipblock].Copy().Union(icmpSet)
				}
				notOverlappingIPs := ipblock.Subtract(overlappingIPs)
				for _, ip := range notOverlappingIPs.Split() {
					newMap[ip] = span[ipblock].Copy()
				}
			}
		}
		span = mergeIcmpMaps(span, newMap)
	}
	return span
}

func mergeTcpudpMaps(m1, m2 map[*netset.IPBlock]*interval.CanonicalSet) map[*netset.IPBlock]*interval.CanonicalSet {
	for ipblock, portsSet := range m2 {
		m1[ipblock] = portsSet.Copy()
	}
	return m1
}

func mergeIcmpMaps(m1, m2 map[*netset.IPBlock]*netset.ICMPSet) map[*netset.IPBlock]*netset.ICMPSet {
	for ipblock, portsSet := range m2 {
		m1[ipblock] = portsSet.Copy()
	}
	return m1
}
