/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"log"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

func makeSGRuleItem(rule *ir.SGRule, i int, sgRemoteRef *vpcv1.SecurityGroupRuleRemoteSecurityGroupReference) vpcv1.SecurityGroupRuleIntf {
	iPVersion := utils.Ptr(ipv4Const)
	direction := direction(rule.Direction)
	ref := allocateRef()

	switch p := rule.Protocol.(type) {
	case ir.TCPUDP:
		data := tcpudp(p)
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp{
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Direction: direction,
			Remote:    sgRemoteRef,

			Protocol: data.Protocol,
			PortMin:  data.remotePortMin(rule.Direction),
			PortMax:  data.remotePortMax(rule.Direction),
		}
	case ir.ICMP:
		data := icmp(p)
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp{
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Direction: direction,
			Remote:    sgRemoteRef,

			Protocol: data.Protocol,
			Type:     data.Type,
			Code:     data.Code,
		}
	case ir.AnyProtocol:
		data := all()
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll{
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Direction: direction,
			Remote:    sgRemoteRef,

			Protocol: data.Protocol,
		}
	default:
		log.Fatalf("Impossible protocol type for sg rule %v: %T", i, p)
		return nil
	}
}

func makeSGItem(sg *ir.SG, sgRemoteRef *vpcv1.SecurityGroupRuleRemoteSecurityGroupReference) *configModel.SecurityGroup {
	ruleItems := make([]vpcv1.SecurityGroupRuleIntf, len(sg.Rules))
	for i := range sg.Rules {
		ruleItems[i] = makeSGRuleItem(&sg.Rules[i], i, sgRemoteRef)
	}

	result := configModel.NewSecurityGroup(&vpcv1.SecurityGroup{
		ID:      sgRemoteRef.ID,
		CRN:     sgRemoteRef.CRN,
		Href:    sgRemoteRef.Href,
		Name:    sgRemoteRef.Name,
		Rules:   ruleItems,
		Targets: []vpcv1.SecurityGroupTargetReferenceIntf{},
	})
	result.Tags = []string{}
	return result
}

func updateSG(model *configModel.ResourcesContainerModel, collection *ir.SGCollection) error {
	for i := range model.InstanceList {
		vpc := model.InstanceList[i].VPC
		sgName := ScopingString(*vpc.Name, *model.InstanceList[i].Name)
		sg := collection.SGs[*vpc.Name][ir.SGName(sgName)]
		ref := allocateRef()
		sgRemoteRef := &vpcv1.SecurityGroupRuleRemoteSecurityGroupReference{
			ID:   ref.ID,
			CRN:  ref.CRN,
			Href: ref.Href,
			Name: utils.Ptr(sgName),
		}
		sgItem := makeSGItem(sg, sgRemoteRef)
		sgItem.ResourceGroup = model.InstanceList[i].ResourceGroup
		sgItem.VPC = vpc
		sgRef := vpcv1.SecurityGroupReference{
			Name: sgRemoteRef.Name,
			Href: sgRemoteRef.Href,
			ID:   sgRemoteRef.ID,
			CRN:  sgRemoteRef.CRN,
		}
		model.SecurityGroupList = append(model.SecurityGroupList, sgItem)

		for j := range model.InstanceList[i].NetworkInterfaces {
			model.InstanceList[i].NetworkInterfaces[j].SecurityGroups = []vpcv1.SecurityGroupReference{sgRef}
		}
	}

	for i := range model.EndpointGWList {
		vpc := model.InstanceList[i].VPC
		sgName := ScopingString(*vpc.Name, *model.EndpointGWList[i].Name)
		sg := collection.SGs[*vpc.Name][ir.SGName(sgName)]
		ref := allocateRef()
		sgRemoteRef := &vpcv1.SecurityGroupRuleRemoteSecurityGroupReference{
			ID:   ref.ID,
			CRN:  ref.CRN,
			Href: ref.Href,
			Name: utils.Ptr(sgName),
		}
		sgItem := makeSGItem(sg, sgRemoteRef)
		sgItem.ResourceGroup = model.EndpointGWList[i].ResourceGroup
		sgItem.VPC = vpc
		sgRef := vpcv1.SecurityGroupReference{
			Name: sgRemoteRef.Name,
			Href: sgRemoteRef.Href,
			ID:   sgRemoteRef.ID,
			CRN:  sgRemoteRef.CRN,
		}
		model.SecurityGroupList = append(model.SecurityGroupList, sgItem)

		model.EndpointGWList[i].SecurityGroups = []vpcv1.SecurityGroupReference{sgRef}
	}

	return nil
}

func (w *Writer) WriteSG(collection *ir.SGCollection, _ string) error {
	if err := updateSG(w.model, collection); err != nil {
		return err
	}
	return w.writeModel()
}
