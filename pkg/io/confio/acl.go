/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

func cidr(address *netset.IPBlock) *string {
	return utils.Ptr(address.ToCidrListString())
}

func makeACLRuleItem(rule *ir.ACLRule, current,
	next *vpcv1.NetworkACLRuleReference) (vpcv1.NetworkACLRuleItemIntf, error) {
	iPVersion := utils.Ptr(ipv4Const)
	direction := direction(rule.Direction)
	action := action(rule.Action)
	source := cidr(rule.Source)
	destination := cidr(rule.Destination)
	switch p := rule.Protocol.(type) {
	case netp.TCPUDP:
		data := tcpudp(p)
		result := &vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp{
			Href:        current.Href,
			ID:          current.ID,
			Name:        current.Name,
			IPVersion:   iPVersion,
			Action:      action,
			Direction:   direction,
			Source:      source,
			Destination: destination,
			Before:      next,

			Protocol:           data.Protocol,
			SourcePortMin:      data.srcPortMin,
			SourcePortMax:      data.srcPortMax,
			DestinationPortMin: data.dstPortMin,
			DestinationPortMax: data.dstPortMax,
		}
		return result, nil
	case netp.ICMP:
		data := icmp(p)
		result := &vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp{
			Href:        current.Href,
			ID:          current.ID,
			Name:        current.Name,
			IPVersion:   iPVersion,
			Action:      action,
			Direction:   direction,
			Source:      source,
			Destination: destination,
			Before:      next,

			Protocol: data.Protocol,
			Type:     data.Type,
			Code:     data.Code,
		}
		return result, nil
	case netp.AnyProtocol:
		data := all()
		result := &vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll{
			Href:        current.Href,
			ID:          current.ID,
			Name:        current.Name,
			IPVersion:   iPVersion,
			Action:      action,
			Direction:   direction,
			Source:      source,
			Destination: destination,
			Before:      next,

			Protocol: data.Protocol,
		}
		return result, nil
	default:
		return nil, fmt.Errorf("impossible protocol type for acl rule %v: %T", current.Name, p)
	}
}

func aclRules(acl *ir.ACL) ([]vpcv1.NetworkACLRuleItemIntf, error) {
	rules := acl.Rules()
	ruleItems := make([]vpcv1.NetworkACLRuleItemIntf, len(rules))

	var next *vpcv1.NetworkACLRuleReference
	for i := len(ruleItems) - 1; i >= 0; i-- {
		name := utils.Ptr(fmt.Sprintf("rule%v", i))
		ref := allocateRef()
		current := &vpcv1.NetworkACLRuleReference{
			Name: name,
			Href: ref.Href,
			ID:   ref.ID,
		}
		rule, err := makeACLRuleItem(rules[i], current, next)
		if err != nil {
			return nil, err
		}
		ruleItems[i] = rule
		next = current
	}

	return ruleItems, nil
}

func updateACLList(model *configModel.ResourcesContainerModel, collection *ir.ACLCollection) error {
	if len(model.SubnetList) == 0 {
		return nil
	}
	subnet := model.SubnetList[0]
	vpcName := *subnet.VPC.Name
	aclName := ScopingString(vpcName, *subnet.Name)

	// decide if we are in a single-ACL mode
	if _, ok := collection.ACLs[vpcName][aclName]; ok {
		return updateACL(model, collection)
	}
	return updateSingleACL(model, collection)
}

func updateACL(model *configModel.ResourcesContainerModel, collection *ir.ACLCollection) error {
	for _, subnet := range model.SubnetList {
		vpcName := *subnet.VPC.Name
		aclName := ScopingString(vpcName, *subnet.Name)

		acl := collection.ACLs[vpcName][aclName]
		aclItem, err := newACLItem(subnet, acl) // create a new ACL item for the subnet
		if err != nil {
			return err
		}
		model.NetworkACLList = append(model.NetworkACLList, aclItem)
		subnet.NetworkACL = newNACLRef(aclItem)
	}
	globalIndex = 0 // making test results more predictable
	return nil
}

func updateSingleACL(model *configModel.ResourcesContainerModel, collection *ir.ACLCollection) error {
	aclItem := &configModel.NetworkACL{}
	var err error

	for i, subnet := range model.SubnetList {
		vpcName := *subnet.VPC.Name
		acl := collection.ACLs[vpcName][ScopingString(vpcName, "singleACL")]

		// if this is the first subnet being added to the ACL, add it to the list of network ACLs
		// otherwise, add the subnet reference to the existing ACL item
		if i == 0 {
			aclItem, err = newACLItem(subnet, acl)
			if err != nil {
				return err
			}
			model.NetworkACLList = append(model.NetworkACLList, aclItem)
		} else {
			aclItem.Subnets = append(aclItem.Subnets, *subnetRef(subnet))
		}
		subnet.NetworkACL = newNACLRef(aclItem)
	}
	globalIndex = 0 // making test results more predictable
	return nil
}

func newNACLRef(aclItem *configModel.NetworkACL) *vpcv1.NetworkACLReference {
	return &vpcv1.NetworkACLReference{
		ID:   aclItem.ID,
		CRN:  aclItem.CRN,
		Href: aclItem.Href,
		Name: utils.Ptr(*aclItem.Name),
	}
}

func newACLItem(subnet *configModel.Subnet, acl *ir.ACL) (*configModel.NetworkACL, error) {
	ref := allocateRef()
	rules, err := aclRules(acl)
	if err != nil {
		return nil, err
	}
	aclItem := configModel.NewNetworkACL(&vpcv1.NetworkACL{
		CRN:           ref.CRN,
		Href:          ref.Href,
		ID:            ref.ID,
		Name:          utils.Ptr(ir.ChangeScoping(acl.Name())),
		ResourceGroup: subnet.ResourceGroup,
		Rules:         rules,
		Subnets:       []vpcv1.SubnetReference{*subnetRef(subnet)},
		VPC:           subnet.VPC,
	})
	aclItem.Tags = []string{}
	return aclItem, nil
}

func subnetRef(subnet *configModel.Subnet) *vpcv1.SubnetReference {
	return &vpcv1.SubnetReference{
		Name:         subnet.Name,
		CRN:          subnet.CRN,
		Href:         subnet.Href,
		ID:           subnet.ID,
		ResourceType: subnet.ResourceType,
	}
}

func (w *Writer) WriteACL(collection *ir.ACLCollection, _ string) error {
	if err := updateACLList(w.model, collection); err != nil {
		return err
	}
	return w.writeModel()
}
