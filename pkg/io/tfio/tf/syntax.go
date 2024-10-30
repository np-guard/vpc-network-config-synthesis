/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package tf represents the general syntax of terraform files
package tf

import (
	"fmt"
	"strings"
)

// Terminology inspired by
// * https://github.com/hashicorp/hcl/blob/main/hclsyntax/ir.md
// * https://developer.hashicorp.com/terraform/language/syntax/configuration
//
// This part knows nothing about ACLs or SGs

type Argument struct {
	Name  string
	Value string
}

type Block struct {
	Comment   string
	Name      string
	Labels    []string
	Arguments []Argument
	Blocks    []Block
}

type ConfigFile struct {
	Resources []Block
}

const indentation = "  "
const space = " "

func (b *Block) print(indent string) string {
	result := ""
	if b.Comment != "" {
		result += indent + fmt.Sprintf("%v\n", b.Comment)
	}
	result += indent + b.Name
	for _, label := range b.Labels {
		result += space + label
	}
	result += " {\n"
	{
		indent := indent + indentation
		longestName := calculateLongestKey(b.Arguments)
		for _, keyValue := range b.Arguments {
			internalIdent := strings.Repeat(space, longestName-len(keyValue.Name))
			result += indent + fmt.Sprintf("%v%s = %v\n", keyValue.Name, internalIdent, keyValue.Value)
		}
		for _, sub := range b.Blocks {
			result += sub.print(indent)
		}
	}
	result += indent + "}\n"
	return result
}

func (c *ConfigFile) Print() string {
	result := ""
	for _, block := range c.Resources {
		result += block.print("")
	}
	return strings.TrimSpace(result) + "\n"
}

func calculateLongestKey(arguments []Argument) int {
	longestKey := 0
	for _, arg := range arguments {
		longestKey = max(longestKey, len(arg.Name))
	}
	return longestKey
}
