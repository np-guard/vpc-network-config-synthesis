/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import "github.com/np-guard/vpc-network-config-synthesis/pkg/ir"

type Optimizer interface {
	// read the collection from the config object file
	ParseCollection(filename string) error

	// optimize number of SG/nACL rules
	Optimize() (ir.Collection, error)

	// returns a slice of all vpc names. used to generate locals file
	VpcNames() []string
}
