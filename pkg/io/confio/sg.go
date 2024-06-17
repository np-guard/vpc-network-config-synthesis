/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"log"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const ResourceTypeSGTarget = "network_interface"

func deleteTargetFromSG(model *configModel.ResourcesContainerModel, sgIndex int, targetID *string) {
	// for i := range model.SecurityGroupList[sgIndex].Targets {
	// 	switch target := model.SecurityGroupList[sgIndex].Targets[i].(type) {
	// 	case *vpcv1.SecurityGroupTargetReferenceEndpointGatewayReference:
	// 		if targetID == target.ID {

	// 		}
	// 	}
	// 	case
	// }
}

func lookupOrCreate(nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference,
	name string) *vpcv1.SecurityGroupRuleRemoteSecurityGroupReference {
	if sgRemoteRef, ok := nameToSGRemoteRef[name]; ok {
		return sgRemoteRef
	}
	ref := allocateRef()
	sgRemoteRef := &vpcv1.SecurityGroupRuleRemoteSecurityGroupReference{
		ID:   ref.ID,
		CRN:  ref.CRN,
		Href: ref.Href,
		Name: utils.Ptr(name),
	}
	return sgRemoteRef
}

func sgRemote(nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference,
	rule *ir.SGRule) vpcv1.SecurityGroupRuleRemoteIntf {
	st := rule.Remote.String()
	switch t := rule.Remote.(type) {
	case *ipblock.IPBlock:
		if ir.IsIPAddress(t) {
			return &vpcv1.SecurityGroupRuleRemoteIP{
				Address: &st,
			}
		}
		return &vpcv1.SecurityGroupRuleRemoteCIDR{
			CIDRBlock: &st,
		}
	case ir.SGName:
		return lookupOrCreate(nameToSGRemoteRef, ir.ChangeScoping(st))
	}
	return nil
}

func makeSGRuleItem(nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference,
	rule *ir.SGRule, i int) vpcv1.SecurityGroupRuleIntf {
	iPVersion := utils.Ptr(ipv4Const)
	direction := direction(rule.Direction)
	cidrAll := ipblock.CidrAll
	local := &vpcv1.SecurityGroupRuleLocal{
		CIDRBlock: &cidrAll,
	}
	ref := allocateRef()

	switch p := rule.Protocol.(type) {
	case ir.TCPUDP:
		data := tcpudp(p)
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp{
			Direction: direction,
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Local:     local,
			Remote:    sgRemote(nameToSGRemoteRef, rule),
			Protocol:  data.Protocol,
			PortMin:   data.remotePortMin(rule.Direction),
			PortMax:   data.remotePortMax(rule.Direction),
		}
	case ir.ICMP:
		data := icmp(p)
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp{
			Direction: direction,
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Local:     local,
			Remote:    sgRemote(nameToSGRemoteRef, rule),
			Protocol:  data.Protocol,
			Type:      data.Type,
			Code:      data.Code,
		}
	case ir.AnyProtocol:
		data := all()
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll{
			Direction: direction,
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Local:     local,
			Remote:    sgRemote(nameToSGRemoteRef, rule),
			Protocol:  data.Protocol,
		}
	default:
		log.Fatalf("Impossible protocol type for sg rule %v: %T", i, p)
		return nil
	}
}

func makeSGRules(nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference,
	sg *ir.SG) []vpcv1.SecurityGroupRuleIntf {
	ruleItems := make([]vpcv1.SecurityGroupRuleIntf, len(sg.Rules))
	for i := range sg.Rules {
		ruleItems[i] = makeSGRuleItem(nameToSGRemoteRef, &sg.Rules[i], i)
	}
	return ruleItems
}

func parseTargetsSGInstance(instance *configModel.Instance) []vpcv1.SecurityGroupTargetReferenceIntf {
	targets := make([]vpcv1.SecurityGroupTargetReferenceIntf, len(instance.NetworkInterfaces))
	for i := range instance.NetworkInterfaces {
		sgTargetRef := &vpcv1.SecurityGroupTargetReference{
			Name:         instance.NetworkInterfaces[i].Name,
			Href:         instance.NetworkInterfaces[i].Href,
			ID:           instance.NetworkInterfaces[i].ID,
			ResourceType: utils.Ptr(ResourceTypeSGTarget),
		}
		targets[i] = sgTargetRef
	}

	return targets
}

