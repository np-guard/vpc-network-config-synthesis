package tfio

import (
	"fmt"
	"log"
	"strings"

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
		return fmt.Sprintf("ibm_is_security_group.%v.id", v)
	default:
		log.Fatalf("invalid terraform value %v", v)
	}
	return ""
}

func sgProtocol(t ir.Protocol, d ir.Direction) []tf.Block {
	switch p := t.(type) {
	case ir.TCPUDP:
		var remotePort ir.PortRange
		if d == ir.Inbound {
			remotePort = p.PortRangePair.SrcPort
		} else {
			remotePort = p.PortRangePair.DstPort
		}
		return []tf.Block{{
			Name:      strings.ToLower(string(p.Protocol)),
			Arguments: portRange(remotePort, "port"),
		}}
	case ir.ICMP:
		return []tf.Block{{
			Name:      "icmp",
			Arguments: codeTypeArguments(p.ICMPCodeType),
		}}
	case ir.AnyProtocol:
		return []tf.Block{}
	}
	return nil
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
		Blocks: sgProtocol(rule.Protocol, rule.Direction),
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
