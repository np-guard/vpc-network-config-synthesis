/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type (
	sgOptimizer struct {
		sgCollection *ir.SGCollection
		sgName       ir.SGName
		sgVPC        *string
	}

	sgRuleGroups struct {
		sgRemoteRules *sgRulesPerProtocol
		ipRemoteRules *sgRulesPerProtocol
	}

	sgRulesPerProtocol struct {
		tcp  []*ir.SGRule
		udp  []*ir.SGRule
		icmp []*ir.SGRule
		all  []*ir.SGRule
	}

	sgCubesPerProtocol struct {
		tcp  map[ir.SGName]*netset.PortSet
		udp  map[ir.SGName]*netset.PortSet
		icmp map[ir.SGName]*netset.ICMPSet
		all  []ir.SGName
	}

	ipCubesPerProtocol struct {
		tcp  []ds.Pair[*netset.IPBlock, *netset.PortSet]
		udp  []ds.Pair[*netset.IPBlock, *netset.PortSet]
		icmp []ds.Pair[*netset.IPBlock, *netset.ICMPSet]
		all  *netset.IPBlock
	}
)

func NewSGOptimizer(collection ir.Collection, sgName string) optimize.Optimizer {
	components := ir.ScopingComponents(sgName)
	if len(components) == 1 {
		return &sgOptimizer{sgCollection: collection.(*ir.SGCollection), sgName: ir.SGName(sgName), sgVPC: nil}
	}
	return &sgOptimizer{sgCollection: collection.(*ir.SGCollection), sgName: ir.SGName(components[1]), sgVPC: &components[0]}
}

// Optimize attempts to reduce the number of SG rules
// if -n was supplied, it will attempt to reduce the number of rules only in the requested SG
// otherwise, it will attempt to reduce the number of rules in all SGs
func (s *sgOptimizer) Optimize() (ir.Collection, error) {
	if s.sgName != "" {
		for _, vpcName := range utils.SortedMapKeys(s.sgCollection.SGs) {
			if _, ok := s.sgCollection.SGs[vpcName][s.sgName]; ok {
				s.optimizeSG(vpcName, s.sgName)
				return s.sgCollection, nil
			}
		}
		return nil, fmt.Errorf("could no find %s sg", s.sgName)
	}

	for _, vpcName := range utils.SortedMapKeys(s.sgCollection.SGs) {
		for _, sgName := range utils.SortedMapKeys(s.sgCollection.SGs[vpcName]) {
			s.optimizeSG(vpcName, sgName)
		}
	}
	return s.sgCollection, nil
}

// optimizeSG attempts to reduce the number of SG rules
// the algorithm attempts to reduce both inbound and outbound rules separately
// A message is printed to the log at the end of the algorithm
func (s *sgOptimizer) optimizeSG(vpcName string, sgName ir.SGName) {
	sg := s.sgCollection.SGs[vpcName][sgName]
	reducedRules := 0

	// reduce inbound rules first
	for l, rules := range sg.InboundRules {
		local, _ := netset.IPBlockFromCidr(l)
		newInboundRules := s.reduceSGRules(rules, ir.Inbound, local)
		if len(rules) > len(newInboundRules) {
			reducedRules += len(rules) - len(newInboundRules)
			sg.InboundRules[l] = newInboundRules
		}
	}

	// reduce outbound rules second
	for l, rules := range sg.OutboundRules {
		local, _ := netset.IPBlockFromCidr(l)
		newOutboundRules := s.reduceSGRules(rules, ir.Outbound, local)
		if len(rules) > len(newOutboundRules) {
			reducedRules += len(rules) - len(newOutboundRules)
			sg.OutboundRules[l] = newOutboundRules
		}
	}

	// print a message to the log
	switch {
	case reducedRules == 0:
		log.Printf("no rules were reduced in sg %s\n", string(sgName))
	case reducedRules == 1:
		log.Printf("1 rule was reduced in sg %s\n", string(sgName))
	default:
		log.Printf("%d rules were reduced in sg %s\n", reducedRules, string(sgName))
	}
}

