/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"fmt"

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

}
