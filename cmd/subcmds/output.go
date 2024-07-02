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

const defaultFilePermission = 0o666
const defaultDirectoryPermission = 0o755

func writeOutput(args *inArgs, collection ir.Collection, defs *ir.ConfigDefs) error {
	if err := updateFormat(args); err != nil {
		return err
	}
	if args.outputDir != "" && args.outputFmt == apiOutputFormat {
		return fmt.Errorf("-d cannot be used with format json")
	}
	if args.outputDir == "" {
		return writeToFile(args, collection, "")
	}
	// create the directory if needed
	err := os.MkdirAll(args.outputDir, defaultDirectoryPermission)
	if err != nil {
		return err
	}

	// write each file
	for vpc := range defs.VPCs {
		suffix := vpc + "." + args.outputFmt
		args.outputFile = args.outputDir + "/" + suffix
		if args.prefix != "" {
			args.outputFile = args.outputDir + "/" + args.prefix + "_" + suffix
		}
		if err := writeToFile(args, collection, vpc); err != nil {
			return err
		}
	}
	return nil
}

func writeToFile(args *inArgs, collection ir.Collection, vpc string) error {
	var data bytes.Buffer
	writer, err := pickWriter(args, &data)
	if err != nil {
		return err
	}
	if err := collection.Write(writer, vpc); err != nil {
		return err
	}

	if args.outputFile == "" {
		fmt.Print(data.String())
		return nil
	}
	return os.WriteFile(args.outputFile, data.Bytes(), defaultFilePermission)
}

func pickWriter(args *inArgs, data *bytes.Buffer) (ir.Writer, error) {
	w := bufio.NewWriter(data)
	switch args.outputFmt {
	case tfOutputFormat:
		return tfio.NewWriter(w), nil
	case csvOutputFormat:
		return csvio.NewWriter(w), nil
	case mdOutputFormat:
		return mdio.NewWriter(w), nil
	case apiOutputFormat:
		return confio.NewWriter(w, args.specFile)
	default:
		return nil, fmt.Errorf("bad output format: %q", args.outputFmt)
	}
}
