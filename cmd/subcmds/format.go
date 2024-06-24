/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import "strings"

var outputFormats = []string{"tf", "csv", "md", "json"}

func updateFormat(args *inArgs) {
	if args.outputFmt == "" {
		args.outputFmt = inferFormatUsingFilename(args.outputFile)
	}
}

func inferFormatUsingFilename(filename string) string {
	switch {
	case filename == "":
		return defaultOutputFormat
	case strings.HasSuffix(filename, ".tf"):
		return tfOutputFormat
	case strings.HasSuffix(filename, ".csv"):
		return csvOutputFormat
	case strings.HasSuffix(filename, ".md"):
		return mdOutputFormat
	case strings.HasSuffix(filename, ".json"):
		return apiOutputFormat
	default:
		return ""
	}
}
