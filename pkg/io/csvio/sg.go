/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package csvio

import (
	"errors"
	"fmt"

	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func (w *Writer) WriteSG(collection *ir.SGCollection, vpc string) error {
	if err := w.w.WriteAll(sgHeader()); err != nil {
		return err
	}
	for _, sgName := range collection.SortedSGNames(vpc) {
		vpcName := ir.VpcFromScopedResource(string(sgName))
		sgTable, err := makeSGTable(collection.SGs[vpcName][sgName], sgName)
		if err != nil {
			return err
		}
		if err := w.w.WriteAll(sgTable); err != nil {
			return err
		}
	}
	return nil
}

func sgHeader() [][]string {
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

func sgPort(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return "any port"
	default:
		return fmt.Sprintf("Ports %v-%v", p.Min, p.Max)
	}
}

func sgRemoteType(t ir.RemoteType) (string, error) {
	switch tr := t.(type) {
	case *ipblock.IPBlock:
		if ir.IsIPAddress(tr) {
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
	case *ipblock.IPBlock:
		s := tr.String()
		if s == ipblock.CidrAll {
			return "Any IP", nil
		}
		return s, nil
	case ir.SGName:
		return tr.String(), nil
	}
	return "", fmt.Errorf("impossible remote %v (%T)", r, r)
}

func printProtocolParams(protocol ir.Protocol) (string, error) {
	switch p := protocol.(type) {
	case ir.ICMP:
		return printICMPTypeCode(protocol), nil
	case ir.TCPUDP:
		return sgPort(p.PortRangePair.DstPort), nil
	case ir.AnyProtocol:
		return "", nil
	}
	return "", fmt.Errorf("impossible protocol %T", protocol)
}
