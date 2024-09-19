/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"fmt"
	"strings"
	"testing"

	m "github.com/np-guard/vpc-network-config-synthesis/cmd/_vpcgen"
)

func TestErrors(t *testing.T) {
	for _, tt := range errorTestsList() {
		t.Run(tt.testName, func(t *testing.T) {
			// run command
			cmd := fmt.Sprintf(tt.command, dataErrorsFolder, dataErrorsFolder, resultsFolder)
			err := m.Main(strings.Split(cmd, " "))
			strings.Contains("something", "some") // true
			if err == nil || !strings.Contains(err.Error(), tt.err) {
				res := "nil"
				if err != nil {
					res = err.Error()
				}
				t.Errorf("Bad test %s: expected: %s, got %s", tt.testName, tt.err, res)
			}
		})
	}

}
