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

func TestACLCIDR(t *testing.T) {
	_, err := readSpec(&TestCase{folder: "acl_cidr"})
	if err == nil || err.Error() != "unsupported resource type cidr" {
		t.Errorf("No failure for unsupported type; got %v", err)
	}
}

func aclTestCase(folder, configName, outputFormat string, single bool) TestCase {
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
			return synth.MakeACL(s, synth.Options{SingleACL: single})
		},
	}
}

func sgTestCase(folder, configName, outputFormat string) TestCase {
	return TestCase{
		folder:       folder,
		configName:   configName,
		outputFormat: outputFormat,
		expectedName: fmt.Sprintf(defaultExpectedFormat, "sg", outputFormat),
		maker: func(s *ir.Spec) ir.Collection {
			return synth.MakeSG(s, synth.Options{})
		},
	}
}

func TestCSVCompare(t *testing.T) {
	suite := map[string]TestCase{
		"acl duplication csv":          aclTestCase("acl_dup", "", "csv", false),
		"acl duplication single csv":   aclTestCase("acl_dup", "", "csv", true),
		"acl conn1 csv":                aclTestCase("acl_single_conn1", "", "csv", false),
		"acl conn1 single csv":         aclTestCase("acl_single_conn1", "", "csv", true),
		"acl conn1 tf":                 aclTestCase("acl_single_conn1", "", "tf", false),
		"acl conn1 single tf":          aclTestCase("acl_single_conn1", "", "tf", true),
		"acl conn2 csv":                aclTestCase("acl_single_conn2", "", "csv", false),
		"acl conn2 single csv":         aclTestCase("acl_single_conn2", "", "csv", true),
		"acl conn2 tf":                 aclTestCase("acl_single_conn2", "", "tf", false),
		"acl conn2 single tf":          aclTestCase("acl_single_conn2", "", "tf", true),
		"acl_testing5 csv":             aclTestCase("acl_testing5", "config_object.json", "csv", false),
		"acl_testing5 tf":              aclTestCase("acl_testing5", "config_object.json", "tf", false),
		"acl_testing5 single csv":      aclTestCase("acl_testing5", "config_object.json", "csv", true),
		"acl_testing5 single tf":       aclTestCase("acl_testing5", "config_object.json", "tf", true),
		"acl_cidr_segments csv":        aclTestCase("acl_cidr_segments", "config_object.json", "csv", false),
		"acl_cidr_segments tf":         aclTestCase("acl_cidr_segments", "config_object.json", "tf", false),
		"acl_cidr_segments single csv": aclTestCase("acl_cidr_segments", "config_object.json", "csv", true),
		"acl_cidr_segments single tf":  aclTestCase("acl_cidr_segments", "config_object.json", "tf", true),
		"sg single connection 1 csv":   sgTestCase("sg_single_conn1", "", "csv"),
		"sg single connection 1 tf":    sgTestCase("sg_single_conn1", "", "tf"),
		"sg_testing2 csv":              sgTestCase("sg_testing2", "config_object.json", "csv"),
		"sg_testing2 tf":               sgTestCase("sg_testing2", "config_object.json", "tf"),
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
	if outputFormat == "csv" {
		writer = csvio.NewWriter(buf)
	} else {
		writer = tfio.NewWriter(buf)
	}
	err = collection.Write(writer)
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
