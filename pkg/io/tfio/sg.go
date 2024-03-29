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
		return quote(v.String())
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
	ruleName := fmt.Sprintf("%v-%v", sgName, i)
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

func sg(sgName, comment string) tf.Block {
	verifyName(sgName)
	return tf.Block{
		Name:    "resource", //nolint:revive  // obvious false positive
		Labels:  []string{quote("ibm_is_security_group"), quote(sgName)},
		Comment: comment,
		Arguments: []tf.Argument{
			{Name: "name", Value: quote("sg-" + sgName)},
			{Name: "resource_group", Value: "local.sg_synth_resource_group_id"},
			{Name: "vpc", Value: "local.sg_synth_vpc_id"},
		},
	}
}

func sgCollection(t *ir.SGCollection) *tf.ConfigFile {
	var resources []tf.Block //nolint:prealloc  // nontrivial to calculate, and an unlikely performance bottleneck
	for _, sgName := range t.SortedSGNames() {
		comment := ""
		rules := t.SGs[sgName].Rules
		if len(rules) == 0 {
			continue
		}
		comment = fmt.Sprintf("\n### SG attached to %v", sgName)
		resources = append(resources, sg(sgName.String(), comment))
		for i := range rules {
			rule := sgRule(&rules[i], sgName, i)
			resources = append(resources, rule)
		}
	}
	return &tf.ConfigFile{
		Resources: resources,
	}
}
