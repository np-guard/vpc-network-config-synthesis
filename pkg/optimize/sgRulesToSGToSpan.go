/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

// calculate all spans and set them in sgRulesPerProtocol struct
func sgRulesToSGToSpans(rules *sgRulesPerProtocol) *sgSpansToSGPerProtocol {
	tcpSpan := tcpudpRulesToSGToPortsSpan(rules.tcp)
	udpSpan := tcpudpRulesToSGToPortsSpan(rules.udp)
	icmpSpan := icmpRulesToSGToSpan(rules.icmp)
	allSpan := allProtocolRulesToSGToSpan(rules.all)
	return &sgSpansToSGPerProtocol{tcp: tcpSpan, udp: udpSpan, icmp: icmpSpan, all: allSpan}
}

// tcp/udp rules to a span -- map where the key is the SG name and the value is the protocol ports
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

// icmp rules to a span -- map where the key is the SG name and the value is icmp set
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

// all protocol rules to a span of SG names slice
func allProtocolRulesToSGToSpan(rules []ir.SGRule) []*ir.SGName {
	result := make(map[ir.SGName]struct{})
	for i := range rules {
		remote := rules[i].Remote.(ir.SGName)
		result[remote] = struct{}{}
	}
	return utils.ToPtrSlice(utils.SortedMapKeys(result))
}
