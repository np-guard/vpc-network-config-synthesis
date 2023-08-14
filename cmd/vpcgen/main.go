/*
VpcGen reads specification files adhering to spec_schema.json and generates network ACLs
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"

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
	connectivityFilename := flag.String("spec", "", "JSON file containing connectivity spec")
	configFilename := flag.String("config", "", "JSON file containing config spec")
	flag.Parse()

	s, err := spec.Unmarshal(*connectivityFilename)
	if err != nil {
		log.Fatalf("Could not parse connectivity file %s: %s", *connectivityFilename, err)
	}
	if *configFilename != "" {
		subnetMap, err := readSubnetMap(*configFilename)
		if err != nil {
			log.Fatalf("Could not parse config file %s: %s", *configFilename, err)
		}
		err = s.SetSubnets(subnetMap)
		if err != nil {
			log.Fatalf("Bad subnets: %v", err)
		}
	}
	fmt.Println(synth.MakeACL(s))
}
