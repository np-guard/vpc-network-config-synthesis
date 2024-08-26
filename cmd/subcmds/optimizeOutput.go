/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func writeOptimizeOutput(args *inArgs, collection ir.OptimizeCollection, vpcNames []string) error {
	if err := checkOutputFlags(args); err != nil {
		return err
	}
	_, isACLCollection := collection.(*ir.ACLCollection)
	if err := writeLocals(args, vpcNames, isACLCollection); err != nil {
		return err
	}
	return nil
}

func pickOptimizeWriter(args *inArgs, data *bytes.Buffer) (ir.SynthWriter, error) {
	w := bufio.NewWriter(data)
	switch args.outputFmt {
	case tfOutputFormat:
		return tfio.NewWriter(w), nil
	default:
		return nil, fmt.Errorf("bad output format: %q", args.outputFmt)
	}
}
