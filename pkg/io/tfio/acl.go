// Package tfio implements output of ACLs in terraform format
package tfio

import (
	"fmt"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteACL prints an entire collection of acls as a sequence of terraform resources.
func (w *Writer) WriteACL(c *ir.ACLCollection) error {
	output := aclCollection(c).Print()
	_, err := w.w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
}

func aclRule(rule *ir.ACLRule, name string) tf.Block {
	verifyName(name)
	arguments := []tf.Argument{
		{Name: "name", Value: quote(name)},
		{Name: "action", Value: quote(action(rule.Action))},
		{Name: "direction", Value: quote(direction(rule.Direction))},
		{Name: "source", Value: quote(rule.Source.String())},
		{Name: "destination", Value: quote(rule.Destination.String())},
	}
	return tf.Block{Name: "rules",
		Comment:   fmt.Sprintf("# %v", rule.Explanation),
		Arguments: arguments,
		Blocks:    protocol(rule.Protocol),
	}
}

func singleACL(t *ir.ACL, comment string) tf.Block {
	rules := t.Rules()
	blocks := make([]tf.Block, len(rules))
	for i := range rules {
		blocks[i] = aclRule(&rules[i], fmt.Sprintf("rule%v", i))
	}
	return tf.Block{
		Comment: comment,
		Name:    "resource",
		Labels:  []string{quote("ibm_is_network_acl"), quote(t.Name())},
		Arguments: []tf.Argument{
			{Name: "name", Value: quote(t.Name())}, //nolint:revive  // obvious false positive
			{Name: "resource_group", Value: "var.resource_group_id"},
			{Name: "vpc", Value: "var.vpc_id"},
		},
		Blocks: blocks,
	}
}

func aclCollection(t *ir.ACLCollection) *tf.ConfigFile {
	var acls = make([]tf.Block, len(t.ACLs))
	i := 0
	for _, subnet := range t.SortedACLSubnets() {
		comment := ""
		if len(acls) > 1 {
			comment = fmt.Sprintf("\n# %v [%v]", subnet, t.ACLs[subnet].Internal[0].Target())
		}
		acls[i] = singleACL(t.ACLs[subnet], comment)
		i += 1
	}
	return &tf.ConfigFile{
		Resources: acls,
	}
}
