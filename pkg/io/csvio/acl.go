/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package csvio

import (
	"fmt"
	"log"
	"strconv"

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

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

func ACLPort(p interval.Interval) string {
	switch {
	case p.Start() == netp.MinPort && p.End() == netp.MaxPort:
		return "any port" //nolint:goconst // independent decision for SG and ACL
	default:
		return fmt.Sprintf("ports %v-%v", p.Start(), p.End())
	}
}

// Write prints an entire collection of acls as a single CSV table.
func (w *Writer) WriteSynthACL(collection *ir.ACLCollection, vpc string) error {
	if err := w.w.WriteAll(aclHeader()); err != nil {
		return err
	}
	for _, subnet := range collection.SortedACLSubnets(vpc) {
		vpcName := ir.VpcFromScopedResource(subnet)
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

func printIP(ip *netset.IPBlock, protocol netp.Protocol, isSource bool) string {
	ipString := ip.String()
	if ip.Equal(netset.GetCidrAll()) {
		ipString = "Any IP" //nolint:goconst // independent decision for SG and ACL
	}
	switch p := protocol.(type) {
	case netp.ICMP:
		return ipString
	case netp.TCPUDP:
		var r interval.Interval
		if isSource {
			r = p.SrcPorts()
		} else {
			r = p.DstPorts()
		}
		return fmt.Sprintf("%v, %v", ipString, ACLPort(r))
	case netp.AnyProtocol:
		return ipString
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}
