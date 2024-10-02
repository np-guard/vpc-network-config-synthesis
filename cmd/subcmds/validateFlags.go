/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import "fmt"

func validateFlags(args *inArgs) error {
	if args.outputDir != "" && args.outputFile != "" {
		return fmt.Errorf("specifying both -d and -o is not allowed")
	}
	if err := updateOutputFormat(args); err != nil {
		return err
	}
	if args.outputDir != "" && args.outputFmt == jsonOutputFormat {
		return fmt.Errorf("-d cannot be used with format json")
	}
	if args.locals && args.outputFmt != tfOutputFormat {
		return fmt.Errorf("--locals flag requires setting the output format to tf")
	}
	return nil
}
