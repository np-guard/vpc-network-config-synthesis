/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package tfio generates a locals.tf file that sets up some of the tf variables
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
	return data, w.Flush()
}

func locals(defs *ir.ConfigDefs, acl bool) string {
	result := "locals {\n"
	for _, vpcName := range utils.SortedMapKeys(defs.VPCs) {
		line := indentation + fmt.Sprintf("name_%s_id = <%s ID>", vpcName, vpcName) + "\n"
		result += line
	}
	prefix := "sg"
	if acl {
		prefix = "acl"
	}
	result += "\n" + indentation + fmt.Sprintf("%s_synth_resource_group_id = <RESOURCE-GROUP ID>", prefix) + "\n"
	result += "}"
	return result
}
