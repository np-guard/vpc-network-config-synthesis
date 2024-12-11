/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func aclRulesToCubes(rules []*ir.ACLRule) *aclRulesPerProtocol {
	res := &aclRulesPerProtocol{
		tcpAllow:         ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.PortSet](),
		tcpDeny:          ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.PortSet](),
		udpAllow:         ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.PortSet](),
		udpDeny:          ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.PortSet](),
		icmpAllow:        ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet](),
		icmpDeny:         ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet](),
		anyProtocolAllow: ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock](),
		anyProtocolDeny:  ds.NewProductLeft[*netset.IPBlock, *netset.IPBlock](),
	}

	for _, rule := range rules {
		switch p := rule.Protocol.(type) {
		case netp.TCPUDP:
			tcpudpRuleToCubes(res, rule, p.ProtocolString() == "TCP")
		case netp.ICMP:
			icmpRuleToCubes(res, rule)
		case netp.AnyProtocol:
			allRuleToCubes(res, rule)
		}
	}
	return res
}

func tcpudpRuleToCubes(rules *aclRulesPerProtocol, rule *ir.ACLRule, isTCP bool) {
	_ = ds.CartesianLeftTriple(rule.Source, rule.Destination, rule.Protocol.(netp.TCPUDP).DstPorts().ToSet())

}

func icmpRuleToCubes(rules *aclRulesPerProtocol, rule *ir.ACLRule) {
	t := ds.CartesianLeftTriple(rule.Source, rule.Destination, optimize.IcmpRuleToIcmpSet(rule.Protocol.(netp.ICMP)))
	if rule.Action == ir.Allow {
		r := t.Subtract(rules.icmpDeny)
		rules.icmpAllow = rules.icmpAllow.Union(r)
	} else {
		r := t.Subtract(rules.icmpAllow)
		rules.icmpDeny = rules.icmpDeny.Union(r)
	}
}

func allRuleToCubes(rules *aclRulesPerProtocol, rule *ir.ACLRule) {

}
