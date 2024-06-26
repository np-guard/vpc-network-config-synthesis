/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

type Packet struct {
	Src, Dst    *netset.IPBlock
	Protocol    netp.Protocol
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
	cidr1, err := netset.IPBlockFromCidr("10.0.0.0/8")
	if err != nil {
		log.Fatal(err)
	}
	cidr2, err := netset.IPBlockFromCidr("172.16.0.0/12")
	if err != nil {
		log.Fatal(err)
	}
	cidr3, err := netset.IPBlockFromCidr("192.168.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	localCidrs := []*netset.IPBlock{cidr1, cidr2, cidr3} // https://datatracker.ietf.org/doc/html/rfc1918#section-3
	var denyInternal []ACLRule
	for i, anyLocalCidrSrc := range localCidrs {
		for j, anyLocalCidrDst := range localCidrs {
			explanation := fmt.Sprintf("Deny other internal communication; see rfc1918#3; item %v,%v", i, j)
			denyInternal = append(denyInternal, []ACLRule{
				*packetACLRule(Packet{Src: anyLocalCidrSrc, Dst: anyLocalCidrDst, Protocol: netp.AnyProtocol{},
					Explanation: explanation}, Outbound, Deny),
				*packetACLRule(Packet{Src: anyLocalCidrDst, Dst: anyLocalCidrSrc, Protocol: netp.AnyProtocol{},
					Explanation: explanation}, Inbound, Deny),
			}...)
		}
	}
	return denyInternal
}

func DenyAllSend(subnetName ID, cidr *netset.IPBlock) *ACLRule {
	explanation := DenyAllExplanation(subnetName, cidr)
	ACLRule := *packetACLRule(Packet{Src: cidr, Dst: netset.GetCidrAll(), Protocol: netp.AnyProtocol{},
		Explanation: explanation}, Outbound, Deny)
	return &ACLRule
}

func DenyAllReceive(subnetName ID, cidr *netset.IPBlock) *ACLRule {
	explanation := DenyAllExplanation(subnetName, cidr)
	ACLRule := *packetACLRule(Packet{Src: netset.GetCidrAll(), Dst: cidr, Protocol: netp.AnyProtocol{},
		Explanation: explanation}, Inbound, Deny)
	return &ACLRule
}

func DenyAllExplanation(subnetName ID, cidr *netset.IPBlock) string {
	return fmt.Sprintf("Deny all communication; subnet %s[%s] does not have required connections", subnetName, cidr.String())
}
