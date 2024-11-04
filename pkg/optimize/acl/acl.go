/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type (
	aclOptimizer struct {
		aclCollection *ir.ACLCollection
		aclName       string
		aclVPC        *string
	}

	aclRulesPerProtocol struct {
		tcp  []*ir.ACLRule
		udp  []*ir.ACLRule
		icmp []*ir.ACLRule
		all  []*ir.ACLRule
	}
)

func NewACLOptimizer(collection ir.Collection, aclName string) optimize.Optimizer {
	components := ir.ScopingComponents(aclName)
	if len(components) == 1 {
		return &aclOptimizer{aclCollection: collection.(*ir.ACLCollection), aclName: aclName, aclVPC: nil}
	}
	return &aclOptimizer{aclCollection: collection.(*ir.ACLCollection), aclName: components[1], aclVPC: &components[0]}
}

func (a *aclOptimizer) Optimize() (ir.Collection, error) {
	if a.aclName != "" {
		for _, vpcName := range utils.SortedMapKeys(a.aclCollection.ACLs) {
			if a.aclVPC != nil || a.aclVPC != &vpcName {
				continue
			}
			if _, ok := a.aclCollection.ACLs[vpcName][a.aclName]; ok {
				a.optimizeACL(vpcName, a.aclName)
				return a.aclCollection, nil
			}
		}
		return nil, fmt.Errorf("could no find %s acl", a.aclName)
	}

	for _, vpcName := range utils.SortedMapKeys(a.aclCollection.ACLs) {
		for _, aclName := range utils.SortedMapKeys(a.aclCollection.ACLs[vpcName]) {
			a.optimizeACL(vpcName, aclName)
		}
	}
	return a.aclCollection, nil
}

func (a *aclOptimizer) optimizeACL(vpcName, aclName string) {
	acl := a.aclCollection.ACLs[vpcName][aclName]
	reducedRules := 0

	// reduce inbound rules first
	newInboundRules := a.reduceACLRules(acl.Inbound, ir.Inbound)
	if len(acl.Inbound) > len(newInboundRules) {
		reducedRules += len(acl.Inbound) - len(newInboundRules)
		acl.Inbound = newInboundRules
	}

	// reduce outbound rules second
	newOutboundRules := a.reduceACLRules(acl.Outbound, ir.Outbound)
	if len(acl.Outbound) > len(newOutboundRules) {
		reducedRules += len(acl.Outbound) - len(newOutboundRules)
		acl.Outbound = newOutboundRules
	}

	// print a message to the log
	switch {
	case reducedRules == 0:
		log.Printf("no rules were reduced in sg %s\n", aclName)
	case reducedRules == 1:
		log.Printf("1 rule was reduced in sg %s\n", aclName)
	default:
		log.Printf("%d rules were reduced in sg %s\n", reducedRules, aclName)
	}
}

func (a *aclOptimizer) reduceACLRules(rules []*ir.ACLRule, direction ir.Direction) []*ir.ACLRule {
	_ = divideACLRules(rules)
	return []*ir.ACLRule{}
}

func divideACLRules(rules []*ir.ACLRule) *aclRulesPerProtocol {
	res := &aclRulesPerProtocol{tcp: make([]*ir.ACLRule, 0), udp: make([]*ir.ACLRule, 0),
		icmp: make([]*ir.ACLRule, 0), all: make([]*ir.ACLRule, 0)}
	for _, rule := range rules {
		switch p := rule.Protocol.(type) {
		case netp.TCPUDP:
			if p.ProtocolString() == "TCP" {
				res.tcp = append(res.tcp, rule)
			} else {
				res.udp = append(res.udp, rule)
			}
		case netp.ICMP:
			res.icmp = append(res.icmp, rule)
		case netp.AnyProtocol:
			res.all = append(res.all, rule)
		}
	}
	return res
}
