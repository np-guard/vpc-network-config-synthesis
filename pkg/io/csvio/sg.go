/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package csvio

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func (w *Writer) WriteSG(collection *ir.SGCollection, vpc string) error {
	if err := w.w.WriteAll(sgHeader()); err != nil {
		return err
	}
	var sortedCollection []ir.SGName
	if vpc == "" {
		sortedCollection = collection.SortedSGNames()
	} else {
		sortedCollection = collection.SortedSGNamesInVPC(vpc)
	}
	for _, sgName := range sortedCollection {
		if err := w.w.WriteAll(makeSGTable(collection.SGs[ir.ScopingComponents(string(sgName))[0]][sgName], sgName)); err != nil {
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
		sGRemoteType(rule.Remote),
		sgRemote(rule.Remote),
		printProtocolName(rule.Protocol),
		printProtocolParams(rule.Protocol, rule.Direction == ir.Inbound),
		rule.Explanation,
	}
}

func makeSGTable(t *ir.SG, sgName ir.SGName) [][]string {
	rules := t.Rules
	rows := make([][]string, len(rules))
	for i := range rules {
		rows[i] = makeSGRow(&rules[i], sgName)
	}
	return rows
}

func sGPort(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return "any port"
	default:
		return fmt.Sprintf("Ports %v-%v", p.Min, p.Max)
	}
}

func sGRemoteType(t ir.RemoteType) string {
	switch t.(type) {
	case ir.IP:
		return "IP address"
	case ir.CIDR:
		return "CIDR block"
	case ir.SGName:
		return "Security group"
	}
	log.Fatalf("impossible remote type %T", t)
	return ""
}

func sgRemote(r ir.RemoteType) string {
	switch tr := r.(type) {
	case ir.IP:
		s := tr.String()
		if s == ir.AnyIP {
			return "Any IP"
		}
	case ir.CIDR:
		s := tr.String()
		if s == ir.AnyCIDR {
			return "Any CIDR"
		}
		return s
	case ir.SGName:
		return tr.String()
	default:
		log.Panicf("Impossible remote %v (%T)", r, r)
	}
	return ""
}

func printProtocolParams(protocol ir.Protocol, isSource bool) string {
	switch p := protocol.(type) {
	case ir.ICMP:
		return printICMPTypeCode(protocol)
	case ir.TCPUDP:
		var r ir.PortRange
		if isSource {
			r = p.PortRangePair.SrcPort
		} else {
			r = p.PortRangePair.DstPort
		}
		return sGPort(r)
	case ir.AnyProtocol:
		return ""
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}
