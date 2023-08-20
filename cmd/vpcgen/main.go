/*
VpcGen reads specification files adhering to spec_schema.json and generates network ACLs
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/csvio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/jsonio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

// Output formats
const (
	tfOutputFormat  = "tf"
	csvOutputFormat = "csv"
)

// Input formats
const (
	jsonInputFormat = "json"
)

func pickWriter(format string) (spec.Writer, error) {
	switch format {
	case tfOutputFormat:
		return tfio.NewWriter(os.Stdout), nil
	case csvOutputFormat:
		return csvio.NewWriter(os.Stdout), nil
	default:
		return nil, fmt.Errorf("bad output format: %q", format)
	}
}

func pickReader(format string) (spec.Reader, error) {
	switch format {
	case jsonInputFormat:
		return jsonio.NewReader(), nil
	default:
		return nil, fmt.Errorf("bad input format: %q", format)
	}
}

func main() {
	connectivityFilename := flag.String("spec", "", "JSON file containing connectivity spec")
	configFilename := flag.String("config", "", "JSON file containing config spec")
	outputFormat := flag.String("fmt", tfOutputFormat, fmt.Sprintf("Output format. One of %q, %q", tfOutputFormat, csvOutputFormat))
	inputFormat := flag.String("inputfmt", jsonInputFormat, fmt.Sprintf("Output format. Must be %q", jsonInputFormat))
	flag.Parse()

	writer, err := pickWriter(*outputFormat)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := pickReader(*inputFormat)
	if err != nil {
		log.Fatal(err)
	}

	var subnets map[string]string
	if *configFilename != "" {
		subnets, err = jsonio.ReadSubnetMap(*configFilename)
		if err != nil {
			log.Fatalf("could not parse config file %v: %v", *configFilename, err)
		}
	}

	s, err := reader.ReadSpec(*connectivityFilename, subnets)
	if err != nil {
		log.Fatalf("Could not parse connectivity file %s: %s", *connectivityFilename, err)
	}

	finalACL := synth.MakeACL(s)

	if err := writer.Write(finalACL); err != nil {
		log.Fatal(err)
	}
}
