/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"testing"

// 	"github.com/np-guard/vpc-network-config-synthesis/cmd/subcmds"
// )

// func TestUpdate(t *testing.T) {
// 	for _, tt := range allMainTests() {
// 		t.Run(tt.testName, func(t *testing.T) {
// 			// create a sub folder
// 			if err := os.MkdirAll(filepath.Join(expectedFolder, tt.testName), defaultDirectoryPermission); err != nil {
// 				t.Errorf("Bad test %s: %s", tt.testName, err)
// 			}

// 			cmd := fmt.Sprintf(tt.command, dataFolder, dataFolder, expectedFolder)
// 			err := subcmds.Main(strings.Split(cmd, " "))
// 			if err != nil {
// 				t.Errorf("Bad test %s: %s", tt.testName, err)
// 			}
// 		})
// 	}
// }
