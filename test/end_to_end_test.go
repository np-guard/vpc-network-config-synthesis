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

func TestACLCIDR(t *testing.T) {
	_, err := readSpec(TestCase{folder: "acl_cidr"})
	if err == nil || err.Error() != "unsupported endpoint type cidr" {
		t.Errorf("No failure for unsupported type; got %v", err)
	}
}

func aclTestCase(folder, configName string, single bool) TestCase {
	expectedFormat := defaultExpectedMultipleFormat
	if single {
		expectedFormat = defaultExpectedSingleFormat
	}
	return TestCase{
		folder:       folder,
		configName:   configName,
		expectedName: fmt.Sprintf(expectedFormat, "nacl"),
		maker: func(s *ir.Spec) ir.Collection {
			return synth.MakeACL(s, synth.Options{SingleACL: single})
		},
	}
}

func sgTestCase(folder, configName string) TestCase {
	return TestCase{
		folder:       folder,
		configName:   configName,
		expectedName: fmt.Sprintf(defaultExpectedMultipleFormat, "sg"),
		maker: func(s *ir.Spec) ir.Collection {
			return synth.MakeSG(s, synth.Options{})
		},
	}
}

func TestACLCSVCompare(t *testing.T) {
	suite := map[string]TestCase{
		"acl conn1":              aclTestCase("acl_single_conn1", "", false),
		"acl conn1 single":       aclTestCase("acl_single_conn1", "", true),
		"acl conn2":              aclTestCase("acl_single_conn2", "", false),
		"acl conn2 single":       aclTestCase("acl_single_conn2", "", true),
		"acl duplication":        aclTestCase("acl_dup", "", false),
		"acl duplication single": aclTestCase("acl_dup", "", true),
		"acl_testing5":           aclTestCase("acl_testing5", "config_object.json", false),
		"acl_testing5 single":    aclTestCase("acl_testing5", "config_object.json", true),
		"sg single connection 1": sgTestCase("sg_single_conn1", ""),
		"sg_testing2":            sgTestCase("sg_testing2", "config_object.json"),
	}
	for testname, c := range suite {
		testcase := c
		t.Run(testname, func(t *testing.T) {
			s, err := readSpec(c)
			if err != nil {
				t.Error(err)
			}
			collection := c.maker(s)
			if err != nil {
				t.Error(err)
			}
			actualCSVString, err := writeCSV(collection)
			if err != nil {
				t.Error(err)
			}
			expectedFile := testcase.at(testcase.expectedName, c.expectedName)
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
