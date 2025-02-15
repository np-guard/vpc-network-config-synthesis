/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package tfio implements output of ACLs in terraform format
package tfio

import (
	"fmt"
	"slices"
	"strings"

	"github.com/np-guard/models/pkg/netp"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteACL prints an entire collection of acls as a sequence of terraform resources.
func (w *Writer) WriteACL(c *ir.ACLCollection, vpc string, _ bool) error {
	collection, err := aclCollection(c, vpc)
	if err != nil {
		return err
	}
	if _, err := w.w.WriteString(collection.Print()); err != nil {
		return err
	}
	return w.w.Flush()
}

func aclCollection(collection *ir.ACLCollection, vpc string) (*tf.ConfigFile, error) {
	res := make([]tf.Block, 0)
	for _, vpcName := range collection.VpcNames() {
		if vpc != vpcName && vpc != "" {
			continue
		}
		for _, aclName := range collection.SortedACLNames(vpcName) {
			acl := collection.ACLs[vpcName][aclName]
			aclBlock, err := singleACL(acl, vpcName)
			if err != nil {
				return nil, err
			}
			res = append(res, aclBlock)
		}
	}
	return &tf.ConfigFile{
		Resources: res,
	}, nil
}

func singleACL(acl *ir.ACL, vpcName string) (tf.Block, error) {
	rules := acl.Rules()
	blocks := make([]tf.Block, len(rules))
	for i, rule := range rules {
		rule, err := aclRule(rule, fmt.Sprintf("rule%v", i))
		if err != nil {
			return tf.Block{}, err
		}
		blocks[i] = rule
	}
	aclName := ir.ChangeScoping(acl.Name)
	if err := verifyName(aclName); err != nil {
		return tf.Block{}, err
	}
	return tf.Block{
		Comment: aclComment(acl),
		Name:    resourceConst,
		Labels:  []string{quote("ibm_is_network_acl"), quote(aclName)},
		Arguments: []tf.Argument{
			{Name: nameConst, Value: quote(aclName)},
			{Name: "resource_group", Value: "local.acl_synth_resource_group_id"},
			{Name: "vpc", Value: fmt.Sprintf("local.acl_synth_%s_id", vpcName)},
		},
		Blocks: blocks,
	}, nil
}

func aclRule(rule *ir.ACLRule, name string) (tf.Block, error) {
	if err := verifyName(name); err != nil {
		return tf.Block{}, err
	}
	arguments := []tf.Argument{
		{Name: nameConst, Value: quote(name)},
		{Name: "action", Value: quote(action(rule.Action))},
		{Name: "direction", Value: quote(direction(rule.Direction))},
		{Name: "source", Value: quote(rule.Source.String())},
		{Name: "destination", Value: quote(rule.Destination.String())},
	}

	comment := ""
	if rule.Explanation != "" {
		comment = fmt.Sprintf("# %v", rule.Explanation)
	}

	return tf.Block{Name: "rules",
		Comment:   comment,
		Arguments: arguments,
		Blocks:    aclProtocol(rule.Protocol),
	}, nil
}

func aclProtocol(t netp.Protocol) []tf.Block {
	switch p := t.(type) {
	case netp.TCPUDP:
		return []tf.Block{{
			Name:      strings.ToLower(string(p.ProtocolString())),
			Arguments: slices.Concat(portRange(p.DstPorts(), "port"), portRange(p.SrcPorts(), "source_port")),
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

func aclComment(acl *ir.ACL) string {
	if len(acl.Subnets) == 0 {
		return "\n# No attached subnets"
	}
	return fmt.Sprintf("\n# Attached subnets: %s", acl.AttachedSubnetsString())
}
