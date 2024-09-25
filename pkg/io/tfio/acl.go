/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package tfio implements output of ACLs in terraform format
package tfio

import (
	"fmt"
	"strings"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// WriteACL prints an entire collection of acls as a sequence of terraform resources.
func (w *Writer) WriteACL(c *ir.ACLCollection, vpc string) error {
	collection, err := aclCollection(c, vpc)
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

func aclProtocol(t netp.Protocol) []tf.Block {
	switch p := t.(type) {
	case netp.TCPUDP:
		return []tf.Block{{
			Name: strings.ToLower(string(p.ProtocolString())),
			Arguments: append(
				portRange(p.DstPorts(), "port"),
				portRange(p.SrcPorts(), "source_port")...,
			),
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

func aclRule(rule *ir.ACLRule, name string) (tf.Block, error) {
	if err := verifyName(name); err != nil {
		return tf.Block{}, err
	}
	arguments := []tf.Argument{
		{Name: "name", Value: quote(name)},
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

func singleACL(t *ir.ACL, comment string) (tf.Block, error) {
	rules := t.Rules()
	blocks := make([]tf.Block, len(rules))
	for i := range rules {
		rule, err := aclRule(&rules[i], fmt.Sprintf("rule%v", i))
		if err != nil {
			return tf.Block{}, err
		}
		blocks[i] = rule
	}
	aclName := ir.ChangeScoping(t.Name())
	if err := verifyName(aclName); err != nil {
		return tf.Block{}, err
	}
	return tf.Block{
		Comment: comment,
		Name:    "resource",
		Labels:  []string{quote("ibm_is_network_acl"), quote(aclName)},
		Arguments: []tf.Argument{
			{Name: "name", Value: quote(aclName)}, //nolint:revive  // obvious false positive
			{Name: "resource_group", Value: "local.acl_synth_resource_group_id"},
			{Name: "vpc", Value: fmt.Sprintf("local.acl_synth_%s_id", ir.VpcFromScopedResource(t.Subnet))},
		},
		Blocks: blocks,
	}, nil
}

func aclCollection(t *ir.ACLCollection, vpc string) (*tf.ConfigFile, error) {
	sortedACLs := t.SortedACLSubnets(vpc)
	var acls = make([]tf.Block, len(sortedACLs))
	i := 0
	for _, subnet := range sortedACLs {
		comment := ""
		vpcName := ir.VpcFromScopedResource(subnet)
		acl := t.ACLs[vpcName][subnet]
		if len(sortedACLs) > 1 { // not a single nacl
			comment = fmt.Sprintf("\n# %v [%v]", subnet, subnetCidr(acl))
		}
		singleACL, err := singleACL(acl, comment)
		if err != nil {
			return nil, err
		}
		acls[i] = singleACL
		i += 1
	}
	return &tf.ConfigFile{
		Resources: acls,
	}, nil
}

func subnetCidr(acl *ir.ACL) *netset.IPBlock {
	if len(acl.Internal) > 0 {
		return acl.Internal[0].Target()
	}
	return acl.External[0].Target()
}
