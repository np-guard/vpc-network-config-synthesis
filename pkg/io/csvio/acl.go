/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package csvio

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Write prints an entire collection of acls as a single CSV table.
func (w *Writer) WriteACL(collection *ir.ACLCollection, vpc string) error {
	if err := w.w.WriteAll(aclHeader()); err != nil {
		return err
	}
	for _, subnet := range collection.SortedACLSubnets(vpc) {
		vpcName := ir.VpcFromScopedResource(subnet)
		aclTable, err := makeACLTable(collection.ACLs[vpcName][subnet], subnet)
		if err != nil {
			return err
		}
		if err := w.w.WriteAll(aclTable); err != nil {
			return err
		}
	}
	return nil
}

func makeACLTable(t *ir.ACL, subnet string) ([][]string, error) {
	rules := t.Rules()
	rows := make([][]string, len(rules))
	for i, rule := range rules {
		aclRow, err := makeACLRow(i+1, rule, t.Name(), subnet)
		if err != nil {
			return nil, err
		}
		rows[i] = aclRow
	}
	return rows, nil
}

func aclPort(p interval.Interval) string {
	if p.Equal(netp.AllPorts()) {
		return "any port" //nolint:goconst // independent decision for SG and ACL
	}
	return fmt.Sprintf("ports %v-%v", p.Start(), p.End())
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

func makeACLRow(priority int, rule *ir.ACLRule, aclName, subnet string) ([]string, error) {
	srcProtocol, err1 := printIP(rule.Source, rule.Protocol, true)
	dstProtocol, err2 := printIP(rule.Destination, rule.Protocol, false)
	if errors.Join(err1, err2) != nil {
		return nil, errors.Join(err1, err2)
	}

	return []string{
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
	}, nil
}

func printIP(ip *netset.IPBlock, protocol netp.Protocol, isSource bool) (string, error) {
	ipString := ip.String()
	if ip.Equal(netset.GetCidrAll()) {
		ipString = "Any IP" //nolint:goconst // independent decision for SG and ACL
	}
	switch p := protocol.(type) {
	case netp.ICMP:
		return ipString, nil
	case netp.TCPUDP:
		var r interval.Interval
		if isSource {
			r = p.SrcPorts()
		} else {
			r = p.DstPorts()
		}
		return fmt.Sprintf("%v, %v", ipString, aclPort(r)), nil
	case netp.AnyProtocol:
		return ipString, nil
	}
	return "", fmt.Errorf("impossible protocol %T", protocol)
}
