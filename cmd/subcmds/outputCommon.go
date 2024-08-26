/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func checkOutputFlags(args *inArgs) error {
	if err := updateOutputFormat(args); err != nil {
		return err
	}
	if args.outputDir != "" && args.outputFmt == apiOutputFormat {
		return fmt.Errorf("-d cannot be used with format json")
	}
	return nil
}

func writeToFile(outputFile string, data *bytes.Buffer) error {
	if outputFile == "" {
		fmt.Println(data.String())
		return nil
	}
	return os.WriteFile(outputFile, data.Bytes(), defaultFilePermission)
}

func writeLocals(args *inArgs, vpcNames []ir.ID, isACL bool) error {
	if !args.locals {
		return nil
	}
	if args.outputFmt != tfOutputFormat {
		return fmt.Errorf("--locals flag requires setting the output format to tf")
	}

	var data *bytes.Buffer
	var err error
	if data, err = tfio.WriteLocals(vpcNames, isACL); err != nil {
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
