/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

/*
VpcGen reads specification files adhering to spec_schema.json and generates network ACLs and security groups
*/
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/confio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/csvio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/jsonio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/mdio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

// Output targets
const (
	aclTarget       = "acl"
	sgTarget        = "sg"
	singleaclTarget = "singleacl"
	defaultTarget   = aclTarget
)

// Output formats
const (
	tfOutputFormat      = "tf"
	csvOutputFormat     = "csv"
	mdOutputFormat      = "md"
	apiOutputFormat     = "api"
	defaultOutputFormat = csvOutputFormat
)

// Input formats
const (
	jsonInputFormat = "json"
)

const defaultFilePermission = 0o666

func inferFormatUsingFilename(filename string) string {
	switch {
	case filename == "":
		return defaultOutputFormat
	case strings.HasSuffix(filename, ".tf"):
		return tfOutputFormat
	case strings.HasSuffix(filename, ".csv"):
		return csvOutputFormat
	case strings.HasSuffix(filename, ".md"):
		return mdOutputFormat
	case strings.HasSuffix(filename, ".json"):
		return apiOutputFormat
	default:
		return ""
	}
}

func pickOutputFormat(outputFormat, outputFile string) (string, error) {
	inferredOutputFormat := inferFormatUsingFilename(outputFile)
	if outputFormat != "" {
		if outputFile != "" && inferredOutputFormat != "" && inferredOutputFormat != outputFormat {
			return "", fmt.Errorf("output file %v is expected to use format %v, but -fmt %v is supplied",
				outputFile, inferredOutputFormat, outputFormat)
		}
		return outputFormat, nil
	}
	if inferredOutputFormat == "" {
		return "", fmt.Errorf("unknown format for file %v. Please supply format using -fmt flag, or use a known extension", outputFile)
	}
	return inferredOutputFormat, nil
}

func pickWriter(format string, data *bytes.Buffer, connectivityFilename string) (ir.Writer, error) {
	w := bufio.NewWriter(data)
	switch format {
	case tfOutputFormat:
		return tfio.NewWriter(w), nil
	case csvOutputFormat:
		return csvio.NewWriter(w), nil
	case mdOutputFormat:
		return mdio.NewWriter(w), nil
	case apiOutputFormat:
		return confio.NewWriter(w, connectivityFilename)
	default:
		return nil, fmt.Errorf("bad output format: %q", format)
	}
}

func pickReader(format string) (ir.Reader, error) {
	switch format {
	case jsonInputFormat:
		return jsonio.NewReader(), nil
	default:
		return nil, fmt.Errorf("bad input format: %q", format)
	}
}

func generate(model *ir.Spec, target string) ir.Collection {
	switch target {
	case sgTarget:
		model.ComputeBlockedResources()
		return synth.MakeSG(model, synth.Options{})
	case singleaclTarget:
		model.ComputeBlockedSubnets(true)
		return synth.MakeACL(model, synth.Options{SingleACL: true})
	case aclTarget:
		model.ComputeBlockedSubnets(false)
		return synth.MakeACL(model, synth.Options{SingleACL: false})
	default:
		log.Fatalf("Impossible target: %v", target)
	}
	return nil
}

func main() {
	configFilename := flag.String("config", "",
		"JSON file containing config spec")
	target := flag.String("target", defaultTarget,
		fmt.Sprintf("Target resource to generate. One of %q, %q, %q.", aclTarget, sgTarget, singleaclTarget))
	outputFormat := flag.String("fmt", "",
		fmt.Sprintf("Output format. One of %q, %q, %q; must not contradict output file suffix. Default: %q",
			tfOutputFormat, csvOutputFormat, mdOutputFormat, defaultOutputFormat))
	outputFile := flag.String("o", "",
		"Output to file. If unspecified, use stdout. If specified, also determines output format.")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `VpcGen translates connectivity spec to network ACLs or Security Groups.
Usage:
	%s [flags] SPEC_FILE

SPEC_FILE: JSON file containing connectivity spec

Flags:
`, "vpcgen")
		flag.PrintDefaults()
	}
	flag.Parse()

	connectivityFilename := flag.Arg(0)
	if connectivityFilename == "" || flag.NArg() != 1 {
		flag.Usage()
		os.Exit(0)
	}

	var err error

	*outputFormat, err = pickOutputFormat(*outputFormat, *outputFile)
	if err != nil {
		log.Fatal(err)
	}
	reader, err := pickReader(jsonInputFormat)
	if err != nil {
		log.Fatal(err)
	}

	var defs *ir.ConfigDefs
	if *configFilename != "" {
		defs, err = confio.ReadDefs(*configFilename)
		if err != nil {
			log.Fatalf("could not parse config file %v: %v", *configFilename, err)
		}
	} else if *outputFormat == apiOutputFormat {
		log.Fatal("-config parameter must be supplied when using -fmt=api or exporting JSON")
	}

	model, err := reader.ReadSpec(connectivityFilename, defs)
	if err != nil {
		log.Fatalf("Could not parse connectivity file %s: %s", connectivityFilename, err)
	}

	collection := generate(model, *target)

	var data bytes.Buffer
	writer, err := pickWriter(*outputFormat, &data, *configFilename)
	if err != nil {
		log.Fatal(err)
	}
	if err = collection.Write(writer); err != nil {
		log.Fatal(err)
	}

	if *outputFile == "" {
		fmt.Print(data.String())
	} else {
		err = os.WriteFile(*outputFile, data.Bytes(), defaultFilePermission)
		if err != nil {
			log.Fatal(err)
		}
	}
}
