/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package io

import (
	"errors"
	"fmt"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func WriteSG(collection *ir.SGCollection, vpc string) ([][]string, error) {
	res := make([][]string, 0)
	for _, vpcName := range collection.VpcNames() {
		if vpc != vpcName && vpc != "" {
			continue
		}
		for _, sgName := range collection.SortedSGNames(vpcName) {
			sgTable, err := makeSGTable(collection.SGs[vpcName][sgName], sgName)
			if err != nil {
				return nil, err
			}
			res = append(res, sgTable...)
		}
	}
	return res, nil
}

func SGHeader() [][]string {
	return [][]string{{
		"SG",
		"Direction",
		"Remote type",
		"Remote",
		"Protocol",
		"Protocol params",
		"Description",
	}}
}

func makeSGTable(t *ir.SG, sgName ir.SGName) ([][]string, error) {
	rules := t.AllRules()
	rows := make([][]string, len(rules))
	for i, rule := range rules {
		sgRow, err := makeSGRow(rule, sgName)
		if err != nil {
			return nil, err
		}
		rows[i] = sgRow
	}
	return rows, nil
}

func makeSGRow(rule *ir.SGRule, sgName ir.SGName) ([]string, error) {
	remoteType, err1 := sgRemoteType(rule.Remote)
	remote, err2 := sgRemote(rule.Remote)
	protocolParams, err3 := printProtocolParams(rule.Protocol)
	if err := errors.Join(err1, err2, err3); err != nil {
		return nil, err
	}

	return []string{
		string(sgName),
		direction(rule.Direction),
		remoteType,
		remote,
		printProtocolName(rule.Protocol),
		protocolParams,
		rule.Explanation,
	}, nil
}

func sgRemoteType(t ir.RemoteType) (string, error) {
	switch p := t.(type) {
	case *netset.IPBlock:
		if ipString := p.ToIPAddressString(); ipString != "" { // single IP address
			return "IP address", nil
		}
		return "CIDR block", nil
	case ir.SGName:
		return "Security group", nil
	}
	return "", fmt.Errorf("impossible remote type %T", t)
}

func sgRemote(r ir.RemoteType) (string, error) {
	switch tr := r.(type) {
	case *netset.IPBlock:
		s := tr.String()
		if s == netset.CidrAll {
			return "Any IP", nil
		}
		return s, nil
	case ir.SGName:
		return tr.String(), nil
	}
	return "", fmt.Errorf("impossible remote %v (%T)", r, r)
}

func printProtocolParams(protocol netp.Protocol) (string, error) {
	switch p := protocol.(type) {
	case netp.ICMP:
		return printICMPTypeCode(protocol), nil
	case netp.TCPUDP:
		return printPorts(p.DstPorts()), nil
	case netp.AnyProtocol:
		return "", nil
	}
	return "", fmt.Errorf("impossible protocol %T", protocol)
}
