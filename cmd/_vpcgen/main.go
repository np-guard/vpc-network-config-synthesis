/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package _main

import (
	"github.com/np-guard/vpc-network-config-synthesis/cmd/subcmds"
)

func Main(args []string) error {
	rootCmd := subcmds.NewRootCommand()
	rootCmd.SetArgs(args[1:])
	return rootCmd.Execute()
}
