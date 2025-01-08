/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type (
	aclOptimizer struct {
		aclCollection *ir.ACLCollection
		aclName       string
		aclVPC        string
	}

	protocolTripleSet = ds.TripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TransportSet]
	srcDstProduct     = ds.Product[*netset.IPBlock, *netset.IPBlock]

	aclCubesPerProtocol struct {
		tcpAllow protocolTripleSet
		tcpDeny  protocolTripleSet

		udpAllow protocolTripleSet
		udpDeny  protocolTripleSet

		icmpAllow protocolTripleSet
		icmpDeny  protocolTripleSet

		// initialized in reduceCubes func
		anyProtocolAllow srcDstProduct
		anyProtocolDeny  srcDstProduct
	}
)

func NewACLOptimizer(collection ir.Collection, aclName string) optimize.Optimizer {
	components := ir.ScopingComponents(aclName)
	if len(components) == 1 {
		return &aclOptimizer{aclCollection: collection.(*ir.ACLCollection), aclName: aclName, aclVPC: ""}
	}
	return &aclOptimizer{aclCollection: collection.(*ir.ACLCollection), aclName: components[1], aclVPC: components[0]}
}

func (a *aclOptimizer) Optimize() (ir.Collection, error) {
	if a.aclName != "" {
		for _, vpcName := range utils.SortedMapKeys(a.aclCollection.ACLs) {
			if a.aclVPC != "" && a.aclVPC != vpcName {
				continue
			}
			if _, ok := a.aclCollection.ACLs[vpcName][a.aclName]; ok {
				a.optimizeACL(vpcName, a.aclName)
				return a.aclCollection, nil
			}
		}
		return nil, fmt.Errorf("could not find nACL %s", a.aclName)
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
	if reducedRules == 0 {
		log.Printf("no rules were reduced in acl %s\n", a.aclName)
	} else {
		log.Printf("the number of rules in acl %s was reduced by %d\n", a.aclName, reducedRules)
	}
}

func (a *aclOptimizer) reduceACLRules(rules []*ir.ACLRule, direction ir.Direction) []*ir.ACLRule {
	optimizedRules := aclCubesToRules(aclRulesToCubes(rules), direction)
	if len(rules) > len(optimizedRules) {
		return optimizedRules
	}
	return rules
}
