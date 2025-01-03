/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"fmt"
	"strings"
)

const (
	tfOutputFormat      = "tf"
	csvOutputFormat     = "csv"
	mdOutputFormat      = "md"
	jsonOutputFormat    = "json"
	defaultOutputFormat = csvOutputFormat
)

var outputFormats = []string{tfOutputFormat, csvOutputFormat, mdOutputFormat, jsonOutputFormat}

func updateOutputFormat(args *inArgs) error {
	var err error
	if args.outputFmt != "" {
		return nil
	}
	args.outputFmt, err = inferFormatUsingFilename(args.outputFile)
	return err
}

func inferFormatUsingFilename(filename string) (string, error) {
	switch {
	case filename == "":
		return defaultOutputFormat, nil
	case strings.HasSuffix(filename, ".tf"):
		return tfOutputFormat, nil
	case strings.HasSuffix(filename, ".csv"):
		return csvOutputFormat, nil
	case strings.HasSuffix(filename, ".md"):
		return mdOutputFormat, nil
	case strings.HasSuffix(filename, ".json"):
		return jsonOutputFormat, nil
	default:
		return "", fmt.Errorf("bad output format")
	}
}
