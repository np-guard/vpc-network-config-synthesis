/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import "fmt"

type testCase struct {
	testName    string
	args        *command
	expectedErr string
}

type command struct {
	cmd        string
	subcmd     string
	singleacl  bool
	config     string
	spec       string
	outputFile string
	outputDir  string
	prefix     string
	format     string
	locals     bool
}

const (
	dataFolder                 string = "data"
	dataForTestingErrorsFolder string = "data_for_testing_errors"
	resultsFolder              string = "results"
	expectedFolder             string = "expected"

	defaultDirectoryPermission = 0o755

	synth string = "synth"
	acl   string = "acl"
	sg    string = "sg"
)

func (c *command) Args(dataFolder, resultsFolder string) []string {
	res := []string{"./bin/vpcgen"}

	if c.cmd != "" {
		res = append(res, c.cmd)
	}
	if c.subcmd != "" {
		res = append(res, c.subcmd)
	}
	if c.singleacl {
		res = append(res, "--single")
	}
	if c.config != "" {
		res = append(res, "-c", fmt.Sprintf(c.config, dataFolder))
	}
	if c.spec != "" {
		res = append(res, "-s", fmt.Sprintf(c.spec, dataFolder))
	}
	if c.outputFile != "" {
		res = append(res, "-o", fmt.Sprintf(c.outputFile, resultsFolder))
	}
	if c.outputDir != "" {
		res = append(res, "-d", fmt.Sprintf(c.outputDir, resultsFolder))
	}
	if c.format != "" {
		res = append(res, "-f", c.format)
	}
	if c.prefix != "" {
		res = append(res, "-p", c.prefix)
	}
	if c.locals {
		res = append(res, "-l")
	}

	return res
}
