/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/np-guard/vpc-network-config-synthesis/cmd/subcmds"
)

func TestMain(t *testing.T) {
	for _, tt := range allMainTests() {
		t.Run(tt.testName, func(t *testing.T) {
			// create a sub folder
			if err := os.MkdirAll(filepath.Join(resultsFolder, tt.testName), defaultDirectoryPermission); err != nil {
				t.Errorf("Bad test %s: %s", tt.testName, err)
			}

			// run command
			cmd := fmt.Sprintf(tt.command, dataFolder, dataFolder, resultsFolder)
			if err := subcmds.Main(strings.Split(cmd, " ")); err != nil {
				t.Errorf("Bad test %s: %s", tt.testName, err)
			}

			// compare results
			compareTestResults(t, tt.testName)
		})
	}
	removeGeneratedFiles()
}

func compareTestResults(t *testing.T, testName string) {
	expectedSubDirPath := filepath.Join(expectedFolder, testName)
	resultsSubDirPath := filepath.Join(resultsFolder, testName)

	expectedDirFiles := readDir(t, expectedSubDirPath)
	resultsDirFiles := readDir(t, resultsSubDirPath)

	if len(expectedDirFiles) != len(resultsDirFiles) {
		t.Fatalf("Bad test: %s", testName)
	}

	for _, file := range expectedDirFiles {
		if readFile(t, filepath.Join(expectedSubDirPath, file)) != readFile(t, filepath.Join(resultsSubDirPath, file)) {
			t.Fatalf("Bad test: %s", testName)
		}
	}
}

func readDir(t *testing.T, dirName string) []string {
	entries, err := os.ReadDir(dirName)
	if err != nil {
		t.Errorf("Bad test %s: %s", dirName, err)
	}

	result := make([]string, len(entries))
	for i := range entries {
		result[i] = entries[i].Name()
	}
	return result
}

func readFile(t *testing.T, file string) string {
	buf, err := os.ReadFile(file)
	if err != nil {
		t.Errorf("Bad test: %v", err)
	}
	return shrinkWhitespace(string(buf))
}

func shrinkWhitespace(s string) string {
	return regexp.MustCompile(`[ \t]+`).ReplaceAllString(s, " ")
}

func removeGeneratedFiles() {
	err := os.RemoveAll(resultsFolder)
	if err != nil {
		panic(err)
	}
}
