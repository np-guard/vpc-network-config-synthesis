/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

type mainTestCase struct {
	testName string
	command  string
}

type errorTestCase struct {
	testName string
	command  string
	err      string
}

const (
	dataFolder       = "data"
	dataErrorsFolder = "data_errors"
	resultsFolder    = "results"
	expectedFolder   = "expected"

	defaultDirectoryPermission = 0o755
)
