/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"fmt"
	"log"
	"slices"

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
		sgVPC        string
	}

	ruleGroups struct {
		sgRemoteRules *rulesPerProtocol
		ipRemoteRules *rulesPerProtocol
	}

	rulesPerProtocol struct {
		tcp         []*ir.SGRule
		udp         []*ir.SGRule
		icmp        []*ir.SGRule
		anyProtocol []*ir.SGRule
	}

	// ir.SGName refers to the remote SG
	sgCubesPerProtocol struct {
		tcp         map[ir.SGName]*netset.PortSet
		udp         map[ir.SGName]*netset.PortSet
		icmp        map[ir.SGName]*netset.ICMPSet
		anyProtocol []ir.SGName
	}

	// ipblocks refers to remote IPs
	ipCubesPerProtocol struct {
		tcp         []ds.Pair[*netset.IPBlock, *netset.PortSet]
		udp         []ds.Pair[*netset.IPBlock, *netset.PortSet]
		icmp        []ds.Pair[*netset.IPBlock, *netset.ICMPSet]
		anyProtocol *netset.IPBlock
	}
)

func NewSGOptimizer(collection ir.Collection, sgName string) optimize.Optimizer {
	components := ir.ScopingComponents(sgName)
	if len(components) == 1 {
		return &sgOptimizer{sgCollection: collection.(*ir.SGCollection), sgName: ir.SGName(sgName), sgVPC: ""}
	}
	return &sgOptimizer{sgCollection: collection.(*ir.SGCollection), sgName: ir.SGName(components[1]), sgVPC: components[0]}
}

// Optimize attempts to reduce the number of SG rules
// if -n was supplied, it will attempt to reduce the number of rules only in the requested SG
// otherwise, it will attempt to reduce the number of rules in all SGs
func (s *sgOptimizer) Optimize() (ir.Collection, error) {
	if s.sgName != "" {
		for _, vpcName := range utils.SortedMapKeys(s.sgCollection.SGs) {
			if s.sgVPC != "" && s.sgVPC != vpcName {
				continue
			}
			if _, ok := s.sgCollection.SGs[vpcName][s.sgName]; ok {
				s.optimizeSG(s.sgCollection.SGs[vpcName][s.sgName])
				return s.sgCollection, nil
			}
		}
		return nil, fmt.Errorf("could not find %s sg", s.sgName)
	}

	for _, vpcName := range utils.SortedMapKeys(s.sgCollection.SGs) {
		for _, sgName := range utils.SortedMapKeys(s.sgCollection.SGs[vpcName]) {
			s.optimizeSG(s.sgCollection.SGs[vpcName][sgName])
		}
	}
	return s.sgCollection, nil
}

// optimizeSG attempts to reduce the number of SG rules
// the algorithm attempts to reduce both inbound and outbound rules separately
// A message is printed to the log at the end of the algorithm
func (s *sgOptimizer) optimizeSG(sg *ir.SG) {
	reducedRules := 0

	// reduce inbound rules first
	newInboundRules := s.reduceSGRules(sg.InboundRules, ir.Inbound)
	if len(sg.InboundRules) > len(newInboundRules) {
		reducedRules += len(sg.InboundRules) - len(newInboundRules)
		sg.InboundRules = newInboundRules
	}

	// reduce outbound rules second
	newOutboundRules := s.reduceSGRules(sg.OutboundRules, ir.Outbound)
	if len(sg.OutboundRules) > len(newOutboundRules) {
		reducedRules += len(sg.OutboundRules) - len(newOutboundRules)
		sg.OutboundRules = newOutboundRules
	}

	// print a message to the log
	if reducedRules == 0 {
		log.Printf("no rules were reduced in sg %s\n", string(sg.SGName))
	} else {
		log.Printf("the number of rules in sg %s was reduced by %d\n", string(sg.SGName), reducedRules)
	}
}

