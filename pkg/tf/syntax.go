// Package tf represents the general syntax of terraform files
package tf

import "fmt"

// Terminology inspired by
// * https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md
// * https://developer.hashicorp.com/terraform/language/syntax/configuration
//
// This part knows nothing about ACLs

type Argument struct {
	Name  string
	Value string
}

type Block struct {
	Name      string
	Labels    []string
	Arguments []Argument
	Blocks    []Block
}

type ConfigFile struct {
	Resources []Block
}

const indentation = "  "

func (b *Block) print(indent string) string {
	result := ""
	result += indent + b.Name
	for _, label := range b.Labels {
		result += " " + label
	}
	result += " {\n"
	{
		indent := indent + indentation //nolint:govet  // intentionally shadow
		for _, keyValue := range b.Arguments {
			result += indent + fmt.Sprintf("%v = %v\n", keyValue.Name, keyValue.Value)
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
	return result
}
