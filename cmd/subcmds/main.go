/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"bytes"
)

// also returns a warning as string
func Main(args []string) (string, error) {
	var outBuffer bytes.Buffer
	rootCmd := newRootCommand()
	rootCmd.SetArgs(args[1:])
	rootCmd.SetOut(&outBuffer)
	err := rootCmd.Execute()
	return outBuffer.String(), err
}
