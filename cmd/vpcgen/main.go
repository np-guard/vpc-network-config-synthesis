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

func pickWriter(format string, data *bytes.Buffer) (ir.Writer, error) {
	w := bufio.NewWriter(data)
	switch format {
	case tfOutputFormat:
		return tfio.NewWriter(w), nil
	case csvOutputFormat:
		return csvio.NewWriter(w), nil
	case mdOutputFormat:
		return mdio.NewWriter(w), nil
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

func writeOutput(collection ir.Collection, defs *ir.ConfigDefs, outputDirectory, outputFormat, outputFile, prefixOfFileNames *string) {
	if *outputDirectory == "" {
		writeToFile(collection, "", outputFormat, outputFile)
	} else {
		// create a directory
		if err := os.Mkdir(*outputDirectory, defaultFilePermission); err != nil {
			log.Fatal(err)
		}

		// write each file
		for vpc := range defs.VPCs {
			suffix := vpc + "." + *outputFormat
			outputPath := *outputDirectory + "/" + suffix
			if *prefixOfFileNames != "" {
				outputPath = *outputDirectory + "/" + *prefixOfFileNames + "_" + suffix
			}
			writeToFile(collection, vpc, outputFormat, &outputPath)
		}
	}
}

func writeToFile(collection ir.Collection, vpc string, outputFormat, outputFile *string) {
	var data bytes.Buffer
	writer, err := pickWriter(*outputFormat, &data)
	if err != nil {
		log.Fatal(err)
	}
	if err = collection.Write(writer, vpc); err != nil {
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

func main() {
	configFilename := flag.String("config", "",
		"JSON file containing config spec")
	target := flag.String("target", defaultTarget,
		fmt.Sprintf("Target resource to generate. One of %q, %q, %q.", aclTarget, sgTarget, singleaclTarget))
	outputFormat := flag.String("fmt", "",
		fmt.Sprintf("Output format. One of %q, %q, %q; must not contradict output file suffix. (default %q)",
			tfOutputFormat, csvOutputFormat, mdOutputFormat, defaultOutputFormat))
	outputFile := flag.String("output-file", "",
		"Output to file. If specified, also determines output format.")
	outputDirectory := flag.String("output-dir", "",
		"Output Directory. If unspecified, output will be written to one file.")
	prefixOfFileNames := flag.String("prefix", "", "The prefix of the files that will be created.")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `VpcGen translates connectivity spec to network ACLs or Security Groups.
Usage:
	%s [flags] SPEC_FILE

SPEC_FILE: JSON file containing connectivity spec, and segments.

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

	if *outputDirectory != "" && *outputFile != "" {
		log.Fatal(fmt.Errorf("could not determine whether to create a folder or not"))
	}

	*outputFormat, err = pickOutputFormat(*outputFormat, *outputFile)
	if err != nil {
		log.Fatal(err)
	} else if *outputFormat == "" {
		log.Fatal("unknown format. Please supply format using -fmt flag, or use a known extension")
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

	writeOutput(collection, defs, outputDirectory, outputFormat, outputFile, prefixOfFileNames)
}
