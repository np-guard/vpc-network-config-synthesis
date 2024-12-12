/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"github.com/spf13/cobra"

	acloptimizer "github.com/np-guard/vpc-network-config-synthesis/pkg/optimize/acl"
)

const aclNameFlag = "acl-name"

func newOptimizeACLCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "OptimizeACL attempts to reduce the number of nACL rules in an nACL without changing the semantic.",
		Long:  `OptimizeACL attempts to reduce the number of nACL rules in an nACL without changing the semantic.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return optimization(cmd, args, acloptimizer.NewACLOptimizer, false)
		},
	}

	// flags
	cmd.PersistentFlags().StringVarP(&args.firewallName, aclNameFlag, "n", "", "which nacl to optimize")

	return cmd
}
