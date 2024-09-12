/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package mdio

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func (w *Writer) WriteSG(collection *ir.SGCollection, vpc string) error {
	if err := w.writeAll(sgHeader()); err != nil {
		return err
	}
	for _, sgName := range collection.SortedSGNames(vpc) {
		vpcName := ir.VpcFromScopedResource(string(sgName))
		if err := w.writeAll(makeSGTable(collection.SGs[vpcName][sgName], sgName)); err != nil {
			return err
		}
	}
	return nil
}

func sgHeader() [][]string {
	return [][]string{{
		"",
		"SG",
		"Direction",
		"Remote type",
		"Remote",
		"Protocol",
		"Protocol params",
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
		"",
	}}
}

func makeSGRow(rule *ir.SGRule, sgName ir.SGName) []string {
	return []string{
		"",
		string(sgName),
		direction(rule.Direction),
		sgRemoteType(rule.Remote),
		sgRemote(rule.Remote),
		printProtocolName(rule.Protocol),
		printProtocolParams(rule.Protocol),
		rule.Explanation,
		"",
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

func sgPort(p interval.Interval) string {
	switch {
	case p.Start() == netp.MinPort && p.End() == netp.MaxPort:
		return "any port"
	default:
		return fmt.Sprintf("ports %v-%v", p.Start(), p.End())
	}
}

func sgRemoteType(t ir.RemoteType) string {
	switch r := t.(type) {
	case *netset.IPBlock:
		if ipString := r.ToIPAddressString(); ipString != "" {
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
	case *netset.IPBlock:
		s := tr.String()
		if s == netset.CidrAll {
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

func printProtocolParams(protocol netp.Protocol) string {
	switch p := protocol.(type) {
	case netp.ICMP:
		return printICMPTypeCode(protocol)
	case netp.TCPUDP:
		return sgPort(p.DstPorts())
	case netp.AnyProtocol:
		return ""
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}
