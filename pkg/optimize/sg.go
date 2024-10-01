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
}

func NewSGOptimizer(collection ir.Collection, sgName string) Optimizer {
	return &SGOptimizer{sgCollection: collection.(*ir.SGCollection), sgName: ir.SGName(sgName)}
}

func (s *SGOptimizer) Optimize() (ir.Collection, error) {
	return s.sgCollection, nil
}
