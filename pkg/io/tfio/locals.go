/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package tfio generates a locals.tf file that sets up all the local variables that
// were used in the generated SGs/nACLs.
package tfio

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const indentation = "  "

func WriteLocals(defs *ir.ConfigDefs, acl bool) error {
	data := new(bytes.Buffer)
	w := bufio.NewWriter(data)
	// writer := NewWriter(w)

	output := locals(defs, acl)
	_, err := w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.Flush()

	// if err := os.WriteFile(args.outputFile, data.Bytes(), defaultFilePermission); err != nil {
	//		return err
	//	}
	return err
}

func locals(defs *ir.ConfigDefs, acl bool) string {
	result := "locals {\n"
	for vpcName := range defs.VPCs {
		line := indentation + fmt.Sprintf("name_%s_id = <%s ID>", vpcName, vpcName) + "\n"
		result = result + line
	}
	prefix := "sg"
	if acl {
		prefix = "acl"
	}
	result = result + "\n" + indentation + fmt.Sprintf("%s_synth_resource_group_id = <RESOURCE-GROUP ID>", prefix) + "\n"
	result = result + "}"
	return result
}
