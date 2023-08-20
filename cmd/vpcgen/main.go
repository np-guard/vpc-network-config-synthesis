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

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl/aclcsv"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl/acltf"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

// Output formats
const (
	tfFormat  = "tf"
	csvFormat = "csv"
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

func pickWriter(format string) (acl.Writer, error) {
	switch format {
	case tfFormat:
		return acltf.NewWriter(os.Stdout), nil
	case csvFormat:
		return aclcsv.NewWriter(os.Stdout), nil
	default:
		return nil, fmt.Errorf("bad format: %q", format)
	}
}

func setSubnets(configFilename string, s *spec.Spec) error {
	if configFilename != "" {
		subnetMap, err := readSubnetMap(configFilename)
		if err != nil {
			return fmt.Errorf("could not parse config file %v: %w", configFilename, err)
		}
		err = s.SetSubnets(subnetMap)
		if err != nil {
			return fmt.Errorf("bad subnets: %w", err)
		}
	}
	return nil
}

func main() {
	connectivityFilename := flag.String("spec", "", "JSON file containing connectivity spec")
	configFilename := flag.String("config", "", "JSON file containing config spec")
	outputFormat := flag.String("fmt", tfFormat, fmt.Sprintf("Output format. One of %q, %q", tfFormat, csvFormat))
	flag.Parse()

	writer, err := pickWriter(*outputFormat)
	if err != nil {
		log.Fatal(err)
	}

	s, err := spec.Unmarshal(*connectivityFilename)
	if err != nil {
		log.Fatalf("Could not parse connectivity file %s: %s", *connectivityFilename, err)
	}

	if err := setSubnets(*configFilename, s); err != nil {
		log.Fatal(err)
	}

	finalACL := synth.MakeACL(s)

	if err := writer.Write(finalACL); err != nil {
		log.Fatal(err)
	}
}
