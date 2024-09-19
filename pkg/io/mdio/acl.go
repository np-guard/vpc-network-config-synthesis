/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package mdio

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Write prints an entire collection of acls as a single MD table.
func (w *Writer) WriteACL(collection *ir.ACLCollection, vpc string) error {
	if err := w.writeAll(aclHeader()); err != nil {
		return err
	}
	for _, subnet := range collection.SortedACLSubnets(vpc) {
		vpcName := ir.VpcFromScopedResource(subnet)
		aclTable, err := makeACLTable(collection.ACLs[vpcName][subnet], subnet)
		if err != nil {
			return err
		}
		if err := w.writeAll(aclTable); err != nil {
			return err
		}
	}
	return nil
}

func makeACLTable(t *ir.ACL, subnet string) ([][]string, error) {
	rules := t.Rules()
	rows := make([][]string, len(rules))
	for i := range rules {
		aclRow, err := makeACLRow(i+1, &rules[i], t.Name(), subnet)
		if err != nil {
			return nil, err
		}
		rows[i] = aclRow
	}
	return rows, nil
}

func aclPort(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return "any port" //nolint:goconst // independent decision for SG and ACL
	default:
		return fmt.Sprintf("ports %v-%v", p.Min, p.Max)
	}
}

func action(a ir.Action) string {
	if a == ir.Deny {
		return "Deny"
	}
	return "Allow"
}

func aclHeader() [][]string {
	return [][]string{{
		"",
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
		"",
	}, {
		"",
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		"",
	}}
}

func makeACLRow(priority int, rule *ir.ACLRule, aclName, subnet string) ([]string, error) {
	srcProtocol, err1 := printIP(rule.Source, rule.Protocol, true)
	dstProtocol, err2 := printIP(rule.Destination, rule.Protocol, false)
	if err := errors.Join(err1, err2); err != nil {
		return nil, err
	}

	return []string{
		"",
		aclName,
		subnet,
		direction(rule.Direction),
		strconv.Itoa(priority),
		action(rule.Action),
		printProtocolName(rule.Protocol),
		srcProtocol,
		dstProtocol,
		printICMPTypeCode(rule.Protocol),
		rule.Explanation,
		"",
	}, nil
}

func printIP(ip *ipblock.IPBlock, protocol ir.Protocol, isSource bool) (string, error) {
	ipString := ip.String()
	if ip.Equal(ipblock.GetCidrAll()) {
		ipString = "Any IP" //nolint:goconst // independent decision for SG and ACL
	}
	switch p := protocol.(type) {
	case ir.ICMP:
		return ipString, nil
	case ir.TCPUDP:
		var r ir.PortRange
		if isSource {
			r = p.PortRangePair.SrcPort
		} else {
			r = p.PortRangePair.DstPort
		}
		return fmt.Sprintf("%v, %v", ipString, aclPort(r)), nil
	case ir.AnyProtocol:
		return ipString, nil
	}
	return "", fmt.Errorf("impossible protocol %T", protocol)
}
