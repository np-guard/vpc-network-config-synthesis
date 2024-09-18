/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	m "github.com/np-guard/vpc-network-config-synthesis/cmd/_vpcgen"
)

func TestUpdate(t *testing.T) {
	for _, testCase := range allMainTests() {
		t.Run(testCase.testName, func(t *testing.T) {
			// create a sub folder
			if err := os.MkdirAll(filepath.Join(expectedFolder, testCase.testName), defaultDirectoryPermission); err != nil {
				log.Printf("Bad test %s: %s", testCase.testName, err)
			}

			cmd := strings.ReplaceAll(testCase.command, resultsFolder, expectedFolder)
			_ = m.Main(strings.Split(cmd, " "))
		})
	}
	log.Printf("done")
}
