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

const dataFolder = "data/"

func TestInvalidConnectivity(t *testing.T) {
	folder := dataFolder + "acl_testing4"
	_, err := makeACL(folder)
	if err.Error() != "unsupported endpoint type cidr" {
		t.Errorf("No failure for unsupported type; got %v", err)
	}
}

func Test_single_conn1(t *testing.T) {
	subfolder := "single_conn1"
	actualCSVString, err := makeACL(dataFolder + subfolder)
	if err != nil {
		t.Error(err)
	}
	expectedCSVString := readExpectedCSV(dataFolder + subfolder + "/nacl_single_expected.csv")
	if expectedCSVString != actualCSVString {
		t.Errorf("%v != %v", expectedCSVString, actualCSVString)
	}
}

func Test_single_conn2(t *testing.T) {
	subfolder := "single_conn2"
	actualCSVString, err := makeACL(dataFolder + subfolder)
	if err != nil {
		t.Error(err)
	}
	expectedCSVString := readExpectedCSV(dataFolder + subfolder + "/nacl_single_expected.csv")
	if expectedCSVString != actualCSVString {
		t.Errorf("%v != %v", expectedCSVString, actualCSVString)
	}
}

func Test_acl_testing5(t *testing.T) {
	subfolder := "acl_testing5"
	actualCSVString, err := makeACL(dataFolder + subfolder)
	if err != nil {
		t.Error(err)
	}
	expectedCSVString := readExpectedCSV(dataFolder + subfolder + "/nacl_single_expected.csv")
	if expectedCSVString != actualCSVString {
		t.Errorf("%v != %v", expectedCSVString, actualCSVString)
	}
}

func makeACL(folder string) (string, error) {
	reader := jsonio.NewReader()

	var subnets map[string]string
	configFilename := folder + "/config_object.json"
	subnets, err := jsonio.ReadSubnetMap(configFilename)
	if err != nil {
		return "", err
	}

	connectivityFilename := folder + "/conn_spec.json"
	s, err := reader.ReadSpec(connectivityFilename, subnets)
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
