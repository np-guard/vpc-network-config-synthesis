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
	outputFmtFlag  = "format"
	outputFileFlag = "output-file"
	outputDirFlag  = "output-dir"
	prefixFlag     = "prefix"
	singleACLFlag  = "single"
	localsFlag     = "locals"
)

type inArgs struct {
	configFile   string
	specFile     string
	outputFmt    string
	outputFile   string
	outputDir    string
	prefix       string
	firewallName string
	singleacl    bool
	locals       bool
}

func newRootCommand() *cobra.Command {
	args := &inArgs{}

	// allow PersistentPreRunE
	cobra.EnableTraverseRunHooks = true

	rootCmd := &cobra.Command{
		Use:   "vpcgen",
		Short: "A tool for synthesizing and optimizing VPC network configurations",
		Long:  `A tool for synthesizing and optimizing VPC network configurations, namely Network ACLs and Security Groups.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return validateFlags(args)
		},
	}

	// flags
	rootCmd.PersistentFlags().StringVarP(&args.configFile, configFlag, "c", "",
		"JSON file containing a configuration object of existing resources")
	rootCmd.PersistentFlags().StringVarP(&args.outputFmt, outputFmtFlag, "f", "", "Output format; "+mustBeOneOf(outputFormats))
	rootCmd.PersistentFlags().StringVarP(&args.outputFile, outputFileFlag, "o", "", "Write all generated resources to the specified file.")
	rootCmd.PersistentFlags().BoolVarP(&args.locals, localsFlag, "l", false,
		"whether to generate a locals.tf file (only possible when the output format is tf)")

	// flags set for all commands
	rootCmd.PersistentFlags().SortFlags = false
	_ = rootCmd.MarkPersistentFlagRequired(configFlag)

	// sub cmds
	rootCmd.AddCommand(NewSynthCommand(args))
	rootCmd.AddCommand(NewOptimizeCommand(args))

	// prevent Cobra from creating a default 'completion' command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// disable help command. should use --help flag instead
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	return rootCmd
}

func mustBeOneOf(values []string) string {
	return fmt.Sprintf("must be one of [%s]", strings.Join(values, ", "))
}
