/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"github.com/spf13/cobra"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

func NewSGCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sg",
		Short: "SG generation for nifs and vpes",
		Long:  `The input supports Instances (VSIs), NIFs, VPEs and externals.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			spec, err := unmarshal(args)
			if err != nil {
				return err
			}
			collection := synth.MakeSG(spec)
			err = writeOutput(args, collection, &spec.Defs.ConfigDefs)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
