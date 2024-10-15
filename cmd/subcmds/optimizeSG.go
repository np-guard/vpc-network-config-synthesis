/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import "github.com/spf13/cobra"

func newOptimizeSGCommand(_ *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sg",
		Short: "OptimizeSG is not supported yet",
		Long:  `OptimizeSG is not supported yet`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
	}

	return cmd
}
