package tf

import "fmt"

// Terminology inspired by
// * https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md
// * https://developer.hashicorp.com/terraform/language/syntax/configuration
//
// This part knows nothing about ACLs

type block struct {
	Name      string
	Labels    []string
	Arguments map[string]string
	Blocks    []block
}

type configFile struct {
	Resources []block
}

type blockable interface {
	terraform() block
}

func blocks[T blockable](items []T) []block {
	result := make([]block, len(items))
	for i := range items {
		result[i] = items[i].terraform()
	}
	return result
}

func (b *block) print(indent string) string {
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

func (c *configFile) print() string {
	result := ""
	for _, block := range c.Resources {
		result += block.print("")
	}
	return result
}
