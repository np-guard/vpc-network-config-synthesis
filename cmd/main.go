/*
VpcGen reads specification files adhering to spec_schema.json and generates NetworkACLs
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
	"log"
	"os"
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
}
