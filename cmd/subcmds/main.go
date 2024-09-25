/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

func Main(args []string) error {
	rootCmd := NewRootCommand()
	rootCmd.SetArgs(args[1:])
	return rootCmd.Execute()
}
