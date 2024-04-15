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
	vpc := model.SubnetList[0].VPC
	resourceGroup := model.SubnetList[0].ResourceGroup
	for _, sgName := range collection.SortedSGNames() {
		sg := collection.SGs[sgName]
		if sg == nil {
			continue
		}
		ref := allocateRef()
		sgRemoteRef := &vpcv1.SecurityGroupRuleRemoteSecurityGroupReference{
			ID:   ref.ID,
			CRN:  ref.CRN,
			Href: ref.Href,
			Name: utils.Ptr(sgName.String()),
		}
		sgItem := makeSGItem(sg, sgRemoteRef)
		sgItem.ResourceGroup = resourceGroup
		sgItem.VPC = vpc

		sgRef := vpcv1.SecurityGroupReference{
			Name: sgRemoteRef.Name,
			Href: sgRemoteRef.Href,
			ID:   sgRemoteRef.ID,
			CRN:  sgRemoteRef.CRN,
		}
		model.SecurityGroupList = append(model.SecurityGroupList, sgItem)

		for _, attached := range sg.Attached {
			for _, instance := range model.InstanceList {
				if *instance.Name == string(attached) {
					for j := range instance.NetworkInterfaces {
						instance.NetworkInterfaces[j].SecurityGroups = []vpcv1.SecurityGroupReference{sgRef}
					}
					break
				}
			}
		}
	}
	return nil
}

func (w *Writer) WriteSG(collection *ir.SGCollection) error {
	if err := updateSG(w.model, collection); err != nil {
		return err
	}
	return w.writeModel()
}
