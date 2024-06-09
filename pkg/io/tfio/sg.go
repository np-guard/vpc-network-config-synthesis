/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tfio

import (
	"fmt"
	"log"
	"strings"

	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteSG prints an entire collection of security groups as a sequence of terraform resources.
func (w *Writer) WriteSG(c *ir.SGCollection, vpc string) error {
	output := sgCollection(c, vpc).Print()
	_, err := w.w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
}

func value(x interface{}) string {
	switch v := x.(type) {
	case *ipblock.IPBlock:
		return quote(v.String())
	case ir.SGName:
		return changeScoping(fmt.Sprintf("ibm_is_security_group.%v.id", v))
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
		Labels:  []string{quote("ibm_is_security_group_rule"), changeScoping(quote(ruleName))},
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
		Labels:  []string{quote("ibm_is_security_group"), changeScoping(quote(sgName))},
		Comment: comment,
		Arguments: []tf.Argument{
			{Name: "name", Value: changeScoping(quote("sg-" + sgName))},
			{Name: "resource_group", Value: "local.sg_synth_resource_group_id"},
			{Name: "vpc", Value: fmt.Sprintf("local.name_%s_id", ir.VpcFromScopedResource(sgName))},
		},
	}
}

func sgCollection(t *ir.SGCollection, vpc string) *tf.ConfigFile {
	var resources []tf.Block //nolint:prealloc  // nontrivial to calculate, and an unlikely performance bottleneck
	for _, sgName := range t.SortedSGNames(vpc) {
		comment := ""
		vpcName := ir.VpcFromScopedResource(string(sgName))
		rules := t.SGs[vpcName][sgName].Rules
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
