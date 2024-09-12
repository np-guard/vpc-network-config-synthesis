/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tfio

import (
	"fmt"
	"log"
	"strings"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func (w *Writer) WriteSynthSG(c *ir.SGCollection, vpc string) error {
	return w.writeSGCollection(c, vpc, true)
}

func (w *Writer) WriteOptimizeSG(c *ir.SGCollection) error {
	return w.writeSGCollection(c, "", false)
}

// writeSGCollection prints an entire collection of Security Groups as a sequence of terraform resources.
func (w *Writer) writeSGCollection(c *ir.SGCollection, vpc string, writeComments bool) error {
	output := sgCollection(c, vpc, writeComments).Print()
	_, err := w.w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
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

func sgProtocol(t netp.Protocol) []tf.Block {
	switch p := t.(type) {
	case netp.TCPUDP:
		return []tf.Block{{
			Name:      strings.ToLower(string(p.ProtocolString())),
			Arguments: portRange(p.DstPorts(), "port"),
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

func sgRule(rule *ir.SGRule, sgName ir.SGName, i int, writeComment bool) tf.Block {
	ruleName := fmt.Sprintf("%v-%v", sgName, i)
	verifyName(ruleName)
	comment := ""
	if writeComment {
		comment = fmt.Sprintf("# %v", rule.Explanation)
	}
	return tf.Block{
		Name:    "resource",
		Labels:  []string{quote("ibm_is_security_group_rule"), ir.ChangeScoping(quote(ruleName))},
		Comment: comment,
		Arguments: []tf.Argument{
			{Name: "group", Value: value(sgName)},
			{Name: "direction", Value: quote(direction(rule.Direction))},
			{Name: "remote", Value: value(rule.Remote)},
		},
		Blocks: sgProtocol(rule.Protocol),
	}
}

func sgBlock(sg *ir.SG, comment string) tf.Block {
	sgName := sg.SGName.String()
	verifyName(sgName)
	return tf.Block{
		Name:    "resource", //nolint:revive  // obvious false positive
		Labels:  []string{quote("ibm_is_security_group"), ir.ChangeScoping(quote(sgName))},
		Comment: comment,
		Arguments: []tf.Argument{
			{Name: "name", Value: ir.ChangeScoping(quote("sg-" + sgName))},
			{Name: "resource_group", Value: "local.sg_synth_resource_group_id"},
			{Name: "vpc", Value: fmt.Sprintf("local.sg_synth_%s_id", sg.VpcName)},
		},
	}
}

func sgCollection(t *ir.SGCollection, vpc string, writeComments bool) *tf.ConfigFile {
	var resources []tf.Block
	for _, vpcName := range t.VpcNames() {
		if vpc != vpcName && vpc != "" {
			continue
		}
		for _, sgName := range t.SortedSGNames(vpcName) {
			sg := t.SGs[vpcName][sgName]
			comment := "\n"
			if writeComments {
				comment = fmt.Sprintf("\n### SG attached to %v", sgName)
			}
			resources = append(resources, sgBlock(sg, comment))
			rules := sg.AllRules()
			for i := range rules {
				resources = append(resources, sgRule(&rules[i], sgName, i, writeComments))
			}
		}
	}
	return &tf.ConfigFile{
		Resources: resources,
	}
}
