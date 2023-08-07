/*
VpcGen reads specification files adhering to spec_schema.json and generates NetworkACLs
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	vpc1 "github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

func main() {
	filename := os.Args[1]
	spec, err := synth.UnmarshalSpec(filename)
	if err != nil {
		log.Fatalf("Could not parse %s: %s", filename, err)
	}
	dataJSON, err := json.MarshalIndent(spec, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(dataJSON))

	var sgResource vpc1.SecurityGroup
	fmt.Printf("%v", sgResource)
}
