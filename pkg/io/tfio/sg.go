/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tfio

import (
	"errors"
	"fmt"
	"strings"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteSG prints an entire collection of Security Groups as a sequence of terraform resources.
func (w *Writer) WriteSG(c *ir.SGCollection, vpc string, _ bool) error {
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

func sgCollection(collection *ir.SGCollection, vpc string) (*tf.ConfigFile, error) {
	var resources []tf.Block

	for _, vpcName := range collection.VpcNames() {
		if vpc != vpcName && vpc != "" {
			continue
		}
		for _, sgName := range collection.SortedSGNames(vpcName) {
			sgObject := collection.SGs[vpcName][sgName]
			sgTf, err := sg(sgObject, vpcName)
			if err != nil {
				return nil, err
			}
			resources = append(resources, sgTf)
			for i, rule := range sgObject.AllRules() {
				rule, err := sgRule(rule, sgName, i)
				if err != nil {
					return nil, err
				}
				resources = append(resources, rule)
			}
		}
	}
	return &tf.ConfigFile{
		Resources: resources,
	}, nil
}

func sg(sG *ir.SG, vpcName string) (tf.Block, error) {
	sgName := ir.ChangeScoping(sG.SGName.String())
	comment := fmt.Sprintf("\n### SG %s is attached to %s", sgName, strings.Join(sG.Targets, ", "))
	if len(sG.Targets) == 0 {
		comment = fmt.Sprintf("\n### SG %s is not attached to anything", sgName)
	}
	if err := verifyName(sgName); err != nil {
		return tf.Block{}, err
	}
	return tf.Block{
		Name:    "resource",
		Labels:  []string{quote("ibm_is_security_group"), quote(sgName)},
		Comment: comment,
		Arguments: []tf.Argument{
			{Name: "name", Value: quote("sg-" + sgName)},
			{Name: "resource_group", Value: "local.sg_synth_resource_group_id"},
			{Name: "vpc", Value: fmt.Sprintf("local.sg_synth_%s_id", vpcName)},
		},
	}, nil
}

func sgRule(rule *ir.SGRule, sgName ir.SGName, i int) (tf.Block, error) {
	ruleName := fmt.Sprintf("%s-%v", ir.ChangeScoping(sgName.String()), i)
	if err := verifyName(ruleName); err != nil {
		return tf.Block{}, err
	}

	group, err1 := value(sgName)
	remote, err2 := value(rule.Remote)
	if err := errors.Join(err1, err2); err != nil {
		return tf.Block{}, err
	}

	comment := ""
	if rule.Explanation != "" {
		comment = fmt.Sprintf("# %v", rule.Explanation)
	}

	return tf.Block{
		Name:    "resource", //nolint:revive  // obvious false positive
		Labels:  []string{quote("ibm_is_security_group_rule"), ir.ChangeScoping(quote(ruleName))},
		Comment: comment,
		Arguments: []tf.Argument{
			{Name: "group", Value: group},
			{Name: "direction", Value: quote(direction(rule.Direction))},
			{Name: "remote", Value: remote},
		},
		Blocks: sgProtocol(rule.Protocol),
	}, nil
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

func value(x interface{}) (string, error) {
	switch v := x.(type) {
	case *netset.IPBlock:
		return quote(v.String()), nil
	case ir.SGName:
		return ir.ChangeScoping(fmt.Sprintf("ibm_is_security_group.%v.id", v)), nil
	}
	return "", fmt.Errorf("invalid terraform value %v (type %T)", x, x)
}
