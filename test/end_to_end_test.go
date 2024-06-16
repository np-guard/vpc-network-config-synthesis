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
	"strings"
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
	defaultConfigName           = "config_object.json"
	defaultExpectedSingleFormat = "%v_single_expected.%v"
	defaultExpectedFormat       = "%v_expected.%v"
)

type TestCase struct {
	folder       string
	expectedName string
	outputFormat string
	separate     bool
	blocked      func(s *ir.Spec) []ir.ID
	maker        func(s *ir.Spec, blocked []ir.ID) ir.Collection
}

func (c *TestCase) resolve(name string) string {
	return dataFolder + "/" + c.folder + "/" + name
}

func aclTestCase(folder, outputFormat string, single, separateOutputs bool) TestCase {
	expectedFormat := defaultExpectedFormat
	if single {
		expectedFormat = defaultExpectedSingleFormat
	}
	return TestCase{
		folder:       folder,
		expectedName: fmt.Sprintf(expectedFormat, "nacl", outputFormat),
		outputFormat: outputFormat,
		separate:     separateOutputs,
		blocked: func(s *ir.Spec) []ir.ID {
			return s.ComputeBlockedSubnets(false) // don't print warning
		},
		maker: func(s *ir.Spec, blocked []ir.ID) ir.Collection {
			return synth.MakeACL(s, synth.Options{SingleACL: single}, blocked)
		},
	}
}

func sgTestCase(folder, outputFormat string, separateOutputs bool) TestCase {
	return TestCase{
		folder:       folder,
		expectedName: fmt.Sprintf(defaultExpectedFormat, "sg", outputFormat),
		outputFormat: outputFormat,
		separate:     separateOutputs,
		blocked: func(s *ir.Spec) []ir.ID {
			return s.ComputeBlockedResources(false) // don't print warning
		},
		maker: func(s *ir.Spec, blocked []ir.ID) ir.Collection {
			return synth.MakeSG(s, synth.Options{}, blocked)
		},
	}
}

func TestCSVCompare(t *testing.T) {
	suite := map[string]TestCase{
		"acl_testing5 csv":                   aclTestCase("acl_testing5", "csv", false, false),
		"acl_testing5 tf":                    aclTestCase("acl_testing5", "tf", false, false),
		"acl_testing5 json":                  aclTestCase("acl_testing5", "json", false, false),
		"acl_testing5 single csv":            aclTestCase("acl_testing5", "csv", true, false),
		"acl_testing5 single tf":             aclTestCase("acl_testing5", "tf", true, false),
		"acl_single_conn csv":                aclTestCase("acl_single_conn", "csv", false, false),
		"acl_single_conn tf":                 aclTestCase("acl_single_conn", "tf", false, false),
		"acl_single_conn single csv":         aclTestCase("acl_single_conn", "csv", true, false),
		"acl_single_conn single tf":          aclTestCase("acl_single_conn", "tf", true, false),
		"acl_cidr_segments1 tf":              aclTestCase("acl_cidr_segments1", "tf", false, false),
		"acl_cidr_segments1 single tf":       aclTestCase("acl_cidr_segments1", "tf", true, false),
		"acl_cidr_segments2 tf":              aclTestCase("acl_cidr_segments2", "tf", false, false),
		"acl_cidr_segments2 single tf":       aclTestCase("acl_cidr_segments2", "tf", true, false),
		"acl_tg_multiple tf":                 aclTestCase("acl_tg_multiple", "tf", false, false),
		"acl_tg_multiple single tf":          aclTestCase("acl_tg_multiple", "tf", true, false),
		"acl_tg_multiple_separate tf":        aclTestCase("acl_tg_multiple", "tf", false, true),
		"acl_tg_multiple_separate single tf": aclTestCase("acl_tg_multiple", "tf", true, true),
		"acl_tg_multiple json":               aclTestCase("acl_tg_multiple", "json", false, false),
		"sg_testing3 csv":                    sgTestCase("sg_testing3", "csv", false),
		"sg_testing3 tf":                     sgTestCase("sg_testing3", "tf", false),
		"sg_single_conn csv":                 sgTestCase("sg_single_conn", "csv", false),
		"sg_single_conn tf":                  sgTestCase("sg_single_conn", "tf", false),
		"sg_tg_multiple csv":                 sgTestCase("sg_tg_multiple", "csv", false),
		"sg_tg_multiple tf":                  sgTestCase("sg_tg_multiple", "tf", false),
		"sg_tg_multiple_separate csv":        sgTestCase("sg_tg_multiple", "csv", true),
		"sg_tg_multiple_separate tf":         sgTestCase("sg_tg_multiple", "tf", true),
		"sg_tg_multiple json":                sgTestCase("sg_tg_multiple", "json", false),
		"sg_externals csv":                   sgTestCase("sg_externals", "csv", false),
		"sg_externals tf":                    sgTestCase("sg_externals", "tf", false),
		"sg_externals md":                    sgTestCase("sg_externals", "md", false),
	}
	for testName := range suite {
		testCase := suite[testName]
		t.Run(testName, func(t *testing.T) {
			s, err := readSpec(&testCase)
			if err != nil {
				t.Fatal(err)
				return
			}
			blocked := testCase.blocked(s)
			collection := testCase.maker(s, blocked)
			if testCase.separate {
				writeMultipleFiles(testCase, t, collection, s)
			} else {
				writeSingleFile(testCase, t, collection)
			}
		})
	}
}

