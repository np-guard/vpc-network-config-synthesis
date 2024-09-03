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

// Rules with SG remote
func sgRulesToSGToSpans(rules *sgRulesPerProtocol) *sgRulesToSGSpans {
	tcpSpan := tcpudpRulesToSGToPortsSpan(rules.tcp)
	udpSpan := tcpudpRulesToSGToPortsSpan(rules.udp)
	icmpSpan := icmpRulesToSGToSpan(rules.icmp)
	allSpan := allProtocolRulesToSGToSpan(rules.all)
	return &sgRulesToSGSpans{tcp: tcpSpan, udp: udpSpan, icmp: icmpSpan, all: allSpan}
}

func tcpudpRulesToSGToPortsSpan(rules []ir.SGRule) map[*ir.SGName]*interval.CanonicalSet {
	result := make(map[*ir.SGName]*interval.CanonicalSet)
	for i := range rules {
		p := rules[i].Protocol.(netp.TCPUDP)             // already checked
		remote := utils.Ptr(rules[i].Remote.(ir.SGName)) // already checked
		if result[remote] == nil {
			result[remote] = interval.NewCanonicalSet()
		}
		result[remote].AddInterval(p.DstPorts())
	}
	return result
}

func icmpRulesToSGToSpan(rules []ir.SGRule) map[*ir.SGName]*icmp {
	result := make(map[*ir.SGName]*icmp)
	for i := range rules {
		p := rules[i].Protocol.(netp.ICMP)               // already checked
		remote := utils.Ptr(rules[i].Remote.(ir.SGName)) // already checked
		if result[remote] == nil {
			result[remote] = newIcmp()
		}
		result[remote].add(p.TypeCode)
	}
	return result
}

func allProtocolRulesToSGToSpan(rules []ir.SGRule) []*ir.SGName {
	result := make([]*ir.SGName, len(rules))
	for i := range rules {
		result[i] = rules[i].Remote.(*ir.SGName)
	}
	return result
}

// Rules with IPAddrs remote
func sgRulesToIPAddrsToSpans(rules *sgRulesPerProtocol) *sgRulesToIPAddrsSpans {
	tcpSpan := tcpudpRulesToIPAddrsToPortsSpan(rules.tcp)
	udpSpan := tcpudpRulesToIPAddrsToPortsSpan(rules.udp)
	icmpSpan := icmpRulesToIPAddrsToSpan(rules.icmp)
	allSpan := allProtocolRulesToIPAddrsToSpan(rules.all)
	return &sgRulesToIPAddrsSpans{tcp: tcpSpan, udp: udpSpan, icmp: icmpSpan, all: allSpan}
}

// converts []ir.SGRule (where all rules or either TCP/UDP but not both) to a span of (IPBlock X ports)
func tcpudpRulesToIPAddrsToPortsSpan(rules []ir.SGRule) []ds.Pair[*netset.IPBlock, *interval.CanonicalSet] {
	span := ds.NewProductLeft[*netset.IPBlock, *interval.CanonicalSet]()
	for i := range rules {
		p := rules[i].Protocol.(netp.TCPUDP) // already checked
		r := ds.CartesianPairLeft(rules[i].Remote.(*netset.IPBlock), p.DstPorts().ToSet())
		span = span.Union(r).(*ds.ProductLeft[*netset.IPBlock, *interval.CanonicalSet])
	}
	// a := ds.NewLeftTripleSet[*netset.IPBlock, *interval.CanonicalSet, bool]()

	return sortPartitionsByIPAddrs(span.Partitions())
}

func icmpRulesToIPAddrsToSpan(rules []ir.SGRule) []ds.Pair[*netset.IPBlock, *icmp] {
	return sortPartitionsByIPAddrs([]ds.Pair[*netset.IPBlock, *icmp]{})
}

func allProtocolRulesToIPAddrsToSpan(rules []ir.SGRule) []*netset.IPBlock {
	result := make([]*netset.IPBlock, len(rules))
	for i := range rules {
		result[i] = rules[i].Remote.(*netset.IPBlock)
	}
	return sortIPBlockSlice(result)
}
