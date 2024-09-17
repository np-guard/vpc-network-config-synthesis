/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"os/exec"
	"strings"
)

func updateTests() {
	for _, testCase := range allMainTests() {
		cmd := strings.ReplaceAll(testCase.command, resultsFolder, expectedFolder)
		_ = exec.Command("bash", "-c", cmd).Run()
	}
}
