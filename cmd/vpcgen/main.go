/*
VpcGen reads specification files adhering to spec_schema.json and generates network ACLs
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

func readSubnetMap(filename string) (map[string]string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := map[string][]map[string]interface{}{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	subnetMap := make(map[string]string)
	for _, subnet := range config["subnets"] {
		subnetMap[subnet["name"].(string)] = subnet["ipv4_cidr_block"].(string)
	}
	return subnetMap, nil
}

func main() {
	connectivityFilename := os.Args[1]
	spec, err := synth.UnmarshalSpec(connectivityFilename)
	if err != nil {
		log.Fatalf("Could not parse connectivity file %s: %s", connectivityFilename, err)
	}
	configFilename := os.Args[2]
	subnetMap, err := readSubnetMap(configFilename)
	if err != nil {
		log.Fatalf("Could not parse config file %s: %s", configFilename, err)
	}
	fmt.Println(spec.MakeACL(subnetMap))
}
