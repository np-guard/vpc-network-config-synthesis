/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import "github.com/spf13/cobra"

func NewOptimizeCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimize",
		Short: "Optimize is not supported yet",
		Long:  `Optimize is not supported yet`,
	}

	cmd.AddCommand(NewOptimizeSGCommand(args))
	// Todo: add OptimizeACL cmd

	return cmd
}
