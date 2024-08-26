/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import "github.com/np-guard/vpc-network-config-synthesis/pkg/ir"

func ReadSGs(filename string) (map[ir.SGName]*ir.SG, error) {
	return map[ir.SGName]*ir.SG{}, nil
}
