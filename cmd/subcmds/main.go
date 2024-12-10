/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package subcmds

import (
	"bytes"
	"log"
)

// also returns a warning as string
func Main(args []string) (string, error) {
	var outBuffer bytes.Buffer
	rootCmd := newRootCommand()
	rootCmd.SetArgs(args[1:])
	rootCmd.SetOut(&outBuffer)
	err := rootCmd.Execute()
	res := outBuffer.String()
	log.Println(res) // print usage or blocked resources warning
	return res, err
}
