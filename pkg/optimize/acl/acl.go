/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
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
	return a.aclCollection, nil
}