func readSpec(c *TestCase) (s *ir.Spec, err error) {
	reader := jsonio.NewReader()
	defs, err := confio.ReadDefs(c.resolve(defaultConfigName))
	if err != nil {
		return
	}
	return reader.ReadSpec(c.resolve(defaultSpecName), defs)
}

func shrinkWhitespace(s string) string {
	return regexp.MustCompile(`[ \t]+`).ReplaceAllString(s, " ")
}

func write(collection ir.Collection, outputFormat, conn, vpc string) (text string, err error) {
	buf := new(bytes.Buffer)
	var writer ir.Writer
	switch outputFormat {
	case "csv":
		writer = csvio.NewWriter(buf)
	case "tf":
		writer = tfio.NewWriter(buf)
	case "md":
		writer = mdio.NewWriter(buf)
	case "json":
		writer, err = confio.NewWriter(buf, conn)
	}
	if err != nil {
		return "", err
	}
	err = collection.Write(writer, vpc)
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

func writeSingleFile(testCase TestCase, t *testing.T, collection ir.Collection) {
	actual, err := write(collection, testCase.outputFormat, fmt.Sprintf("../test/%s", testCase.resolve(defaultConfigName)), "")
	if err != nil {
		t.Fatal(err)
		return
	}
	expectedFile := testCase.resolve(testCase.expectedName)
	expected := readExpectedFile(expectedFile)
	if expected != actual {
		t.Errorf("%v != %v", expected, actual)
	}
}

func writeMultipleFiles(testCase TestCase, t *testing.T, collection ir.Collection, s *ir.Spec) {
	for vpcName := range s.Defs.VPCs {
		actual, err := write(collection, testCase.outputFormat, fmt.Sprintf("../test/%s", testCase.resolve(defaultConfigName)), vpcName)
		if err != nil {
			t.Fatal(err)
			return
		}
		expectedFile := changeExpectedFileName(testCase, vpcName)
		expected := readExpectedFile(expectedFile)
		if expected != actual {
			t.Errorf("%v != %v", expected, actual)
		}
	}
}

func changeExpectedFileName(testCase TestCase, vpc string) string {
	suffix := vpc + "." + testCase.outputFormat
	if strings.Contains(testCase.expectedName, "single") { // single nacl
		suffix = "single_" + suffix
	}
	return testCase.resolve(suffix)
}