func updateSGInstances(model *configModel.ResourcesContainerModel, collection *ir.SGCollection,
	nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference, idToSGIndex map[string]int) {
	for i := range model.InstanceList {
		vpc := model.InstanceList[i].VPC
		sgName := ScopingString(*vpc.Name, *model.InstanceList[i].Name)
		sgItemName := utils.Ptr(ir.ChangeScoping(sgName))
		sg := collection.SGs[*vpc.Name][ir.SGName(sgName)]
		ref := lookupOrCreate(nameToSGRemoteRef, *sgItemName)

		sgItem := configModel.NewSecurityGroup(&vpcv1.SecurityGroup{
			CRN:           ref.CRN,
			Href:          ref.Href,
			ID:            ref.ID,
			Name:          sgItemName,
			ResourceGroup: model.InstanceList[i].ResourceGroup,
			Rules:         makeSGRules(nameToSGRemoteRef, sg),
			Targets:       parseTargetsSGInstance(model.InstanceList[i]),
			VPC:           vpc,
		})
		sgItem.Tags = []string{}
		model.SecurityGroupList = append(model.SecurityGroupList, sgItem)

		sgRef := vpcv1.SecurityGroupReference{
			CRN:  ref.CRN,
			Href: ref.Href,
			ID:   ref.ID,
			Name: sgItemName,
		}

		for j := range model.InstanceList[i].NetworkInterfaces {
			for k := range model.InstanceList[i].NetworkInterfaces[j].SecurityGroups {
				sgID := model.InstanceList[i].NetworkInterfaces[j].SecurityGroups[k].ID
				nifID := model.InstanceList[i].NetworkInterfaces[j].ID
				if _, ok := idToSGIndex[*sgID]; !ok {
					log.Fatalf("ERROR")
				}
				deleteTargetFromSG(model, idToSGIndex[*sgID], nifID)
			}
			model.InstanceList[i].NetworkInterfaces[j].SecurityGroups = []vpcv1.SecurityGroupReference{sgRef}
		}
	}
}

func updateSGEndpointGW(model *configModel.ResourcesContainerModel, collection *ir.SGCollection,
	nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference, idToSGIndex map[string]int) {
	for i := range model.EndpointGWList {
		vpc := model.EndpointGWList[i].VPC
		sgName := ScopingString(*vpc.Name, *model.EndpointGWList[i].Name)
		sgItemName := utils.Ptr(ir.ChangeScoping(sgName))
		sg := collection.SGs[*vpc.Name][ir.SGName(sgName)]
		ref := lookupOrCreate(nameToSGRemoteRef, *sgItemName)
		target := &vpcv1.SecurityGroupTargetReference{
			Name:         model.EndpointGWList[i].Name,
			Href:         model.EndpointGWList[i].Href,
			ID:           model.EndpointGWList[i].ID,
			ResourceType: utils.Ptr(ResourceTypeSGTarget),
		}

		sgItem := configModel.NewSecurityGroup(&vpcv1.SecurityGroup{
			CRN:           ref.CRN,
			Href:          ref.Href,
			ID:            ref.ID,
			Name:          sgItemName,
			ResourceGroup: model.EndpointGWList[i].ResourceGroup,
			Rules:         makeSGRules(nameToSGRemoteRef, sg),
			Targets:       []vpcv1.SecurityGroupTargetReferenceIntf{target},
			VPC:           vpc,
		})
		sgItem.Tags = []string{}
		model.SecurityGroupList = append(model.SecurityGroupList, sgItem)

		sgRef := vpcv1.SecurityGroupReference{
			CRN:  ref.CRN,
			Href: ref.Href,
			ID:   ref.ID,
			Name: sgItemName,
		}

		for j := range model.EndpointGWList[i].SecurityGroups {
			sgID := model.EndpointGWList[i].SecurityGroups[j].ID
			endpointGatewayID := model.EndpointGWList[i].ID
			deleteTargetFromSG(model, idToSGIndex[*sgID], endpointGatewayID)
		}
		model.EndpointGWList[i].SecurityGroups = []vpcv1.SecurityGroupReference{sgRef}
	}
}

func updateSG(model *configModel.ResourcesContainerModel, collection *ir.SGCollection) error {
	// model.SecurityGroupList = make([]*configModel.SecurityGroup, 0) // delete old SGs
	nameToSGRemoteRef := make(map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference)
	idToSGIndex := make(map[string]int, len(model.SecurityGroupList))
	for i := range model.SecurityGroupList {
		idToSGIndex[*model.SecurityGroupList[i].ID] = i
	}

	updateSGInstances(model, collection, nameToSGRemoteRef, idToSGIndex)
	updateSGEndpointGW(model, collection, nameToSGRemoteRef, idToSGIndex)

	GlobalIndex = 0 // for tests
	return nil
}

func (w *Writer) WriteSG(collection *ir.SGCollection, _ string) error {
	if err := updateSG(w.model, collection); err != nil {
		return err
	}
	return w.writeModel()
}
