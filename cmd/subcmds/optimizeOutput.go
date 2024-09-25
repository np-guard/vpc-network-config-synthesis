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

func writeOptimizeOutput(args *inArgs, collection ir.Collection, vpcNames []string) error {
	if err := checkOutputFlags(args); err != nil {
		return err
	}
	_, isACLCollection := collection.(*ir.ACLCollection)
	if err := writeLocals(args, vpcNames, isACLCollection); err != nil {
		return err
	}
	data, err := writeOptimizeCollection(args, collection)
	if err != nil {
		return err
	}
	return writeToFile(args.outputFile, data)
}

func writeOptimizeCollection(args *inArgs, collection ir.Collection) (*bytes.Buffer, error) {
	var data bytes.Buffer
	writer, err := pickOptimizeWriter(args, &data)
	if err != nil {
		return nil, err
	}
	if err := collection.Write(writer, ""); err != nil {
		return nil, err
	}
	return &data, nil
}

func pickOptimizeWriter(args *inArgs, data *bytes.Buffer) (ir.Writer, error) {
	w := bufio.NewWriter(data)
	switch args.outputFmt {
	case tfOutputFormat:
		return tfio.NewWriter(w), nil
	default:
		return nil, fmt.Errorf("bad output format: %q", args.outputFmt)
	}
}
