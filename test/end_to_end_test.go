package test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/csvio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/jsonio"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
)

const (
	dataFolder          = "data"
	defaultSpecName     = "conn_spec.json"
	defaultExpectedName = "nacl_single_expected.csv"
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
	_, err := makeACLCSV(TestCase{folder: "cidr"})
	if err.Error() != "unsupported endpoint type cidr" {
		t.Errorf("No failure for unsupported type; got %v", err)
	}
}

func TestCSVCompare(t *testing.T) {
	suite := map[string]TestCase{
		"single connection 1": {folder: "single_conn1"},
		"single connection 2": {folder: "single_conn2"},
		"duplication":         {folder: "dup"},
		"acl_testing5":        {folder: "acl_testing5", configName: "config_object.json"},
	}
	for testname, c := range suite {
		testcase := c
		t.Run(testname, func(t *testing.T) {
			actualCSVString, err := makeACLCSV(testcase)
			if err != nil {
				t.Error(err)
			}
			expectedFile := testcase.at(testcase.expectedName, defaultExpectedName)
			expectedCSVString := readExpectedCSV(expectedFile)
			if expectedCSVString != actualCSVString {
				t.Errorf("%v != %v", expectedCSVString, actualCSVString)
			}
		})
	}
}

func makeACLCSV(c TestCase) (csvString string, err error) {
	reader := jsonio.NewReader()

	var subnets map[string]string
	if c.configName != "" {
		subnets, err = jsonio.ReadSubnetMap(c.resolve(c.configName))
		if err != nil {
			return
		}
	}
	s, err := reader.ReadSpec(c.at(c.specName, defaultSpecName), subnets)
	if err != nil {
		return "", err
	}
	acl := synth.MakeACL(s)

	buf := new(bytes.Buffer)
	writer := csvio.NewWriter(buf)
	err = writer.Write(acl)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func readExpectedCSV(filename string) string {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Bad test: %v", err)
	}
	return string(buf)
}
