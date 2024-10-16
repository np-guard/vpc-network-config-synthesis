/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"github.com/spf13/cobra"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

func NewSynthCommand(args *inArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synth",
		Short: "generate a SG/nACL collection",
		Long: `Generate nACLS or Security Groups to only allow the specified connectivity.
		--config and --spec parameters must be supplied.`,
	}

	cmd.PersistentFlags().StringVarP(&args.specFile, specFlag, "s", "", "JSON file containing spec file")
	_ = cmd.MarkPersistentFlagRequired(specFlag)

	cmd.AddCommand(NewSynthACLCommand(args))
	cmd.AddCommand(NewSynthSGCommand(args))

	return cmd
}

func synthesis(cmd *cobra.Command, args *inArgs, newSynthesizer func(*ir.Spec, bool) synth.Synthesizer, single bool) error {
	cmd.SilenceUsage = true // if we got this far, flags are syntactically correct, so no need to print usage
	spec, err := unmarshal(args)
	if err != nil {
		return err
	}
	synthesizer := newSynthesizer(spec, single)
	collection, err := synthesizer.Synth()
	if err != nil {
		return err
	}
	return writeOutput(args, collection, utils.MapKeys(spec.Defs.ConfigDefs.VPCs))
}
