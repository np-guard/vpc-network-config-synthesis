/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"log"
	"sort"

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
		tcpToAddrs  []ir.SGRule
		tcpToSG     []ir.SGRule
		udpToAddrs  []ir.SGRule
		udpToSG     []ir.SGRule
		icmpToAddrs []ir.SGRule
		icmpToSG    []ir.SGRule
		allToAddrs  []ir.SGRule
		allToSG     []ir.SGRule
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
		newInboundRules := s.reduceSGRules(sg.InboundRules, ir.Inbound)
		if len(sg.InboundRules) > len(newInboundRules) {
			reducedRules += len(sg.InboundRules) - len(newInboundRules)
			s.sgCollection.SGs[vpcName][s.sgName].InboundRules = newInboundRules
		}
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
	tcpToAddrs := s.reduceSGRulesTcpudpToAddrs(rulesToIPAddrsToPortsSpan(ruleGroups.tcpToAddrs), direction, true)
	tcpToSg := s.reduceSGRulesTcpudpToSG(rulesToSGToPortsSpan(rules), direction, true)
	return append(tcpToAddrs, tcpToSg...)
}

// reduceSGRulesTcpudpToAddrs attempts to reduce the number of rules of tcp/udp rules with ipAddrs as remote
func (s *SGOptimizer) reduceSGRulesTcpudpToAddrs(span []ds.Pair[*netset.IPBlock, *interval.CanonicalSet],
	direction ir.Direction, isTCP bool) []ir.SGRule {
	result := make([]ir.SGRule, len(span))
	for i := range span {
		for _, dstPorts := range span[i].Right.Intervals() {
			p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(dstPorts.Start()), int(dstPorts.End()))
			rule := ir.SGRule{
				Direction: direction,
				Remote:    span[i].Left,
				Protocol:  p,
				Local:     netset.GetCidrAll(),
			}
			result[i] = rule
		}
	}
	return result
}

// reduceSGRulesTcpudpToAddrs attempts to reduce the number of rules of tcp/udp rules with a sg as remote
func (s *SGOptimizer) reduceSGRulesTcpudpToSG(span map[ir.SGName]*interval.CanonicalSet, direction ir.Direction, isTCP bool) []ir.SGRule {
	result := make([]ir.SGRule, 0)
	for sgName, intervals := range span {
		for _, dstPorts := range intervals.Intervals() {
			p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(dstPorts.Start()), int(dstPorts.End()))
			rule := ir.SGRule{
				Direction: direction,
				Remote:    sgName,
				Protocol:  p,
				Local:     netset.GetCidrAll(),
			}
			result = append(result, rule)
		}
	}
	return result
}

// converts []ir.SGRule (where all rules or either TCP/UDP but not both) to a span of (IPBlock X ports)
func rulesToIPAddrsToPortsSpan(rules []ir.SGRule) (p []ds.Pair[*netset.IPBlock, *interval.CanonicalSet]) {
	span := ds.NewProductLeft[*netset.IPBlock, *interval.CanonicalSet]()
	for i := range rules {
		p := rules[i].Protocol.(netp.TCPUDP) // already checked
		r := ds.CartesianPairLeft(rules[i].Remote.(*netset.IPBlock), p.DstPorts().ToSet())
		span = span.Union(r).(*ds.ProductLeft[*netset.IPBlock, *interval.CanonicalSet])
	}
	return sortPartitionsByIPAddrs(span.Partitions())
}

// converts []ir.SGRule (where all rules or either TCP/UDP but not both) to a span intervals for each remote
func rulesToSGToPortsSpan(rules []ir.SGRule) map[ir.SGName]*interval.CanonicalSet {
	result := make(map[ir.SGName]*interval.CanonicalSet)
	for i := range rules {
		p := rules[i].Protocol.(netp.TCPUDP)  // already checked
		remote := rules[i].Remote.(ir.SGName) // already checked
		if result[remote] == nil {
			result[remote] = interval.NewCanonicalSet()
		}
		result[remote].AddInterval(p.DstPorts())
	}
	return result
}

// each IPBlock is a single CIDR/IP address. The IPBlocks are disjoint.
func sortPartitionsByIPAddrs(p []ds.Pair[*netset.IPBlock, *interval.CanonicalSet]) []ds.Pair[*netset.IPBlock, *interval.CanonicalSet] {
	cmp := func(i, j int) bool { return p[i].Left.FirstIPAddress() < p[j].Left.FirstIPAddress() }
	sort.Slice(p, cmp)
	return p
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
