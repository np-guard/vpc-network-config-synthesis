/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"github.com/spf13/cobra"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

func NewSynthACLCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "Generate Networks ACLs from connectivity specification",
		Long: `Generate Network ACLs to only allow the specified connectivity, either for each subnet separately or per VPC.
		Endpoints in the required-connectivity specification may be subnets, subnet segments, CIDR segments and externals.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return synthesis(cmd, args, synth.NewACLSynthesizer, args.singleacl, false)
		},
	}

	cmd.Flags().BoolVar(&args.singleacl, singleACLFlag, false, "whether to generate a single acl")

	return cmd
}
