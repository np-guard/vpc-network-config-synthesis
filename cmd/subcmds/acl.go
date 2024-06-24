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
		Short: "...",
		Long:  `...`,
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

	cmd.Flags().BoolVar(&args.singleacl, singleAclFlag, false, "whether to generate a single acl")

	return cmd
}
