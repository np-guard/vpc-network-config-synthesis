/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"fmt"
	"log"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

func cidr(address *ipblock.IPBlock) *string {
	return utils.Ptr(address.ToCidrListString())
}

func makeACLRuleItem(rule *ir.ACLRule, current,
	next *vpcv1.NetworkACLRuleReference) vpcv1.NetworkACLRuleItemIntf {
	iPVersion := utils.Ptr(ipv4Const)
	direction := direction(rule.Direction)
	action := action(rule.Action)
	source := cidr(rule.Source)
	destination := cidr(rule.Destination)
	switch p := rule.Protocol.(type) {
	case ir.TCPUDP:
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
			SourcePortMin:      data.SourcePortMin,
			SourcePortMax:      data.SourcePortMax,
			DestinationPortMin: data.DestinationPortMin,
			DestinationPortMax: data.DestinationPortMax,
		}
		return result
	case ir.ICMP:
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
		return result
	case ir.AnyProtocol:
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
		return result
	default:
		log.Fatalf("Impossible protocol type for acl rule %v: %T", current.Name, p)
		return nil
	}
}

func aclRules(acl *ir.ACL) []vpcv1.NetworkACLRuleItemIntf {
	ruleItems := make([]vpcv1.NetworkACLRuleItemIntf, len(acl.Rules()))
	rules := acl.Rules()

	var next *vpcv1.NetworkACLRuleReference
	for i := len(ruleItems) - 1; i >= 0; i-- {
		name := utils.Ptr(fmt.Sprintf("rule%v", i))
		ref := allocateRef()
		current := &vpcv1.NetworkACLRuleReference{
			Name: name,
			Href: ref.Href,
			ID:   ref.ID,
		}
		ruleItems[i] = makeACLRuleItem(&rules[i], current, next)
		next = current
	}

	return ruleItems
}

func updateACL(model *configModel.ResourcesContainerModel, collection *ir.ACLCollection) error {
	for i := range model.SubnetList {
		subnet := model.SubnetList[i]
		vpc := subnet.VPC
		aclName := ScopingString(*vpc.Name, *subnet.Name)
		acl := collection.ACLs[*vpc.Name][aclName]
		subnetRef := &vpcv1.SubnetReference{
			Name:         subnet.Name,
			CRN:          subnet.CRN,
			Href:         subnet.Href,
			ID:           subnet.ID,
			ResourceType: subnet.ResourceType,
		}

		ref := allocateRef()
		aclItem := configModel.NewNetworkACL(&vpcv1.NetworkACL{
			CRN:           ref.CRN,
			Href:          ref.Href,
			ID:            ref.ID,
			Name:          utils.Ptr(ir.ChangeScoping(acl.Name())),
			ResourceGroup: subnet.ResourceGroup,
			Rules:         aclRules(acl),
			Subnets:       []vpcv1.SubnetReference{*subnetRef},
			VPC:           subnet.VPC,
		})
		aclItem.Tags = []string{}

		model.NetworkACLList = append(model.NetworkACLList, aclItem)
		subnet.NetworkACL = &vpcv1.NetworkACLReference{
			ID:   aclItem.ID,
			CRN:  aclItem.CRN,
			Href: aclItem.Href,
			Name: utils.Ptr(*aclItem.Name),
		}
	}
	GlobalIndex = 0 // making test results more predictable
	return nil
}

func (w *Writer) WriteACL(collection *ir.ACLCollection, _ string) error {
	if err := updateACL(w.model, collection); err != nil {
		return err
	}
	return w.writeModel()
}
