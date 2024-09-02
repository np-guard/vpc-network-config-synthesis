/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package csvio

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

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
		sGRemoteType(rule.Remote),
		sgRemote(rule.Remote),
		printProtocolName(rule.Protocol),
		printProtocolParams(rule.Protocol, rule.Direction == ir.Inbound),
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

func sGPort(p interval.Interval) string {
	switch {
	case p.Start() == netp.MinPort && p.End() == netp.MaxPort:
		return "any port"
	default:
		return fmt.Sprintf("ports %v-%v", p.Start(), p.End())
	}
}

func sGRemoteType(t ir.RemoteType) string {
	switch t := t.(type) {
	case *netset.IPBlock:
		if ipString := t.ToIPAddressString(); ipString != "" {
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

func printProtocolParams(protocol netp.Protocol, isSource bool) string {
	switch p := protocol.(type) {
	case netp.ICMP:
		return printICMPTypeCode(protocol)
	case netp.TCPUDP:
		var r interval.Interval
		if isSource {
			r = p.SrcPorts()
		} else {
			r = p.DstPorts()
		}
		return sGPort(r)
	case netp.AnyProtocol:
		return ""
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}
