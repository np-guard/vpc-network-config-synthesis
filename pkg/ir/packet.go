/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
)

type Packet struct {
	Src, Dst    IP
	Protocol    Protocol
	Explanation string
}

func AllowSend(packet Packet) *ACLRule {
	return packetACLRule(packet, Outbound, Allow)
}

func AllowReceive(packet Packet) *ACLRule {
	return packetACLRule(packet, Inbound, Allow)
}

func packetACLRule(packet Packet, direction Direction, action Action) *ACLRule {
	return &ACLRule{
		Action:      action,
		Source:      packet.Src,
		Destination: packet.Dst,
		Direction:   direction,
		Protocol:    packet.Protocol,
		Explanation: packet.Explanation,
	}
}

// makeDenyInternal prevents allowing external communications from accidentally allowing internal communications too
func makeDenyInternal() []ACLRule {
	localIPs := []IP{ // https://datatracker.ietf.org/doc/html/rfc1918#section-3
		{"10.0.0.0/8"},
		{"172.16.0.0/12"},
		{"192.168.0.0/16"},
	}
	var denyInternal []ACLRule
	for i, anyLocalIPSrc := range localIPs {
		for j, anyLocalIPDst := range localIPs {
			explanation := fmt.Sprintf("Deny other internal communication; see rfc1918#3; item %v,%v", i, j)
			denyInternal = append(denyInternal, []ACLRule{
				*packetACLRule(Packet{Src: anyLocalIPSrc, Dst: anyLocalIPDst, Protocol: AnyProtocol{}, Explanation: explanation}, Outbound, Deny),
				*packetACLRule(Packet{Src: anyLocalIPDst, Dst: anyLocalIPSrc, Protocol: AnyProtocol{}, Explanation: explanation}, Inbound, Deny),
			}...)
		}
	}
	return denyInternal
}
