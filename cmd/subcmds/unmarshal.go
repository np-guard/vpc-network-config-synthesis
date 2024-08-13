/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"fmt"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/confio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/jsonio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func unmarshal(args *inArgs) (*ir.Spec, error) {
	defs, err := confio.ReadDefs(args.configFile)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file %v: %w", args.configFile, err)
	}

	model, err := jsonio.NewReader().ReadSpec(args.specFile, defs)
	if err != nil {
		return nil, fmt.Errorf("could not parse connectivity file %s: %w", args.specFile, err)
	}

	return model, nil
}
