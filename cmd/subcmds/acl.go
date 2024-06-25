/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"github.com/spf13/cobra"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

func NewACLCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "nACL generation for subnets",
		Long: `generate an nACL for each subnet separately, or to generate a single nACL for all subnets in 
			the same VPC. The input supports subnets, subnet segments, CIDR segments and externals.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			spec, err := unmarshal(args)
			if err != nil {
				return err
			}
			var collection *ir.ACLCollection
			if args.singleacl {
				collection = synth.MakeACL(spec, synth.Options{SingleACL: true})
			} else {
				collection = synth.MakeACL(spec, synth.Options{SingleACL: false})
			}
			err = writeOutput(args, collection, &spec.Defs.ConfigDefs)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&args.singleacl, singleACLFlag, false, "whether to generate a single acl")

	return cmd
}
