/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

type testCase struct {
	testName    string
	command     string
	expectedErr string
}

const (
	dataFolder       = "data"
	dataErrorsFolder = "data_for_testing_errors"
	resultsFolder    = "results"
	expectedFolder   = "expected"

	defaultDirectoryPermission = 0o755
)
