/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import "github.com/np-guard/vpc-network-config-synthesis/pkg/ir"

type Optimizer interface {
	// attempts to reduce number of SG/nACL rules
	Optimize() (ir.Collection, error)
}
