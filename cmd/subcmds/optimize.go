/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import "github.com/spf13/cobra"

func NewOptimizeCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimize",
		Short: "optimization of existing SG (nACLS are not supported yet)",
		Long:  `optimization of existing SG (nACLS are not supported yet)`,
	}

	cmd.AddCommand(NewOptimizeSGCommand(args))
	// Todo: add OptimizeACL cmd

	return cmd
}
