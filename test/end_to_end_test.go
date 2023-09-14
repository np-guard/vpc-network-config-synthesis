package test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/csvio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/jsonio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

const (
	dataFolder                    = "data"
	defaultSpecName               = "conn_spec.json"
	defaultExpectedSingleFormat   = "%v_single_expected.csv"
	defaultExpectedMultipleFormat = "%v_expected.csv"
)

type TestCase struct {
	folder       string
	specName     string
	expectedName string
	configName   string
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

func TestACLCIDR(t *testing.T) {
	_, err := readSpec(TestCase{folder: "acl_cidr"})
	if err == nil || err.Error() != "unsupported endpoint type cidr" {
		t.Errorf("No failure for unsupported type; got %v", err)
	}
}

func TestACLCSVCompare(t *testing.T) {
	suite := map[string]TestCase{
		"acl single connection 1": {folder: "acl_single_conn1"},
		"acl single connection 2": {folder: "acl_single_conn2"},
		"acl duplication":         {folder: "acl_dup"},
		"acl_testing5":            {folder: "acl_testing5", configName: "config_object.json"},
	}
	for testname, c := range suite {
		testcase := c
		for _, single := range []bool{false, true} {
			expectedFormat := defaultExpectedMultipleFormat
			if single {
				expectedFormat = defaultExpectedSingleFormat
			}
			expectedName := fmt.Sprintf(expectedFormat, "nacl")
			t.Run(fmt.Sprintf("%v-%v", testname, single), func(t *testing.T) {
				s, err := readSpec(c)
				if err != nil {
					t.Error(err)
				}
				acl := synth.MakeACL(s, synth.Options{SingleACL: single})
				if err != nil {
					t.Error(err)
				}
				actualCSVString, err := writeCSV(acl)
				if err != nil {
					t.Error(err)
				}
				expectedFile := testcase.at(testcase.expectedName, expectedName)
				expectedCSVString := readExpectedCSV(expectedFile)
				if expectedCSVString != actualCSVString {
					t.Errorf("%v != %v", expectedCSVString, actualCSVString)
				}
			})
		}
	}
}

func TestSGCSVCompare(t *testing.T) {
	suite := map[string]TestCase{
		"sg single connection 1": {folder: "sg_single_conn1"},
		"sg_testing5":            {folder: "sg_testing2", configName: "config_object.json"},
	}
	for testname, c := range suite {
		testcase := c
		expectedName := fmt.Sprintf(defaultExpectedMultipleFormat, "sg")
		t.Run(testname, func(t *testing.T) {
			s, err := readSpec(c)
			if err != nil {
				t.Error(err)
			}
			sg := synth.MakeSG(s, synth.Options{})
			if err != nil {
				t.Error(err)
			}
			actualCSVString, err := writeCSV(sg)
			if err != nil {
				t.Error(err)
			}
			expectedFile := testcase.at(testcase.expectedName, expectedName)
			expectedCSVString := readExpectedCSV(expectedFile)
			if expectedCSVString != actualCSVString {
				t.Errorf("%v != %v", expectedCSVString, actualCSVString)
			}
		})
	}
}

func readSpec(c TestCase) (s *ir.Spec, err error) {
	reader := jsonio.NewReader()

	var defs *ir.ConfigDefs
	if c.configName != "" {
		defs, err = jsonio.ReadDefs(c.resolve(c.configName))
		if err != nil {
			return
		}
	}

	return reader.ReadSpec(c.at(c.specName, defaultSpecName), defs)
}

func writeCSV(collection ir.Collection) (csvString string, err error) {
	buf := new(bytes.Buffer)
	writer := csvio.NewWriter(buf)
	err = collection.Write(writer)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func readExpectedCSV(filename string) string {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Panicf("Bad test: %v", err)
	}
	return string(buf)
}
