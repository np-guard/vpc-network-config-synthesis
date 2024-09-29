/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

type Packet struct {
	Src, Dst    *netset.IPBlock
	Protocol    netp.Protocol
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
func makeDenyInternal() []*ACLRule {
	localCidrs := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} // https://datatracker.ietf.org/doc/html/rfc1918#section-3
	localCidrsIPBlocks, _ := netset.IPBlockFromCidrList(localCidrs)
	localCidrsList := localCidrsIPBlocks.Split() // should be splitted to CIDRs
	var denyInternal []*ACLRule
	for i, localCidrSrc := range localCidrsList {
		for j, localCidrDst := range localCidrsList {
			explanation := fmt.Sprintf("Deny other internal communication; see rfc1918#3; item %v,%v", i, j)
			denyInternal = append(denyInternal,
				packetACLRule(&Packet{Src: localCidrSrc, Dst: localCidrDst, Protocol: netp.AnyProtocol{}, Explanation: explanation}, Outbound, Deny),
				packetACLRule(&Packet{Src: localCidrDst, Dst: localCidrSrc, Protocol: netp.AnyProtocol{}, Explanation: explanation}, Inbound, Deny),
			)
		}
	}
	return denyInternal
}

func DenyAllSend(subnetName ID, cidr *netset.IPBlock) *ACLRule {
	explanation := DenyAllExplanation(subnetName, cidr)
	return packetACLRule(&Packet{Src: cidr, Dst: netset.GetCidrAll(), Protocol: netp.AnyProtocol{}, Explanation: explanation}, Outbound, Deny)
}

func DenyAllReceive(subnetName ID, cidr *netset.IPBlock) *ACLRule {
	explanation := DenyAllExplanation(subnetName, cidr)
	return packetACLRule(&Packet{Src: netset.GetCidrAll(), Dst: cidr, Protocol: netp.AnyProtocol{}, Explanation: explanation}, Inbound, Deny)
}

func DenyAllExplanation(subnetName ID, cidr *netset.IPBlock) string {
	return fmt.Sprintf("Deny all communication; subnet %s[%s] does not have required connections", subnetName, cidr.String())
}
