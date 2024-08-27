/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/confio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type SGOptimizer struct {
	sgCollection *ir.SGCollection
	sgName       string
}

func NewSGOptimizer(sgName string) Optimizer {
	return &SGOptimizer{sgCollection: nil, sgName: sgName}
}

func (s *SGOptimizer) ParseCollection(filename string) error {
	c, err := confio.ReadSGs(filename)
	if err != nil {
		return err
	}
	s.sgCollection = c
	return nil
}

func (s *SGOptimizer) Optimize() ir.OptimizeCollection {
	return s.sgCollection
}

func (s *SGOptimizer) VpcNames() []string {
	return utils.MapKeys(s.sgCollection.SGs)
}
