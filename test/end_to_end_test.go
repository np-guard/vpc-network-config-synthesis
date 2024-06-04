/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/confio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/csvio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/jsonio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/mdio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

const (
	dataFolder                  = "data"
	defaultSpecName             = "conn_spec.json"
	defaultExpectedSingleFormat = "%v_single_expected.%v"
	defaultExpectedFormat       = "%v_expected.%v"
)

const configName string = "config_object.json"

type TestCase struct {
	folder       string
	specName     string
	expectedName string
	outputFormat string
	configName   string
	maker        func(s *ir.Spec) ir.Collection
}

func (c *TestCase) resolve(name string) string {
	return dataFolder + "/" + c.folder + "/" + name
}

func (c *TestCase) at(name, otherwise string) string {
	if name == "" {
		name = otherwise
	}
	return c.resolve(name)
}

func aclTestCase(folder, outputFormat string, single bool) TestCase {
	expectedFormat := defaultExpectedFormat
	if single {
		expectedFormat = defaultExpectedSingleFormat
	}
	return TestCase{
		folder:       folder,
		configName:   configName,
		outputFormat: outputFormat,
		expectedName: fmt.Sprintf(expectedFormat, "nacl", outputFormat),
		maker: func(s *ir.Spec) ir.Collection {
			return synth.MakeACL(s, synth.Options{SingleACL: single}, []ir.ID{})
		},
	}
}

func sgTestCase(folder, outputFormat string) TestCase {
	return TestCase{
		folder:       folder,
		configName:   configName,
		outputFormat: outputFormat,
		expectedName: fmt.Sprintf(defaultExpectedFormat, "sg", outputFormat),
		maker: func(s *ir.Spec) ir.Collection {
			return synth.MakeSG(s, synth.Options{}, []ir.ID{})
		},
	}
}

func TestCSVCompare(t *testing.T) {
	suite := map[string]TestCase{
		"acl_testing5 csv":             aclTestCase("acl_testing5", "csv", false),
		"acl_testing5 tf":              aclTestCase("acl_testing5", "tf", false),
		"acl_testing5 single csv":      aclTestCase("acl_testing5", "csv", true),
		"acl_testing5 single tf":       aclTestCase("acl_testing5", "tf", true),
		"acl_single_conn csv":          aclTestCase("acl_single_conn", "csv", false),
		"acl_single_conn tf":           aclTestCase("acl_single_conn", "tf", false),
		"acl_single_conn single csv":   aclTestCase("acl_single_conn", "csv", true),
		"acl_single_conn single tf":    aclTestCase("acl_single_conn", "tf", true),
		"acl_cidr_segments1 tf":        aclTestCase("acl_cidr_segments1", "tf", false),
		"acl_cidr_segments1 single tf": aclTestCase("acl_cidr_segments1", "tf", true),
		"acl_cidr_segments2 tf":        aclTestCase("acl_cidr_segments2", "tf", false),
		"acl_cidr_segments2 single tf": aclTestCase("acl_cidr_segments2", "tf", true),
		"acl_tg_multiple tf":           aclTestCase("acl_tg_multiple", "tf", false),
		"acl_tg_multiple single tf":    aclTestCase("acl_tg_multiple", "tf", true),
		"sg_testing3 csv":              sgTestCase("sg_testing3", "csv"),
		"sg_testing3 tf":               sgTestCase("sg_testing3", "tf"),
		"sg_single_conn csv":           sgTestCase("sg_single_conn", "csv"),
		"sg_single_conn tf":            sgTestCase("sg_single_conn", "tf"),
		"sg_tg_multiple csv":           sgTestCase("sg_tg_multiple", "csv"),
		"sg_tg_multiple tf":            sgTestCase("sg_tg_multiple", "tf"),
		"sg_externals csv":             sgTestCase("sg_externals", "csv"),
		"sg_externals tf":              sgTestCase("sg_externals", "tf"),
		"sg_externals md":              sgTestCase("sg_externals", "md"),
	}
	for testName := range suite {
		testCase := suite[testName]
		t.Run(testName, func(t *testing.T) {
			s, err := readSpec(&testCase)
			if err != nil {
				t.Fatal(err)
				return
			}
			collection := testCase.maker(s)
			actual, err := write(collection, testCase.outputFormat)
			if err != nil {
				t.Fatal(err)
				return
			}
			expectedFile := testCase.at(testCase.expectedName, testCase.expectedName)
			expected := readExpectedFile(expectedFile)
			if expected != actual {
				t.Errorf("%v != %v", expected, actual)
			}
		})
	}
}

func readSpec(c *TestCase) (s *ir.Spec, err error) {
	reader := jsonio.NewReader()

	var defs *ir.ConfigDefs
	if c.configName != "" {
		defs, err = confio.ReadDefs(c.resolve(c.configName))
		if err != nil {
			return
		}
	}

	return reader.ReadSpec(c.at(c.specName, defaultSpecName), defs)
}

func shrinkWhitespace(s string) string {
	return regexp.MustCompile(`[ \t]+`).ReplaceAllString(s, " ")
}

func write(collection ir.Collection, outputFormat string) (text string, err error) {
	buf := new(bytes.Buffer)
	var writer ir.Writer
	switch outputFormat {
	case "csv":
		writer = csvio.NewWriter(buf)
	case "tf":
		writer = tfio.NewWriter(buf)
	case "md":
		writer = mdio.NewWriter(buf)
	}
	err = collection.Write(writer, "") // write the collection to one file
	if err != nil {
		return "", err
	}
	return shrinkWhitespace(buf.String()), nil
}

func readExpectedFile(filename string) string {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Panicf("Bad test: %v", err)
	}
	return shrinkWhitespace(string(buf))
}
