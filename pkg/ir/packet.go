/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"

	"github.com/np-guard/models/pkg/ipblock"
)

type Packet struct {
	Src, Dst    *ipblock.IPBlock
	Protocol    Protocol
	Explanation string
}

func AllowSend(packet *Packet) *ACLRule {
	return packetACLRule(packet, Outbound, Allow)
}

func AllowReceive(packet *Packet) *ACLRule {
	return packetACLRule(packet, Inbound, Allow)
}

func packetACLRule(packet *Packet, direction Direction, action Action) *ACLRule {
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
	localCidrs := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} // https://datatracker.ietf.org/doc/html/rfc1918#section-3
	localCidrsIPBlocks, _ := ipblock.FromCidrList(localCidrs)
	localCidrsList := localCidrsIPBlocks.Split()
	var denyInternal []ACLRule
	for i, anyLocalCidrSrc := range localCidrsList {
		for j, anyLocalCidrDst := range localCidrsList {
			explanation := fmt.Sprintf("Deny other internal communication; see rfc1918#3; item %v,%v", i, j)
			denyInternal = append(denyInternal,
				*packetACLRule(&Packet{Src: anyLocalCidrSrc, Dst: anyLocalCidrDst, Protocol: AnyProtocol{}, Explanation: explanation}, Outbound, Deny),
				*packetACLRule(&Packet{Src: anyLocalCidrDst, Dst: anyLocalCidrSrc, Protocol: AnyProtocol{}, Explanation: explanation}, Inbound, Deny),
			)
		}
	}
	return denyInternal
}

func DenyAllSend(subnetName ID, cidr *ipblock.IPBlock) *ACLRule {
	explanation := DenyAllExplanation(subnetName, cidr)
	return packetACLRule(&Packet{Src: cidr, Dst: ipblock.GetCidrAll(), Protocol: AnyProtocol{}, Explanation: explanation}, Outbound, Deny)
}

func DenyAllReceive(subnetName ID, cidr *ipblock.IPBlock) *ACLRule {
	explanation := DenyAllExplanation(subnetName, cidr)
	return packetACLRule(&Packet{Src: ipblock.GetCidrAll(), Dst: cidr, Protocol: AnyProtocol{}, Explanation: explanation}, Inbound, Deny)
}

func DenyAllExplanation(subnetName ID, cidr *ipblock.IPBlock) string {
	return fmt.Sprintf("Deny all communication; subnet %s[%s] does not have required connections", subnetName, cidr.String())
}
