/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package io

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func WriteACL(collection *ir.ACLCollection, vpc string) ([][]string, error) {
	res := make([][]string, 0)
	for _, subnet := range collection.SortedACLSubnets(vpc) {
		vpcName := ir.VpcFromScopedResource(subnet)
		aclTable, err := makeACLTable(collection.ACLs[vpcName][subnet], subnet)
		if err != nil {
			return nil, err
		}
		res = slices.Concat(res, aclTable)
	}
	return res, nil
}

func ACLHeader() [][]string {
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

func makeACLRow(priority int, rule *ir.ACLRule, aclName, subnet string) ([]string, error) {
	src, err1 := printIP(rule.Source, rule.Protocol, true)
	dst, err2 := printIP(rule.Destination, rule.Protocol, false)
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
		src,
		dst,
		printICMPTypeCode(rule.Protocol),
		rule.Explanation,
	}, nil
}

func action(a ir.Action) string {
	if a == ir.Deny {
		return "Deny"
	}
	return "Allow"
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
		r := p.DstPorts()
		if isSource {
			r = p.SrcPorts()
		}
		return fmt.Sprintf("%v, %v", ipString, printPorts(r)), nil
	case netp.AnyProtocol:
		return ipString, nil
	}
	return "", fmt.Errorf("impossible protocol %T", protocol)
}
