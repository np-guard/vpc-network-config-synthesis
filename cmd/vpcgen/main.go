/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

/*
VpcGen reads specification files adhering to spec_schema.json and generates network ACLs and security groups
*/
package main

import (
	"log"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/cmd/subcmds"
)

func main() {
	rootCmd := subcmds.NewRootCommand()
	rootCmd.SetArgs(os.Args)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
