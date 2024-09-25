/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"strings"
	"testing"

	"github.com/np-guard/vpc-network-config-synthesis/cmd/subcmds"
)

func TestErrors(t *testing.T) {
	for _, tt := range errorTestsList() {
		t.Run(tt.testName, func(t *testing.T) {
			// run command
			err := subcmds.Main(strings.Split(tt.command, " "))
			if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
				res := "nil"
				if err != nil {
					res = err.Error()
				}
				t.Errorf("Bad test %s:\nexpected err: %s\ngot err: %s", tt.testName, tt.expectedErr, res)
			}
		})
	}
}
