/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tfio

import (
	"errors"
	"fmt"
	"strings"

	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteSG prints an entire collection of Security Groups as a sequence of terraform resources.
func (w *Writer) WriteSG(c *ir.SGCollection, vpc string) error {
	collection, err := sgCollection(c, vpc)
	if err != nil {
		return err
	}
	output := collection.Print()
	_, err = w.w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
}

func value(x interface{}) (string, error) {
	switch v := x.(type) {
	case *ipblock.IPBlock:
		return quote(v.String()), nil
	case ir.SGName:
		return ir.ChangeScoping(fmt.Sprintf("ibm_is_security_group.%v.id", v)), nil
	}
	return "", fmt.Errorf("invalid terraform value %v (type %T)", x, x)
}

func sgProtocol(t ir.Protocol) []tf.Block {
	switch p := t.(type) {
	case ir.TCPUDP:
		return []tf.Block{{
			Name:      strings.ToLower(string(p.Protocol)),
			Arguments: portRange(p.PortRangePair.DstPort, "port"),
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

func sgRule(rule *ir.SGRule, sgName ir.SGName, i int) (tf.Block, error) {
	ruleName := fmt.Sprintf("%v-%v", sgName, i)
	if err := verifyName(ruleName); err != nil {
		return tf.Block{}, err
	}

	group, err1 := value(sgName)
	remote, err2 := value(rule.Remote)
	if err := errors.Join(err1, err2); err != nil {
		return tf.Block{}, err
	}

	return tf.Block{
		Name:    "resource",
		Labels:  []string{quote("ibm_is_security_group_rule"), ir.ChangeScoping(quote(ruleName))},
		Comment: fmt.Sprintf("# %v", rule.Explanation),
		Arguments: []tf.Argument{
			{Name: "group", Value: group},
			{Name: "direction", Value: quote(direction(rule.Direction))},
			{Name: "remote", Value: remote},
		},
		Blocks: sgProtocol(rule.Protocol),
	}, nil
}

func sg(sgName, comment string) (tf.Block, error) {
	if err := verifyName(sgName); err != nil {
		return tf.Block{}, err
	}
	return tf.Block{
		Name:    "resource", //nolint:revive  // obvious false positive
		Labels:  []string{quote("ibm_is_security_group"), ir.ChangeScoping(quote(sgName))},
		Comment: comment,
		Arguments: []tf.Argument{
			{Name: "name", Value: ir.ChangeScoping(quote("sg-" + sgName))},
			{Name: "resource_group", Value: "local.sg_synth_resource_group_id"},
			{Name: "vpc", Value: fmt.Sprintf("local.sg_synth_%s_id", ir.VpcFromScopedResource(sgName))},
		},
	}, nil
}

func sgCollection(t *ir.SGCollection, vpc string) (*tf.ConfigFile, error) {
	var resources []tf.Block //nolint:prealloc  // nontrivial to calculate, and an unlikely performance bottleneck
	for _, sgName := range t.SortedSGNames(vpc) {
		comment := ""
		vpcName := ir.VpcFromScopedResource(string(sgName))
		rules := t.SGs[vpcName][sgName].AllRules()
		comment = fmt.Sprintf("\n### SG attached to %v", sgName)
		sg, err := sg(sgName.String(), comment)
		if err != nil {
			return nil, err
		}
		resources = append(resources, sg)
		for i, rule := range rules {
			rule, err := sgRule(rule, sgName, i)
			if err != nil {
				return nil, err
			}
			resources = append(resources, rule)
		}
	}
	return &tf.ConfigFile{
		Resources: resources,
	}, nil
}
