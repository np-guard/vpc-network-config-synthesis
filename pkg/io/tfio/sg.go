/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tfio

import (
	"fmt"
	"log"
	"strings"

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteSG prints an entire collection of Security Groups as a sequence of terraform resources.
func (w *Writer) WriteSynthSG(c *ir.SGCollection, vpc string) error {
	output := sgCollection(c, vpc).Print()
	_, err := w.w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
}

func (w *Writer) WriteOptimizeSG(c *ir.SGCollection) error {
	return nil
}

func value(x interface{}) string {
	switch v := x.(type) {
	case *netset.IPBlock:
		return quote(v.String())
	case ir.SGName:
		return ir.ChangeScoping(fmt.Sprintf("ibm_is_security_group.%v.id", v))
	default:
		log.Fatalf("invalid terraform value %v", v)
	}
	return ""
}

func sgProtocol(t netp.Protocol, d ir.Direction) []tf.Block {
	switch p := t.(type) {
	case netp.TCPUDP:
		var remotePort interval.Interval
		if d == ir.Inbound {
			remotePort = p.SrcPorts()
		} else {
			remotePort = p.DstPorts()
		}
		return []tf.Block{{
			Name:      strings.ToLower(string(p.ProtocolString())),
			Arguments: portRange(remotePort, "port"),
		}}
	case netp.ICMP:
		return []tf.Block{{
			Name:      "icmp",
			Arguments: codeTypeArguments(p.ICMPTypeCode()),
		}}
	case netp.AnyProtocol:
		return []tf.Block{}
	}
	return nil
}

func sgRule(rule *ir.SGRule, sgName ir.SGName, i int) tf.Block {
	ruleName := fmt.Sprintf("%v-%v", sgName, i)
	verifyName(ruleName)
	return tf.Block{
		Name:    "resource",
		Labels:  []string{quote("ibm_is_security_group_rule"), ir.ChangeScoping(quote(ruleName))},
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
		Labels:  []string{quote("ibm_is_security_group"), ir.ChangeScoping(quote(sgName))},
		Comment: comment,
		Arguments: []tf.Argument{
			{Name: "name", Value: ir.ChangeScoping(quote("sg-" + sgName))},
			{Name: "resource_group", Value: "local.sg_synth_resource_group_id"},
			{Name: "vpc", Value: fmt.Sprintf("local.sg_synth_%s_id", ir.VpcFromScopedResource(sgName))},
		},
	}
}

func sgCollection(t *ir.SGCollection, vpc string) *tf.ConfigFile {
	var resources []tf.Block //nolint:prealloc  // nontrivial to calculate, and an unlikely performance bottleneck
	for _, sgName := range t.SortedSGNames(vpc) {
		comment := ""
		vpcName := ir.VpcFromScopedResource(string(sgName))
		rules := t.SGs[vpcName][sgName].AllRules()
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
