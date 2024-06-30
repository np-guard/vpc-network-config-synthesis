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
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const indentation = "  "

func WriteLocals(defs *ir.ConfigDefs, acl bool) (*bytes.Buffer, error) {
	data := new(bytes.Buffer)
	w := bufio.NewWriter(data)

	output := locals(defs, acl)
	if _, err := w.WriteString(output); err != nil {
		return nil, err
	}
	err := w.Flush()
	return data, err
}

func locals(defs *ir.ConfigDefs, acl bool) string {
	result := "locals {\n"
	prefix := "sg"
	if acl {
		prefix = "acl"
	}
	for _, vpcName := range utils.SortedMapKeys(defs.VPCs) {
		line := indentation + prefix + fmt.Sprintf("_synth_%s_id = <%s ID>", vpcName, vpcName) + "\n"
		result += line
	}

	result += "\n" + indentation + fmt.Sprintf("%s_synth_resource_group_id = <RESOURCE-GROUP ID>", prefix) + "\n"
	result += "}"
	return result
}
