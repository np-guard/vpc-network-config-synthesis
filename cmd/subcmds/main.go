/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"bytes"

	"github.com/spf13/cobra"
)

func Main(args []string) (string, error) {
	rootCmd := newRootCommand()
	rootCmd.SetArgs(args[1:])
	return cmdWrapper(rootCmd)
}

// also returns a warning as string
func cmdWrapper(cmd *cobra.Command) (string, error) {
	var outBuffer bytes.Buffer
	cmd.SetOut(&outBuffer)

	err := cmd.Execute()
	return outBuffer.String(), err
}
