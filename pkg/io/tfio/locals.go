/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tfio

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const indentation = "  "

// generate a locals.tf file that sets up some of the tf variables
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
  prefix := "sg"
	if acl {
		prefix = "acl"
	}
  result := []string{"locals {"}
	for _, vpcName := range utils.SortedMapKeys(defs.VPCs) {
		line := indentation + fmt.Sprintf("_synth_%s_id = <%s ID>", vpcName, vpcName)
		result = append(result, line)
	}
	
	result = append(result, "") // empty line
	line := indentation + fmt.Sprintf("%s_synth_resource_group_id = <RESOURCE-GROUP ID>\n}\n", prefix)
	result = append(result, line)
	return strings.Join(result, "\n")
}
