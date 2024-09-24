/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"log"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/confio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type (
	SGOptimizer struct {
		sgCollection *ir.SGCollection
		sgName       ir.SGName
	}

	sgRuleGroups struct {
		rulesToSG      *sgRulesPerProtocol
		rulesToIPAddrs *sgRulesPerProtocol
	}

	sgRulesPerProtocol struct {
		tcp  []ir.SGRule
		udp  []ir.SGRule
		icmp []ir.SGRule
		all  []ir.SGRule
	}

	sgSpansToSGPerProtocol struct {
		tcp  map[ir.SGName]*interval.CanonicalSet
		udp  map[ir.SGName]*interval.CanonicalSet
		icmp map[ir.SGName]*netset.ICMPSet
		all  []*ir.SGName
	}

	sgSpansToIPPerProtocol struct {
		tcp  []ds.Pair[*netset.IPBlock, *interval.CanonicalSet]
		udp  []ds.Pair[*netset.IPBlock, *interval.CanonicalSet]
		icmp []ds.Pair[*netset.IPBlock, *netset.ICMPSet]
		all  *netset.IPBlock
	}
)

func NewSGOptimizer(sgName string) Optimizer {
	return &SGOptimizer{sgCollection: nil, sgName: ir.SGName(sgName)}
}

// read SGs from config file
func (s *SGOptimizer) ParseCollection(filename string) error {
	c, err := confio.ReadSGs(filename)
	if err != nil {
		return err
	}
	s.sgCollection = c
	return nil
}

// returns a sorted slice of the vpc names
func (s *SGOptimizer) VpcNames() []string {
	return utils.SortedMapKeys(s.sgCollection.SGs)
}

// Optimize attempts to reduce the number of SG rules
// the algorithm attempts to reduce both inbound and outbound rules separately
// A message is printed to the log at the end of the algorithm
func (s *SGOptimizer) Optimize() ir.OptimizeCollection {
	for vpcName := range s.sgCollection.SGs {
		var sg *ir.SG
		var ok bool
		if sg, ok = s.sgCollection.SGs[vpcName][s.sgName]; !ok {
			continue
		}
		reducedRules := 0

		// reduce inbound rules first
		newInboundRules := s.reduceSGRules(sg.InboundRules, ir.Inbound)
		if len(sg.InboundRules) > len(newInboundRules) {
			reducedRules += len(sg.InboundRules) - len(newInboundRules)
			s.sgCollection.SGs[vpcName][s.sgName].InboundRules = newInboundRules
		}

		// reduce outbound rules second
		newOutboundRules := s.reduceSGRules(sg.OutboundRules, ir.Outbound)
		if len(sg.OutboundRules) > len(newOutboundRules) {
			reducedRules += len(sg.OutboundRules) - len(newOutboundRules)
			s.sgCollection.SGs[vpcName][s.sgName].OutboundRules = newOutboundRules
		}

		// print a message to the log
		switch {
		case reducedRules == 0:
			log.Printf("no rules were reduced in sg %s", string(s.sgName))
		case reducedRules == 1:
			log.Printf("1 rule was reduced in sg %s", string(s.sgName))
		default:
			log.Printf("%d rules were reduced in sg %s", reducedRules, string(s.sgName))
		}
	}
	return s.sgCollection
}

// reduceSGRules attempts to reduce the number of rules with different remote types separately
func (s *SGOptimizer) reduceSGRules(rules []ir.SGRule, direction ir.Direction) []ir.SGRule {
	// separate all rules to groups of (protocol X remote)
	ruleGroups := divideSGRules(rules)

	// rules with SG as a remote
	optimizedRulesToSG := reduceSGRulesToSG(sgRulesToSGToSpans(ruleGroups.rulesToSG), direction)
	if len(ruleGroups.rulesToSG.allRules()) < len(optimizedRulesToSG) {
		optimizedRulesToSG = ruleGroups.rulesToSG.allRules()
	}

	// rules with IPBlock as a remote
	optimizedRulesToIPAddrs := reduceSGRulesToIPAddrs(sgRulesToIPAddrsToSpans(ruleGroups.rulesToIPAddrs), direction)
	if len(ruleGroups.rulesToIPAddrs.allRules()) < len(optimizedRulesToSG) {
		optimizedRulesToIPAddrs = ruleGroups.rulesToIPAddrs.allRules()
	}

	// append both slices together
	return append(optimizedRulesToSG, optimizedRulesToIPAddrs...)
}

func reduceSGRulesToSG(spans *sgSpansToSGPerProtocol, direction ir.Direction) []ir.SGRule {
	spans = compressSpansToSG(spans)

	// convert spans to SG rules
	tcpRules := tcpudpSGSpanToSGRules(spans.tcp, direction, true)
	udpRules := tcpudpSGSpanToSGRules(spans.udp, direction, false)
	icmpRules := icmpSGSpanToSGRules(spans.icmp, direction)
	allRules := protocolAllSGSpanToSGRules(spans.all, direction)

	// return all rules
	return append(tcpRules, append(udpRules, append(icmpRules, allRules...)...)...)
}

func reduceSGRulesToIPAddrs(spans *sgSpansToIPPerProtocol, direction ir.Direction) []ir.SGRule {
	spans = compressToAllProtocolRule(spans)

	// spans to SG rules
	tcpRules := tcpudpIPSpanToSGRules(spans.tcp, spans.all, direction, true)
	udpRules := tcpudpIPSpanToSGRules(spans.udp, spans.all, direction, false)
	icmpRules := icmpSpanToSGRules(spans.icmp, spans.all, direction)
	allRules := allSpanIPToSGRules(spans.all, direction)

	// return all rules
	return append(tcpRules, append(udpRules, append(icmpRules, allRules...)...)...)
}

// divide SGCollection to TCP/UDP/ICMP/ProtocolALL X SGRemote/IPAddrs rules
func divideSGRules(rules []ir.SGRule) *sgRuleGroups {
	rulesToSG := &sgRulesPerProtocol{tcp: make([]ir.SGRule, 0), udp: make([]ir.SGRule, 0),
		icmp: make([]ir.SGRule, 0), all: make([]ir.SGRule, 0)}
	rulesToIPAddrs := &sgRulesPerProtocol{tcp: make([]ir.SGRule, 0), udp: make([]ir.SGRule, 0),
		icmp: make([]ir.SGRule, 0), all: make([]ir.SGRule, 0)}

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
	return &sgRuleGroups{rulesToSG: rulesToSG, rulesToIPAddrs: rulesToIPAddrs}
}

func (s *sgRulesPerProtocol) allRules() []ir.SGRule {
	return append(s.tcp, append(s.udp, append(s.icmp, s.all...)...)...)
}
