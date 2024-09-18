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

	m "github.com/np-guard/vpc-network-config-synthesis/cmd/_vpcgen"
)

func main() {
	err := m.Main(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
