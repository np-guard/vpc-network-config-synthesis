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

// calculate all spans and set them in sgRulesPerProtocol struct
func sgRulesToIPAddrsToSpans(rules *sgRulesPerProtocol) *sgSpansToIPPerProtocol {
	tcpSpan := tcpudpRulesToIPAddrsToPortsSpan(rules.tcp)
	udpSpan := tcpudpRulesToIPAddrsToPortsSpan(rules.udp)
	icmpSpan := icmpRulesToIPAddrsToSpan(rules.icmp)
	allSpan := allProtocolRulesToIPAddrsToSpan(rules.all)
	return &sgSpansToIPPerProtocol{tcp: tcpSpan, udp: udpSpan, icmp: icmpSpan, all: allSpan}
}

// all protocol rules to a span. The span will be splitted to disjoint CIDRs
func allProtocolRulesToIPAddrsToSpan(rules []ir.SGRule) *netset.IPBlock {
	res := netset.NewIPBlock()
	for i := range rules {
		res.Union(rules[i].Remote.(*netset.IPBlock))
	}
	return res
}

// tcp/udp rules (separately) to a span of (IPBlock X protocol ports).
// all IPBlocks are disjoint
func tcpudpRulesToIPAddrsToPortsSpan(rules []ir.SGRule) []ds.Pair[*netset.IPBlock, *interval.CanonicalSet] {
	span := tcpudpMapSpan(rules)
	result := ds.NewProductLeft[*netset.IPBlock, *interval.CanonicalSet]()
	for ipblock, portsSet := range span {
		r := ds.CartesianPairLeft(ipblock, portsSet)
		result = result.Union(r).(*ds.ProductLeft[*netset.IPBlock, *interval.CanonicalSet])
	}
	return sortPartitionsByIPAddrs(result.Partitions())
}

// icmp rules to a span of (IPBlock X icmp set).
// all IPBlocks are disjoint
func icmpRulesToIPAddrsToSpan(rules []ir.SGRule) []ds.Pair[*netset.IPBlock, *netset.ICMPSet] {
	span := icmpMapSpan((rules))
	result := ds.NewProductLeft[*netset.IPBlock, *netset.ICMPSet]()
	for ipblock, icmpSet := range span {
		r := ds.CartesianPairLeft(ipblock, icmpSet)
		result = result.Union(r).(*ds.ProductLeft[*netset.IPBlock, *netset.ICMPSet])
	}
	return sortPartitionsByIPAddrs(result.Partitions())
}

/* ######################## */
/* #### HELP FUNCTIONS #### */
/* ######################## */

// tcp/udp rules to a span in a map format, where the key is the IPBlock and the value contains the protocol ports
// all ipblocks are disjoint
func tcpudpMapSpan(rules []ir.SGRule) map[*netset.IPBlock]*interval.CanonicalSet {
	span := make(map[*netset.IPBlock]*interval.CanonicalSet, 0) // all keys are disjoint
	for i := range rules {
		portsSet := rules[i].Protocol.(netp.TCPUDP).DstPorts().ToSet() // already checked
		ruleIP := rules[i].Remote.(*netset.IPBlock)                    // already checked
		span = updateSpan(span, portsSet, ruleIP)
	}
	return span
}

// icmp rules to a span in a map format, where the key is the IPBlock and the value contains the icmp set
// all ipblocks are disjoint
func icmpMapSpan(rules []ir.SGRule) map[*netset.IPBlock]*netset.ICMPSet {
	span := make(map[*netset.IPBlock]*netset.ICMPSet, 0) // all keys are disjoint
	for i := range rules {
		icmpSet := netset.NewICMPSet(rules[i].Protocol.(netp.ICMP)) // already checked
		ruleIP := rules[i].Remote.(*netset.IPBlock)                 // already checked
		span = updateSpan(span, icmpSet, ruleIP)
	}
	return span
}

// updateSpan gets the current span, and a rule details (IPBlock and a protocol set)
// if the IPBlock is already in the map, the new protocol set will be unioned with the existing one
// otherwise the rule will be added in the `addRuleToSpan` function.
func updateSpan[T ds.Set[T]](span map[*netset.IPBlock]T, ruleSet T, ruleIP *netset.IPBlock) map[*netset.IPBlock]T {
	if protocolSet, ok := span[ruleIP]; ok {
		span[ruleIP] = protocolSet.Union(ruleSet)
		return span
	}
	span, newMap := addRuleToSpan(span, ruleIP, ruleSet)
	return utils.MergeSetMaps(span, newMap)
}

// any IPblock that overlaps with the new ipblock:
//  1. will be deleted from the new map
//  2. will be splitted into two parts: the part overlapping with the new ipblock and the part that is not
//     a. the overlapping part will enter the new map, where the existing set will be unioned with the new set
//     b. the non overlapping part will enter the new map with the same value he had.
func addRuleToSpan[T ds.Set[T]](span map[*netset.IPBlock]T, ruleIP *netset.IPBlock, ruleSet T) (s, res map[*netset.IPBlock]T) {
	res = make(map[*netset.IPBlock]T, 0)
	for ipblock := range span {
		if !ipblock.Overlap(ruleIP) {
			continue
		}
		overlappingIPs := ruleIP.Subtract(ipblock)
		for _, ip := range overlappingIPs.Split() {
			res[ip] = span[ipblock].Copy().Union(ruleSet)
		}
		notOverlappingIPs := ipblock.Subtract(overlappingIPs)
		for _, ip := range notOverlappingIPs.Split() {
			res[ip] = span[ipblock].Copy()
		}
		delete(span, ipblock)
	}
	return span, res
}
