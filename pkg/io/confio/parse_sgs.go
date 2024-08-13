/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func ReadSGs(filename string) (*ir.Spec, error) {
	_, err := readModel(filename)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
