/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// Output formats
const (
	tfOutputFormat      = "tf"
	csvOutputFormat     = "csv"
	mdOutputFormat      = "md"
	apiOutputFormat     = "api"
	defaultOutputFormat = csvOutputFormat
)

const (
	configFlag     = "config"
	specFlag       = "spec"
	outputFmtFlag  = "format"
	outputFileFlag = "output-file"
	outputDirFlag  = "output-dur"
	prefixFlag     = "prefix"
	singleAclFlag  = "single"
)

type inArgs struct {
	configFile string
	specFile   string
	outputFmt  string
	outputFile string
	outputDir  string
	prefix     string
	singleacl  bool
}

func NewRootCommand() *cobra.Command {
	args := &inArgs{}

	rootCmd := &cobra.Command{
		Use:   "vpc-synthesis",
		Short: "Tool for automatic synthesis of VPC network configurations",
		Long:  `Tool for automatic synthesis of VPC network configurations, namely Network ACLs and Security Groups.`,
	}

	rootCmd.PersistentFlags().StringVarP(&args.configFile, configFlag, "c", "", "JSON file containing config spec")
	rootCmd.PersistentFlags().StringVarP(&args.specFile, specFlag, "s", "", "JSON file containing spec file")
	rootCmd.PersistentFlags().StringVarP(&args.outputFmt, outputFmtFlag, "fmt", "", "Output format; "+mustBeOneOf(outputFormats))
	rootCmd.PersistentFlags().StringVarP(&args.outputFile, outputFileFlag, "o", "", "Write all generated resources to the specified file.")
	rootCmd.PersistentFlags().StringVarP(&args.outputDir, outputDirFlag, "d", "",
		"Write generated resources to files in the specified directory, one file per VPC.")
	rootCmd.PersistentFlags().StringVar(&args.prefix, prefixFlag, "", "The prefix of the files that will be created.")
	rootCmd.PersistentFlags().SortFlags = false

	if err := rootCmd.MarkFlagRequired(configFlag); err != nil {
		log.Fatalf("%w", err)
	}
	rootCmd.MarkFlagsMutuallyExclusive(outputFileFlag, outputDirFlag)
	rootCmd.MarkFlagsMutuallyExclusive(outputFileFlag, outputFmtFlag)
	// make sure fmt flag is set when dir flag is set

	rootCmd.AddCommand(NewACLCommand(args))
	rootCmd.AddCommand(NewSGCommand(args))

	return rootCmd
}

func mustBeOneOf(values []string) string {
	return fmt.Sprintf("must be one of [%s]", strings.Join(values, ", "))
}
