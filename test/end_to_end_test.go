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
	dataFolder                  = "data"
	defaultSpecName             = "conn_spec.json"
	defaultExpectedSingleName   = "nacl_single_expected.csv"
	defaultExpectedMultipleName = "nacl_multiple_expected.csv"
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

func TestCIDR(t *testing.T) {
	_, err := makeACLCSV(TestCase{folder: "acl_cidr"}, synth.Options{Single: true})
	if err.Error() != "unsupported endpoint type cidr" {
		t.Errorf("No failure for unsupported type; got %v", err)
	}
}

func TestCSVCompare(t *testing.T) {
	suite := map[string]TestCase{
		"single connection 1": {folder: "acl_single_conn1"},
		"single connection 2": {folder: "acl_single_conn2"},
		"duplication":         {folder: "acl_dup"},
		"acl_testing5":        {folder: "acl_testing5", configName: "config_object.json"},
	}
	for testname, c := range suite {
		testcase := c
		for _, single := range []bool{false, true} {
			expectedName := defaultExpectedMultipleName
			if single {
				expectedName = defaultExpectedSingleName
			}
			t.Run(fmt.Sprintf("%v-%v", testname, single), func(t *testing.T) {
				actualCSVString, err := makeACLCSV(testcase, synth.Options{Single: single})
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

func makeACLCSV(c TestCase, opt synth.Options) (csvString string, err error) {
	reader := jsonio.NewReader()

	var defs *ir.ConfigDefs
	if c.configName != "" {
		defs, err = jsonio.ReadDefs(c.resolve(c.configName))
		if err != nil {
			return
		}
	}
	s, err := reader.ReadSpec(c.at(c.specName, defaultSpecName), defs)
	if err != nil {
		return "", err
	}
	acl := synth.MakeACL(s, opt)

	buf := new(bytes.Buffer)
	writer := csvio.NewWriter(buf)
	err = writer.WriteACL(acl)
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
