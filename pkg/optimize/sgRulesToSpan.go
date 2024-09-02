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

func icmpRulesToSGToIcmpSpan(rules []ir.SGRule) map[*ir.SGName]*icmp {
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

// converts []ir.SGRule (where all rules or either TCP/UDP but not both) to a span of (IPBlock X ports)
func tcpudpRulesToIPAddrsToPortsSpan(rules []ir.SGRule) (p []ds.Pair[*netset.IPBlock, *interval.CanonicalSet]) {
	span := ds.NewProductLeft[*netset.IPBlock, *interval.CanonicalSet]()
	for i := range rules {
		p := rules[i].Protocol.(netp.TCPUDP) // already checked
		r := ds.CartesianPairLeft(rules[i].Remote.(*netset.IPBlock), p.DstPorts().ToSet())
		span = span.Union(r).(*ds.ProductLeft[*netset.IPBlock, *interval.CanonicalSet])
	}
	return sortPartitionsByIPAddrs(span.Partitions())
}

func icmpRulesToIPAddrsToIcmpSpan(rules []ir.SGRule) bool {
	return true
}

func allProtocolRulesToIPAddrsToSpan(rules []ir.SGRule) []*netset.IPBlock {
	result := make([]*netset.IPBlock, len(rules))
	for i := range rules {
		result[i] = rules[i].Remote.(*netset.IPBlock)
	}
	return result // should sort !!
}
