/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package csvio

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func (w *Writer) WriteSG(collection *ir.SGCollection, vpc string) error {
	if err := w.w.WriteAll(sgHeader()); err != nil {
		return err
	}
	for _, sgName := range collection.SortedSGNames(vpc) {
		vpcName := ir.VpcFromScopedResource(string(sgName))
		if err := w.w.WriteAll(makeSGTable(collection.SGs[vpcName][sgName], sgName)); err != nil {
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

func makeSGRow(rule *ir.SGRule, sgName ir.SGName) []string {
	return []string{
		string(sgName),
		direction(rule.Direction),
		sgRemoteType(rule.Remote),
		sgRemote(rule.Remote),
		printProtocolName(rule.Protocol),
		printProtocolParams(rule.Protocol),
		rule.Explanation,
	}
}

func makeSGTable(t *ir.SG, sgName ir.SGName) [][]string {
	rules := t.AllRules()
	rows := make([][]string, len(rules))
	for i := range rules {
		rows[i] = makeSGRow(&rules[i], sgName)
	}
	return rows
}

func sgPort(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return "any port"
	default:
		return fmt.Sprintf("Ports %v-%v", p.Min, p.Max)
	}
}

func sgRemoteType(t ir.RemoteType) string {
	switch tr := t.(type) {
	case *ipblock.IPBlock:
		if ir.IsIPAddress(tr) {
			return "IP address"
		}
		return "CIDR block"
	case ir.SGName:
		return "Security group"
	}
	log.Fatalf("impossible remote type %T", t)
	return ""
}

func sgRemote(r ir.RemoteType) string {
	switch tr := r.(type) {
	case *ipblock.IPBlock:
		s := tr.String()
		if s == ipblock.CidrAll {
			return "Any IP"
		}
		return s
	case ir.SGName:
		return tr.String()
	default:
		log.Panicf("Impossible remote %v (%T)", r, r)
	}
	return ""
}

func printProtocolParams(protocol ir.Protocol) string {
	switch p := protocol.(type) {
	case ir.ICMP:
		return printICMPTypeCode(protocol)
	case ir.TCPUDP:
		return sgPort(p.PortRangePair.DstPort)
	case ir.AnyProtocol:
		return ""
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}
