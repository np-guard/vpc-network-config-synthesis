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
		Short: "Generate Security Groups from connectivity specification",
		Long: `Generate Security Groups for Network Interfaces and VPEs to only allow the specified connectivity. 
		Endpoints in the required-connectivity specification may be Instances (VSIs), Network Interfaces, VPEs and externals.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SilenceUsage = true
			spec, err := unmarshal(args)
			if err != nil {
				return err
			}
			collection := synth.MakeSG(spec)
			return writeOutput(args, collection, &spec.Defs.ConfigDefs)
		},
	}
	return cmd
}
