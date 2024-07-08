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
	"path/filepath"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/confio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/csvio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/mdio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const defaultFilePermission = 0o644
const defaultDirectoryPermission = 0o755

func writeOutput(args *inArgs, collection ir.Collection, defs *ir.ConfigDefs) error {
	if err := updateOutputFormat(args); err != nil {
		return err
	}
	if args.outputDir != "" && args.outputFmt == apiOutputFormat {
		return fmt.Errorf("-d cannot be used with format json")
	}
	if args.outputDir != "" { // create the directory if needed
		if err := os.MkdirAll(args.outputDir, defaultDirectoryPermission); err != nil {
			return err
		}
	}
	if err := writeLocals(args, collection, defs); err != nil {
		return err
	}

	var data *bytes.Buffer
	var err error
	if args.outputDir == "" {
		if data, err = writeCollection(args, collection, ""); err != nil {
			return err
		}
		return writeToFile(args.outputFile, data)
	}

	// write each file
	for vpc := range defs.VPCs {
		suffix := vpc + "." + args.outputFmt
		args.outputFile = args.outputDir + "/" + suffix
		if args.prefix != "" {
			args.outputFile = args.outputDir + "/" + args.prefix + "_" + suffix
		}
		if data, err = writeCollection(args, collection, vpc); err != nil {
			return err
		}
		if err := writeToFile(args.outputFile, data); err != nil {
			return err
		}
	}
	return nil
}

func writeCollection(args *inArgs, collection ir.Collection, vpc string) (*bytes.Buffer, error) {
	var data bytes.Buffer
	writer, err := pickWriter(args, &data)
	if err != nil {
		return nil, err
	}
	if err := collection.Write(writer, vpc); err != nil {
		return nil, err
	}
	return &data, nil
}

func writeToFile(outputFile string, data *bytes.Buffer) error {
	if outputFile == "" {
		fmt.Println(data.String())
		return nil
	}
	return os.WriteFile(outputFile, data.Bytes(), defaultFilePermission)
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
		return confio.NewWriter(w, args.configFile)
	default:
		return nil, fmt.Errorf("bad output format: %q", args.outputFmt)
	}
}

func writeLocals(args *inArgs, collection ir.Collection, defs *ir.ConfigDefs) error {
	if !args.locals {
		return nil
	}
	if args.outputFmt != tfOutputFormat {
		return fmt.Errorf("--locals flag requires setting the output format to tf")
	}

	_, isACLCollection := collection.(*ir.ACLCollection)
	var data *bytes.Buffer
	var err error
	if data, err = tfio.WriteLocals(defs, isACLCollection); err != nil {
		return err
	}

	outputFile := ""
	suffix := "/locals.tf"
	if args.outputDir != "" {
		outputFile = args.outputDir + suffix
	} else if args.outputFile != "" {
		outputFile = filepath.Dir(args.outputFile) + suffix
	}
	return writeToFile(outputFile, data)
}
