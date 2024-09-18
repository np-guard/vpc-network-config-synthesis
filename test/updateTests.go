/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func updateTests() {
	for _, testCase := range allMainTests() {
		// create a sub folder
		if err := os.MkdirAll(filepath.Join(resultsFolder, testCase.testName), defaultDirectoryPermission); err != nil {
			handleErrors(t, testCase.testName, err)
		}

		cmd := strings.ReplaceAll(testCase.command, resultsFolder, expectedFolder)

		//create sub d

		_ = exec.Command("bash", "-c", cmd).Run()
	}
}
