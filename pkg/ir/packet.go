/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/ipblock"
)

type Packet struct {
	Src, Dst    *ipblock.IPBlock
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
	cidr1, err := ipblock.FromCidr("10.0.0.0/8")
	if err != nil {
		log.Fatal(err)
	}
	cidr2, err := ipblock.FromCidr("172.16.0.0/12")
	if err != nil {
		log.Fatal(err)
	}
	cidr3, err := ipblock.FromCidr("192.168.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	localCidrs := []*ipblock.IPBlock{cidr1, cidr2, cidr3} // https://datatracker.ietf.org/doc/html/rfc1918#section-3
	var denyInternal []ACLRule
	for i, anyLocalCidrSrc := range localCidrs {
		for j, anyLocalCidrDst := range localCidrs {
			explanation := fmt.Sprintf("Deny other internal communication; see rfc1918#3; item %v,%v", i, j)
			denyInternal = append(denyInternal, []ACLRule{
				*packetACLRule(Packet{Src: anyLocalCidrSrc, Dst: anyLocalCidrDst, Protocol: AnyProtocol{}, Explanation: explanation}, Outbound, Deny),
				*packetACLRule(Packet{Src: anyLocalCidrDst, Dst: anyLocalCidrSrc, Protocol: AnyProtocol{}, Explanation: explanation}, Inbound, Deny),
			}...)
		}
	}
	return denyInternal
}
