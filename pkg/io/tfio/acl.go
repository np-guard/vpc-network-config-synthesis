/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package tfio implements output of ACLs in terraform format
package tfio

import (
	"fmt"
	"strings"

	"github.com/np-guard/models/pkg/ipblock"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteACL prints an entire collection of acls as a sequence of terraform resources.
func (w *Writer) WriteACL(c *ir.ACLCollection, vpc string) error {
	output := aclCollection(c, vpc).Print()
	_, err := w.w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
}

func aclProtocol(t ir.Protocol) []tf.Block {
	switch p := t.(type) {
	case ir.TCPUDP:
		return []tf.Block{{
			Name: strings.ToLower(string(p.Protocol)),
			Arguments: append(
				portRange(p.PortRangePair.DstPort, "port"),
				portRange(p.PortRangePair.SrcPort, "source_port")...,
			),
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
		Blocks:    aclProtocol(rule.Protocol),
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
		Labels:  []string{quote("ibm_is_network_acl"), changeScoping(quote(t.Name()))},
		Arguments: []tf.Argument{
			{Name: "name", Value: changeScoping(quote(t.Name()))}, //nolint:revive  // obvious false positive
			{Name: "resource_group", Value: "local.acl_synth_resource_group_id"},
			{Name: "vpc", Value: fmt.Sprintf("local.name_%s_id", vpcFromScopedResource(t.Subnet))},
		},
		Blocks: blocks,
	}
}

func aclCollection(t *ir.ACLCollection, vpc string) *tf.ConfigFile {
	sortedACLs := t.SortedACLSubnets(vpc)
	var acls = make([]tf.Block, len(sortedACLs))
	i := 0
	for _, subnet := range sortedACLs {
		comment := ""
		vpcName := ir.ScopingComponents(subnet)[0]
		acl := t.ACLs[vpcName][subnet]
		if len(sortedACLs) > 1 { // single nacl
			comment = fmt.Sprintf("\n# %v [%v]", subnet, subnetCidr(acl))
		}
		acls[i] = singleACL(acl, comment)
		i += 1
	}
	return &tf.ConfigFile{
		Resources: acls,
	}
}

func subnetCidr(acl *ir.ACL) *ipblock.IPBlock {
	if len(acl.Internal) > 0 {
		return acl.Internal[0].Target()
	}
	return acl.External[0].Target()
}