// reduceSGRules attempts to reduce the number of rules with different remote types separately
func (s *sgOptimizer) reduceSGRules(rules []*ir.SGRule, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	// separate all rules to groups of protocol X remote ([tcp, udp, icmp, protocolAll] X [ip, sg])
	ruleGroups := divideSGRules(rules)

	// rules with SG as a remote
	optimizedRulesToSG := reduceRulesSGRemote(rulesToSGCubes(ruleGroups.sgRemoteRules), direction, l)
	originlRulesToSG := ruleGroups.sgRemoteRules.allRules()
	if len(originlRulesToSG) < len(optimizedRulesToSG) { // failed to reduce number of rules
		optimizedRulesToSG = originlRulesToSG
	}

	// rules with IPBlock as a remote
	optimizedRulesToIPAddrs := reduceRulesIPRemote(rulesToIPCubes(ruleGroups.ipRemoteRules), direction, l)
	originalRulesToIPAddrs := ruleGroups.ipRemoteRules.allRules()
	if len(originalRulesToIPAddrs) < len(optimizedRulesToSG) { // failed to reduce number of rules
		optimizedRulesToIPAddrs = originalRulesToIPAddrs
	}

	return append(optimizedRulesToSG, optimizedRulesToIPAddrs...)
}

func reduceRulesSGRemote(cubes *sgCubesPerProtocol, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	reduceSGCubes(cubes)

	// cubes to SG rules
	tcpRules := tcpudpSGCubesToRules(cubes.tcp, direction, true, l)
	udpRules := tcpudpSGCubesToRules(cubes.udp, direction, false, l)
	icmpRules := icmpSGCubesToRules(cubes.icmp, direction, l)
	allRules := protocolAllCubesToRules(cubes.all, direction, l)

	// return all rules
	return append(tcpRules, append(udpRules, append(icmpRules, allRules...)...)...)
}

func reduceRulesIPRemote(cubes *ipCubesPerProtocol, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	reduceIPCubes(cubes)

	// cubes to SG rules
	tcpRules := tcpudpIPCubesToRules(cubes.tcp, cubes.all, direction, true, l)
	udpRules := tcpudpIPCubesToRules(cubes.udp, cubes.all, direction, false, l)
	icmpRules := icmpIPCubesToRules(cubes.icmp, cubes.all, direction, l)
	allRules := allProtocolIPCubesIPToRules(cubes.all, direction, l)

	// return all rules
	return append(tcpRules, append(udpRules, append(icmpRules, allRules...)...)...)
}

// divide SGCollection to TCP/UDP/ICMP/ProtocolALL X SGRemote/IPAddrs rules
func divideSGRules(rules []*ir.SGRule) *sgRuleGroups {
	rulesToSG := &sgRulesPerProtocol{tcp: make([]*ir.SGRule, 0), udp: make([]*ir.SGRule, 0),
		icmp: make([]*ir.SGRule, 0), all: make([]*ir.SGRule, 0)}
	rulesToIPAddrs := &sgRulesPerProtocol{tcp: make([]*ir.SGRule, 0), udp: make([]*ir.SGRule, 0),
		icmp: make([]*ir.SGRule, 0), all: make([]*ir.SGRule, 0)}

	for _, rule := range rules {
		// TCP rule
		if p, ok := rule.Protocol.(netp.TCPUDP); ok && p.ProtocolString() == "TCP" {
			if _, ok := rule.Remote.(*netset.IPBlock); ok {
				rulesToIPAddrs.tcp = append(rulesToIPAddrs.tcp, rule)
			} else {
				rulesToSG.tcp = append(rulesToSG.tcp, rule)
			}
		}

		// UDP rule
		if p, ok := rule.Protocol.(netp.TCPUDP); ok && p.ProtocolString() == "UDP" {
			if _, ok := rule.Remote.(*netset.IPBlock); ok {
				rulesToIPAddrs.udp = append(rulesToIPAddrs.udp, rule)
			} else {
				rulesToSG.udp = append(rulesToSG.udp, rule)
			}
		}

		// ICMP rule
		if _, ok := rule.Protocol.(netp.ICMP); ok {
			if _, ok := rule.Remote.(*netset.IPBlock); ok {
				rulesToIPAddrs.icmp = append(rulesToIPAddrs.icmp, rule)
			} else {
				rulesToSG.icmp = append(rulesToSG.icmp, rule)
			}
		}

		// all protocol rules
		if _, ok := rule.Protocol.(netp.AnyProtocol); ok {
			if _, ok := rule.Remote.(*netset.IPBlock); ok {
				rulesToIPAddrs.all = append(rulesToIPAddrs.all, rule)
			} else {
				rulesToSG.all = append(rulesToSG.all, rule)
			}
		}
	}
	return &sgRuleGroups{sgRemoteRules: rulesToSG, ipRemoteRules: rulesToIPAddrs}
}

func (s *sgRulesPerProtocol) allRules() []*ir.SGRule {
	return append(s.tcp, append(s.udp, append(s.icmp, s.all...)...)...)
}
