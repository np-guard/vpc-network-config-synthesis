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
	_, err := subcmds.Main(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
