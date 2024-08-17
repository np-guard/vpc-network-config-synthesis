/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/confio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/opt"
)

func NewOptimizeSGCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sg",
		Short: "OptimizeSG attempts to reduce the number of security group rules in a SG without changing the semantic.",
		Long:  `OptimizeSG attempts to reduce the number of security group rules in a SG without changing the semantic.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return optimize(cmd, args)
		},
	}

	cmd.Flags().StringVarP(&args.sgName, sgNameFlag, "s", "", "which SG to optimize")
	_ = cmd.MarkFlagRequired(sgNameFlag) // Todo: delete this line. if sgName flag is not supplied - optimize all SGs

	return cmd
}

func optimize(cmd *cobra.Command, args *inArgs) error {
	cmd.SilenceUsage = true // if we got this far, flags are syntactically correct, so no need to print usage
	sgs, err := confio.ReadSGs(args.configFile)
	if err != nil {
		return fmt.Errorf("could not parse config file %v: %w", args.configFile, err)
	}
	_, err = opt.ReduceSGRules(sgs, args.sgName)
	if err != nil {
		return err
	}

	// Output

	return nil

}
