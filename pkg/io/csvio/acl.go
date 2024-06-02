/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package csvio

import (
	"fmt"
	"log"
	"strconv"

	"github.com/np-guard/models/pkg/ipblock"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func makeACLTable(t *ir.ACL, subnet string) [][]string {
	rules := t.Rules()
	rows := make([][]string, len(rules))
	for i := range rules {
		rows[i] = makeACLRow(i+1, &rules[i], t.Name(), subnet)
	}
	return rows
}

func ACLPort(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return "any port" //nolint:goconst // independent decision for SG and ACL
	default:
		return fmt.Sprintf("ports %v-%v", p.Min, p.Max)
	}
}

// Write prints an entire collection of acls as a single CSV table.
func (w *Writer) WriteACL(collection *ir.ACLCollection, vpc string) error {
	if err := w.w.WriteAll(aclHeader()); err != nil {
		return err
	}
	var sortedACLs []string
	if vpc == "" {
		sortedACLs = collection.SortedACLSubnets()
	} else {
		sortedACLs = collection.SortedACLSubnetsInVPC(vpc)
	}
	for _, subnet := range sortedACLs {
		vpcName := ir.ScopingComponents(subnet)[0]
		if err := w.w.WriteAll(makeACLTable(collection.ACLs[vpcName][subnet], subnet)); err != nil {
			return err
		}
	}
	return nil
}

func action(a ir.Action) string {
	if a == ir.Deny {
		return "Deny"
	}
	return "Allow"
}

func aclHeader() [][]string {
	return [][]string{{
		"Acl",
		"Subnet",
		"Direction",
		"Rule priority",
		"Allow or deny",
		"Protocol",
		"Source",
		"Destination",
		"Value",
		"Description",
	}}
}

func makeACLRow(priority int, rule *ir.ACLRule, aclName, subnet string) []string {
	return []string{
		aclName,
		subnet,
		direction(rule.Direction),
		strconv.Itoa(priority),
		action(rule.Action),
		printProtocolName(rule.Protocol),
		printIP(rule.Source, rule.Protocol, true),
		printIP(rule.Destination, rule.Protocol, false),
		printICMPTypeCode(rule.Protocol),
		rule.Explanation,
	}
}

func printIP(ip *ipblock.IPBlock, protocol ir.Protocol, isSource bool) string {
	ipString := ip.String()
	if ipString == ipblock.CidrAll {
		ipString = "Any IP"
	}
	switch p := protocol.(type) {
	case ir.ICMP:
		return ipString
	case ir.TCPUDP:
		var r ir.PortRange
		if isSource {
			r = p.PortRangePair.SrcPort
		} else {
			r = p.PortRangePair.DstPort
		}
		return fmt.Sprintf("%v, %v", ipString, ACLPort(r))
	case ir.AnyProtocol:
		return ipString
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}
