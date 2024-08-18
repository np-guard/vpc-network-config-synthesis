/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"fmt"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func writeOptimizeOutput(args *inArgs, sgs map[ir.SGName]*ir.SG) error {
	outputFmt, err := inferFormatUsingFilename(args.outputFile)
	if err != nil {
		return err
	}
	if outputFmt != tfOutputFormat {
		return fmt.Errorf("output format must be %s", tfOutputFormat)
	}
	return writeOptimizeTfOutput()
}

func writeOptimizeTfOutput() error {
	return nil
}
