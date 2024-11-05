/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/np-guard/vpc-network-config-synthesis/cmd/subcmds"
)

// comment lines 18-20 and uncomment `update_test.go` file to update all test outputs
func TestMain(t *testing.T) {
	testMain(t)
}

func testMain(t *testing.T) {
	for _, tt := range allMainTests() {
		t.Run(tt.testName, func(t *testing.T) {
			// create a sub folder
			if err := os.MkdirAll(filepath.Join(resultsFolder, tt.testName), defaultDirectoryPermission); err != nil {
				t.Fatalf("Bad test %s; error creating folder for results: %v", tt.testName, err)
			}

			// run command
			warning, err := subcmds.Main(tt.args.Args(dataFolder, resultsFolder))
			if err != nil {
				t.Fatalf("Bad test %s; unexpected err: %v", tt.testName, err)
			}

			if tt.expectedWarning != nil && *tt.expectedWarning != warning {
				t.Errorf("Bad test %s; blocked resources warning is different than expected; \n expected: %s got: %s", tt.testName,
					*tt.expectedWarning, warning)
			}

			// compare results
			compareTestResults(t, tt.testName)
		})
	}
	removeGeneratedFiles()
}

func compareTestResults(t *testing.T, testName string) {
	expectedSubDirPath := filepath.Join(expectedFolder, testName)
	expectedDirFiles := readDir(t, expectedSubDirPath)
	expectedFileNames := strings.Join(expectedDirFiles, ", ")

	resultsSubDirPath := filepath.Join(resultsFolder, testName)
	resultsDirFiles := readDir(t, resultsSubDirPath)
	resultsFileNames := strings.Join(resultsDirFiles, ", ")

	if len(expectedDirFiles) != len(resultsDirFiles) {
		t.Fatalf("Bad test: %s; incorrect number of files created.\nexpected: %s\ngot: %s", testName, expectedFileNames, resultsFileNames)
	}

	for _, file := range expectedDirFiles {
		if readFile(t, filepath.Join(expectedSubDirPath, file), testName) != readFile(t, filepath.Join(resultsSubDirPath, file), testName) {
			t.Fatalf("Bad test %s; The %s file is different than expected", testName, file)
		}
	}
}

func readDir(t *testing.T, dirName string) []string {
	entries, err := os.ReadDir(dirName)
	if err != nil {
		t.Fatalf("Bad test %s; error reading folder: %s", dirName, err)
	}

	result := make([]string, len(entries))
	for i := range entries {
		result[i] = entries[i].Name()
	}
	return result
}

func readFile(t *testing.T, file, testName string) string {
	buf, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Bad test: %s; error reading file %s: %v", testName, file, err)
	}
	return string(buf)
}

func removeGeneratedFiles() {
	err := os.RemoveAll(resultsFolder)
	if err != nil {
		panic(err)
	}
}
