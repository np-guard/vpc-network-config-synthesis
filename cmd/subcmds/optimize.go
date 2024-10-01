/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func NewOptimizeCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimize",
		Short: "optimization of existing SG (nACLS are not supported yet)",
		Long:  `optimization of existing SG (nACLS are not supported yet)`,
	}

	// sub cmds
	cmd.AddCommand(NewOptimizeSGCommand(args))

	return cmd
}

func optimization(cmd *cobra.Command, args *inArgs, newOptimizer func(ir.Collection, string) optimize.Optimizer, isSG bool) error {
	cmd.SilenceUsage = true // if we got this far, flags are syntactically correct, so no need to print usage
	collection, err := parseCollection(args, isSG)
	if err != nil {
		return fmt.Errorf("could not parse config file %v: %w", args.configFile, err)
	}
	optimizer := newOptimizer(collection, args.firewallName)
	optimizedCollection, err := optimizer.Optimize()
	if err != nil {
		return err
	}
	return writeOutput(args, optimizedCollection, collection.VpcNames(), false)
}
