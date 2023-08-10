// Package tf represents the general syntax of terraform files
package tf

import "fmt"

// Terminology inspired by
// * https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md
// * https://developer.hashicorp.com/terraform/language/syntax/configuration
//
// This part knows nothing about ACLs

type Block struct {
	Name      string
	Labels    []string
	Arguments map[string]string
	Blocks    []Block
}

type ConfigFile struct {
	Resources []Block
}

type Blockable interface {
	Terraform() Block
}

func Blocks[T Blockable](items []T) []Block {
	result := make([]Block, len(items))
	for i := range items {
		result[i] = items[i].Terraform()
	}
	return result
}

func (b *Block) print(indent string) string {
	result := ""
	result += indent + b.Name
	for _, label := range b.Labels {
		result += " " + label
	}
	result += " {\n"
	{
		indent := indent + "    " //nolint:govet  // intentionally shadow
		for key, value := range b.Arguments {
			result += indent + fmt.Sprintf("%v = %v\n", key, value)
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