/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/confio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/csvio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/mdio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const defaultDirectoryPermission = 0o755

func writeOutput(args *inArgs, collection ir.SynthCollection, vpcNames []ir.ID) error {
	if err := checkOutputFlags(args); err != nil {
		return err
	}
	if args.outputDir != "" { // create the directory if needed
		if err := os.MkdirAll(args.outputDir, defaultDirectoryPermission); err != nil {
			return err
		}
	}
	_, isACLCollection := collection.(*ir.ACLCollection)
	if err := writeLocals(args, vpcNames, isACLCollection); err != nil {
		return err
	}

	var data *bytes.Buffer
	var err error
	if args.outputDir == "" {
		if data, err = writeSynthCollection(args, collection, ""); err != nil {
			return err
		}
		return writeToFile(args.outputFile, data)
	}

	// write each file
	for _, vpc := range vpcNames {
		suffix := vpc + "." + args.outputFmt
		args.outputFile = args.outputDir + "/" + suffix
		if args.prefix != "" {
			args.outputFile = args.outputDir + "/" + args.prefix + "_" + suffix
		}
		if data, err = writeSynthCollection(args, collection, vpc); err != nil {
			return err
		}
		if err := writeToFile(args.outputFile, data); err != nil {
			return err
		}
	}
	return nil
}

func writeSynthCollection(args *inArgs, collection ir.SynthCollection, vpc string) (*bytes.Buffer, error) {
	var data bytes.Buffer
	writer, err := pickSynthWriter(args, &data)
	if err != nil {
		return nil, err
	}
	if err := collection.WriteSynth(writer, vpc); err != nil {
		return nil, err
	}
	return &data, nil
}

func pickSynthWriter(args *inArgs, data *bytes.Buffer) (ir.SynthWriter, error) {
	w := bufio.NewWriter(data)
	switch args.outputFmt {
	case tfOutputFormat:
		return tfio.NewWriter(w), nil
	case csvOutputFormat:
		return csvio.NewWriter(w), nil
	case mdOutputFormat:
		return mdio.NewWriter(w), nil
	case apiOutputFormat:
		return confio.NewWriter(w, args.configFile)
	default:
		return nil, fmt.Errorf("bad output format: %q", args.outputFmt)
	}
}
