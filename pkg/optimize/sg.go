/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

type SGOptimizer struct {
	sgCollection *ir.SGCollection
	sgName       ir.SGName
	sgVPC        *string
}

func NewSGOptimizer(collection ir.Collection, sgName string) Optimizer {
	components := ir.ScopingComponents(sgName)
	if len(components) == 1 {
		return &SGOptimizer{sgCollection: collection.(*ir.SGCollection), sgName: ir.SGName(sgName), sgVPC: nil}
	}
	return &SGOptimizer{sgCollection: collection.(*ir.SGCollection), sgName: ir.SGName(components[1]), sgVPC: &components[0]}
}

func (s *SGOptimizer) Optimize() (ir.Collection, error) {
	return s.sgCollection, nil
}
