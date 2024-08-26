/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import "github.com/spf13/cobra"

func NewOptimizeSGCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sg",
		Short: "OptimizeSG attempts to reduce the number of security group rules in a SG without changing the semantic.",
		Long:  `OptimizeSG attempts to reduce the number of security group rules in a SG without changing the semantic.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
	}

	return cmd
}
