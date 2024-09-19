/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	m "github.com/np-guard/vpc-network-config-synthesis/cmd/_vpcgen"
)

type testCase struct {
	testName string
	command  string
}

const (
	dataFolder     = "data"
	resultsFolder  = "results"
	expectedFolder = "expected"

	defaultDirectoryPermission = 0o755
)

func TestMain(t *testing.T) {
	for _, testCase := range allMainTests() {
		t.Run(testCase.testName, func(t *testing.T) {
			// create a sub folder
			if err := os.MkdirAll(filepath.Join(resultsFolder, testCase.testName), defaultDirectoryPermission); err != nil {
				handleErrors(t, testCase.testName, err)
			}

			// run command
			cmd := fmt.Sprintf(testCase.command, dataFolder, dataFolder, resultsFolder)
			if err := m.Main(strings.Split(cmd, " ")); err != nil {
				t.Errorf("Bad test %s: %s", testCase.testName, err)
			}

			// compare results
			compareTestResults(t, testCase.testName)
		})
	}
	removeGeneratedFiles()
}

func compareTestResults(t *testing.T, testName string) {
	expectedSubDirPath := filepath.Join(expectedFolder, testName)
	resultsSubDirPath := filepath.Join(resultsFolder, testName)

	for _, file := range readDir(t, expectedSubDirPath) {
		if readFile(t, filepath.Join(expectedSubDirPath, file)) != readFile(t, filepath.Join(resultsSubDirPath, file)) {
			t.Fatalf("Bad test: %s", testName)
		}
	}
}

func readDir(t *testing.T, dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		handleErrors(t, dir, errors.New("error reading "+dir+": "+err.Error()))
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

func handleErrors(t *testing.T, testName string, err error) {
	removeGeneratedFiles()
	t.Fatalf("Bad test: %s. error: %s", testName, err)
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