// reduceSGRules attempts to reduce the number of rules with different remote types separately
func (s *sgOptimizer) reduceSGRules(rules []*ir.SGRule, direction ir.Direction) []*ir.SGRule {
	// separate all rules to groups of protocol X remote ([tcp, udp, icmp, protocolAll] X [ip, sg])
	ruleGroups := divideSGRules(rules)

	// rules with SG as a remote
	optimizedRulesToSG := reduceRulesSGRemote(rulesToSGCubes(ruleGroups.sgRemoteRules), direction)
	originalRulesToSG := ruleGroups.sgRemoteRules.allRules()
	if len(originalRulesToSG) <= len(optimizedRulesToSG) { // failed to reduce number of rules
		optimizedRulesToSG = originalRulesToSG
	}

	// rules with IPBlock as a remote
	optimizedRulesToIPAddrs := reduceRulesIPRemote(rulesToIPCubes(ruleGroups.ipRemoteRules), direction)
	originalRulesToIPAddrs := ruleGroups.ipRemoteRules.allRules()
	if len(originalRulesToIPAddrs) <= len(optimizedRulesToSG) { // failed to reduce number of rules
		optimizedRulesToIPAddrs = originalRulesToIPAddrs
	}

	return slices.Concat(optimizedRulesToSG, optimizedRulesToIPAddrs)
}

func reduceRulesSGRemote(cubes *sgCubesPerProtocol, direction ir.Direction) []*ir.SGRule {
	reduceCubesWithSGRemote(cubes)

	// cubes to SG rules
	tcpRules := tcpudpSGCubesToRules(cubes.tcp, direction, true)
	udpRules := tcpudpSGCubesToRules(cubes.udp, direction, false)
	icmpRules := icmpSGCubesToRules(cubes.icmp, direction)
	anyProtocolRules := anyProtocolCubesToRules(cubes.anyProtocol, direction)

	// return all rules
	return slices.Concat(tcpRules, udpRules, icmpRules, anyProtocolRules)
}

func reduceRulesIPRemote(cubes *ipCubesPerProtocol, direction ir.Direction) []*ir.SGRule {
	reduceIPCubes(cubes)

	// cubes to SG rules
	tcpRules := tcpudpIPCubesToRules(cubes.tcp, cubes.anyProtocol, direction, true)
	udpRules := tcpudpIPCubesToRules(cubes.udp, cubes.anyProtocol, direction, false)
	icmpRules := icmpIPCubesToRules(cubes.icmp, cubes.anyProtocol, direction)
	anyProtocolRules := anyProtocolIPCubesToRules(cubes.anyProtocol, direction)

	// return all rules
	return slices.Concat(tcpRules, udpRules, icmpRules, anyProtocolRules)
}

// divide SGCollection to TCP/UDP/ICMP/anyProtocols X SGRemote/IPAddrs rules
func divideSGRules(rules []*ir.SGRule) *ruleGroups {
	rulesToSG := &rulesPerProtocol{}
	rulesToIPAddrs := &rulesPerProtocol{}

	for _, rule := range rules {
		switch p := rule.Protocol.(type) {
		case netp.TCPUDP:
			//nolint:nestif // if statements
			if p.ProtocolString() == "TCP" {
				if isRemoteIPBlock(rule) {
					rulesToIPAddrs.tcp = append(rulesToIPAddrs.tcp, rule)
				} else {
					rulesToSG.tcp = append(rulesToSG.tcp, rule)
				}
			} else {
				if isRemoteIPBlock(rule) {
					rulesToIPAddrs.udp = append(rulesToIPAddrs.udp, rule)
				} else {
					rulesToSG.udp = append(rulesToSG.udp, rule)
				}
			}
		case netp.ICMP:
			if isRemoteIPBlock(rule) {
				rulesToIPAddrs.icmp = append(rulesToIPAddrs.icmp, rule)
			} else {
				rulesToSG.icmp = append(rulesToSG.icmp, rule)
			}
		case netp.AnyProtocol:
			if isRemoteIPBlock(rule) {
				rulesToIPAddrs.anyProtocol = append(rulesToIPAddrs.anyProtocol, rule)
			} else {
				rulesToSG.anyProtocol = append(rulesToSG.anyProtocol, rule)
			}
		}
	}
	return &ruleGroups{sgRemoteRules: rulesToSG, ipRemoteRules: rulesToIPAddrs}
}

func isRemoteIPBlock(rule *ir.SGRule) bool {
	_, ok := rule.Remote.(*netset.IPBlock)
	return ok
}

func (s *rulesPerProtocol) allRules() []*ir.SGRule {
	return slices.Concat(s.tcp, s.udp, s.icmp, s.anyProtocol)
}
