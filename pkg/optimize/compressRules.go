/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"github.com/np-guard/models/pkg/netset"
)

func compressSpansToSG(spans *sgSpansToSGPerProtocol) *sgSpansToSGPerProtocol {
	spans = deleteOtherProtocolIfAllProtocolExists(spans)
	return compressThreeProtocolsToAllProtocol(spans)
}

// delete other protocols rules if all protocol rule exists
func deleteOtherProtocolIfAllProtocolExists(spans *sgSpansToSGPerProtocol) *sgSpansToSGPerProtocol {
	for _, sgName := range spans.all {
		delete(spans.tcp, sgName)
		delete(spans.udp, sgName)
		delete(spans.icmp, sgName)
	}
	return spans
}

// merge tcp, udp and icmp rules into all protocol rule
func compressThreeProtocolsToAllProtocol(spans *sgSpansToSGPerProtocol) *sgSpansToSGPerProtocol {
	for sgName, tcpPorts := range spans.tcp {
		if udpPorts, ok := spans.udp[sgName]; ok {
			if ic, ok := spans.icmp[sgName]; ok {
				if ic.Equal(netset.AllICMPSet()) && allPorts(tcpPorts) && allPorts(udpPorts) { // all tcp&udp ports and all icmp types&codes
					delete(spans.tcp, sgName)
					delete(spans.udp, sgName)
					delete(spans.icmp, sgName)
					spans.all = append(spans.all, sgName)
				}
			}
		}
	}
	return spans
}

// observation: It pays to switch to all protocol rule when we have rules that cover all other protocols
// on exactly the same cidr (only one protocol can exceed).
func compressSpansToIP(span *sgSpansToIPPerProtocol) *sgSpansToIPPerProtocol {
	t := 0
	u := 0
	i := 0

	for t != len(span.tcp) && u != len(span.udp) && i != len(span.icmp) {
		if !allPorts(span.tcp[t].Right) {
			t++
			continue
		}
		if !allPorts(span.udp[u].Right) {
			u++
			continue
		}
		if !span.icmp[i].Right.Equal(netset.AllICMPSet()) {
			i++
			continue
		}

		if span.tcp[t].Left.Equal(span.udp[u].Left) && span.tcp[t].Left.Equal(span.icmp[i].Left) {
			span.all = span.all.Union(span.tcp[t].Left.Copy())
			span.tcp = append(span.tcp[:t], span.tcp[t+1:]...)
			span.udp = append(span.udp[:u], span.udp[u+1:]...)
			span.icmp = append(span.icmp[:i], span.icmp[i+1:]...)
		}
	}

	return span
}
