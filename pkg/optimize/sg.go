/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"log"

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
		tcpToAddrs  []ir.SGRule
		tcpToSG     []ir.SGRule
		udpToAddrs  []ir.SGRule
		udpToSG     []ir.SGRule
		icmpToAddrs []ir.SGRule
		icmpToSG    []ir.SGRule
		allToAddrs  []ir.SGRule
		allToSG     []ir.SGRule
	}

	sgRemotePortsSpan map[*ir.SGName]*interval.CanonicalSet
	sgRemoteIcmpSpan  map[*ir.SGName]*icmp
	sgRemoteAllSpan   []*ir.SGName
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

// Optimize attempts to reduce thr number of SG rules
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

// divideSGRules divides the rules into groups based on their protocol and remote
// and attempts to reduce each group separately
func (s *SGOptimizer) reduceSGRules(rules []ir.SGRule, direction ir.Direction) []ir.SGRule {
	ruleGroups := divideSGRules(rules)

	// rules with SG as a remote
	tcpToSGSpan := tcpudpRulesToSGToPortsSpan(ruleGroups.tcpToSG)
	udpToSGSpan := tcpudpRulesToSGToPortsSpan(ruleGroups.udpToSG)
	icmpToSGSpan := icmpRulesToSGToIcmpSpan(ruleGroups.icmpToSG)
	protocolAllToSGSpan := allProtocolRulesToSGToSpan(ruleGroups.allToSG)
	rulesToSG := reduceSGRulesToSG(tcpToSGSpan, udpToSGSpan, icmpToSGSpan, protocolAllToSGSpan, direction)

	// rules with IPBlock as a remote
	tcpToIPAddrsSpan := tcpudpRulesToIPAddrsToPortsSpan(ruleGroups.tcpToAddrs)
	udpToIPAddrsSpan := tcpudpRulesToIPAddrsToPortsSpan(ruleGroups.udpToAddrs)
	icmpToAddrsSpan := icmpRulesToIPAddrsToIcmpSpan(ruleGroups.icmpToSG)
	protocolAllToIPAddrsSpan := allProtocolRulesToIPAddrsToSpan(ruleGroups.allToAddrs)
	rulesToIPAddrs := reduceSGRulesToIPAddrs(tcpToIPAddrsSpan, udpToIPAddrsSpan, icmpToAddrsSpan, protocolAllToIPAddrsSpan, direction)

	// append both slices together
	return append(rulesToSG, rulesToIPAddrs...)
}

func reduceSGRulesToSG(tcp, udp sgRemotePortsSpan, icmp sgRemoteIcmpSpan, all sgRemoteAllSpan, direction ir.Direction) []ir.SGRule {
	// delete other protocols if all protocol rule exists
	for _, sgName := range all {
		delete(tcp, sgName)
		delete(udp, sgName)
		delete(icmp, sgName)
	}

	// merge tcp, udp and icmp rules into all protocol rule
	for sgName, tcpPorts := range tcp {
		if udpPorts, ok := udp[sgName]; ok {
			if i, ok := icmp[sgName]; ok {
				if i.all() && allPorts(tcpPorts) && allPorts(udpPorts) { // all tcp&udp ports and all icmp types&codes
					delete(tcp, sgName)
					delete(udp, sgName)
					delete(icmp, sgName)
					all = append(all, sgName)
				}
			}
		}
	}

	// convert to spans to SG rules
	tcpRules := tcpudpToSGSpanToSGRules(tcp, direction, true)
	udpRules := tcpudpToSGSpanToSGRules(tcp, direction, false)
	icmpRules := icmpToSGSpanToSGRules(icmp, direction)
	protocolAll := protocolAllToSGSpanToSGRules(all, direction)

	// merge all rules together
	tcpudp := append(tcpRules, udpRules...)
	icmpAll := append(icmpRules, protocolAll...)
	return append(tcpudp, icmpAll...)
}

func reduceSGRulesToIPAddrs() []ir.SGRule {
	return []ir.SGRule{}
}

// divide SGCollection to TCP/UDP/ICMP/ProtocolALL X SGRemote/IPAddrs rules
func divideSGRules(rules []ir.SGRule) *sgRuleGroups {
	tcpAddrs := make([]ir.SGRule, 0)
	tcpSG := make([]ir.SGRule, 0)
	udpAddrs := make([]ir.SGRule, 0)
	udpSG := make([]ir.SGRule, 0)
	icmpAddrs := make([]ir.SGRule, 0)
	icmpSG := make([]ir.SGRule, 0)
	allAddrs := make([]ir.SGRule, 0)
	allSG := make([]ir.SGRule, 0)

	for _, rule := range rules {
		// TCP rule
		if p, ok := rule.Protocol.(netp.TCPUDP); ok && p.ProtocolString() == "TCP" {
			if _, ok := rule.Remote.(*netset.IPBlock); ok {
				tcpAddrs = append(tcpAddrs, rule)
			} else {
				tcpSG = append(tcpSG, rule)
			}
		}

		// UDP rule
		if p, ok := rule.Protocol.(netp.TCPUDP); ok && p.ProtocolString() == "UDP" {
			if _, ok := rule.Remote.(*netset.IPBlock); ok {
				udpAddrs = append(udpAddrs, rule)
			} else {
				udpSG = append(udpSG, rule)
			}
		}

		// ICMP rule
		if _, ok := rule.Protocol.(netp.ICMP); ok {
			if _, ok := rule.Remote.(*netset.IPBlock); ok {
				icmpAddrs = append(icmpAddrs, rule)
			} else {
				icmpSG = append(icmpSG, rule)
			}
		}

		// all protocol rules
		if _, ok := rule.Protocol.(netp.AnyProtocol); ok {
			if _, ok := rule.Remote.(*netset.IPBlock); ok {
				allAddrs = append(allAddrs, rule)
			} else {
				allSG = append(allSG, rule)
			}
		}
	}
	return utils.Ptr(sgRuleGroups{tcpToAddrs: tcpAddrs, tcpToSG: tcpSG,
		udpToAddrs: udpAddrs, udpToSG: udpSG,
		icmpToAddrs: icmpAddrs, icmpToSG: icmpSG,
		allToAddrs: allAddrs, allToSG: allSG})
}
