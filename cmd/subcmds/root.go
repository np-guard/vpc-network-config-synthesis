/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	configFlag     = "config"
	specFlag       = "spec"
	outputFmtFlag  = "format"
	outputFileFlag = "output-file"
	outputDirFlag  = "output-dir"
	prefixFlag     = "prefix"
	sgNameFlag     = "sg-name"
	singleACLFlag  = "single"
	localsFlag     = "locals"
)

type inArgs struct {
	configFile string
	specFile   string
	outputFmt  string
	outputFile string
	outputDir  string
	prefix     string
	sgName     string
	singleacl  bool
	locals     bool
}

func NewRootCommand() *cobra.Command {
	args := &inArgs{}

	rootCmd := &cobra.Command{
		Use:   "vpcgen",
		Short: "Tool for automatic synthesis of VPC network configurations",
		Long:  `Tool for automatic synthesis of VPC network configurations, namely Network ACLs and Security Groups.`,
	}

	rootCmd.PersistentFlags().StringVarP(&args.configFile, configFlag, "c", "",
		"JSON file containing a configuration object of existing resources")
	rootCmd.PersistentFlags().StringVarP(&args.outputFmt, outputFmtFlag, "f", "", "Output format; "+mustBeOneOf(outputFormats))
	rootCmd.PersistentFlags().StringVarP(&args.outputFile, outputFileFlag, "o", "", "Write all generated resources to the specified file.")
	rootCmd.PersistentFlags().StringVarP(&args.outputDir, outputDirFlag, "d", "",
		"Write generated resources to files in the specified directory, one file per VPC.")
	rootCmd.PersistentFlags().StringVarP(&args.prefix, prefixFlag, "p", "", "The prefix of the files that will be created.")
	rootCmd.PersistentFlags().BoolVarP(&args.locals, localsFlag, "l", false,
		"whether to generate a locals.tf file (only possible when the output format is tf)")
	rootCmd.PersistentFlags().SortFlags = false

	_ = rootCmd.MarkPersistentFlagRequired(configFlag)
	rootCmd.MarkFlagsMutuallyExclusive(outputFileFlag, outputDirFlag)

	rootCmd.AddCommand(NewSynthCommand(args))
	rootCmd.AddCommand(NewOptimizeCommand(args))

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true}) // disable help command. should use --help flag instead

	return rootCmd
}

func mustBeOneOf(values []string) string {
	return fmt.Sprintf("must be one of [%s]", strings.Join(values, ", "))
}
