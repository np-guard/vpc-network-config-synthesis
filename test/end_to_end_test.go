/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"testing"
)

const (
	dataFolder     = "data"
	expectedFolder = "expected"
	resultsFolder  = "results"

	defaultDirectoryPermission = 0o755
)

type testCase struct {
	testName string
	command  string
}

func TestMain(t *testing.T) {
	var wg sync.WaitGroup

	for _, testCase := range allMainTests() {
		wg.Add(1)

		t.Run(testCase.testName, func(t *testing.T) {
			// run all tests in parallel
			t.Parallel()

			// create a sub folder
			if err := os.MkdirAll(filepath.Join(resultsFolder, testCase.testName), defaultDirectoryPermission); err != nil {
				handleErrors(t, testCase.testName, err)
			}

			// run the command
			err := exec.Command("bash", "-c", testCase.command).Run()
			if err != nil {
				log.Fatalf("Bad test %s: %s", testCase.testName, err) // should not occur
			}

			// compare results
			compareSubDirs(t, testCase.testName)

			wg.Done()
		})
	}
	// wait for all tests to complete
	wg.Wait()

	removeGeneratedFiles()
}

func compareSubDirs(t *testing.T, testName string) {
	expectedSubDirPath := filepath.Join(expectedFolder, testName)
	resultsSubDirPath := filepath.Join(resultsFolder, testName)

	for _, file := range readDir(t, expectedSubDirPath) {
		if readFile(filepath.Join(expectedSubDirPath, file)) != readFile(filepath.Join(resultsSubDirPath, file)) {
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

func readFile(file string) string {
	buf, err := os.ReadFile(file)
	if err != nil {
		log.Panicf("Bad test: %v", err)
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
