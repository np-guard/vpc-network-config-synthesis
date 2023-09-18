// Package tfio implements output of ACLs in terraform format
package tfio

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteSG prints an entire collection of security groups as a sequence of terraform resources.
func (w *Writer) WriteSG(c *ir.SGCollection) error {
	output := sgCollection(c).Print()
	_, err := w.w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
}

func value(x interface{}) string {
	switch v := x.(type) {
	case ir.CIDR:
		return quote(string(v))
	case ir.IP:
		return quote(v.String())
	case ir.SGName:
		return string(v) + ".id"
	default:
		log.Fatalf("invalid terraform value %v", v)
	}
	return ""
}

func sgRule(rule *ir.SGRule, sgName ir.SGName, i int) tf.Block {
	ruleName := fmt.Sprintf("sgrule-%v-%v", sgName, i)
	verifyName(ruleName)
	return tf.Block{
		Name:    "resource",
		Labels:  []string{quote("ibm_is_security_group_rule"), quote(ruleName)},
		Comment: fmt.Sprintf("# %v", rule.Explanation),
		Arguments: []tf.Argument{
			{Name: "group", Value: value(sgName)},
			{Name: "direction", Value: quote(direction(rule.Direction))},
			{Name: "remote", Value: value(rule.Remote)},
		},
		Blocks: protocol(rule.Protocol),
	}
}

func sgCollection(t *ir.SGCollection) *tf.ConfigFile {
	var sgRules []tf.Block
	for _, sgName := range t.SortedSGNames() {
		rules := t.SGs[sgName].Rules
		for i := range rules {
			sgRule := sgRule(&rules[i], sgName, i)
			sgRules = append(sgRules, sgRule)
		}
	}
	return &tf.ConfigFile{
		Resources: sgRules,
	}
}
