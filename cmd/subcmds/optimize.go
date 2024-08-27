/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func NewOptimizeCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimize",
		Short: "optimization of existing SG (nACLS are not supported yet)",
		Long:  `optimization of existing SG (nACLS are not supported yet)`,
	}

	// flags
	cmd.Flags().StringVarP(&args.firewallName, firewallNameFlag, "s", "", "which vpcFirewall to optimize")

	// flags settings
	_ = cmd.MarkPersistentFlagRequired(firewallNameFlag) // temporary

	// sub cmds
	cmd.AddCommand(NewOptimizeSGCommand(args))

	return cmd
}

func optimization(cmd *cobra.Command, args *inArgs, newOptimizer optimize.Optimizer) error {
	cmd.SilenceUsage = true // if we got this far, flags are syntactically correct, so no need to print usage
	if err := newOptimizer.ParseCollection(args.configFile); err != nil {
		return fmt.Errorf("could not parse config file %v: %w", args.configFile, err)
	}
	return writeOptimizeOutput(args, newOptimizer.Optimize(), newOptimizer.VpcNames())
}
