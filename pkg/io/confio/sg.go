/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"log"
	"slices"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const ResourceTypeNif = "network_interface"
const ResourceTypeEndpointGateway = "endpoint_gateway"

func findAndDeleteTargetFromSG(model *configModel.ResourcesContainerModel, sgIndex int, id *string) {
	sg := model.SecurityGroupList[sgIndex]
	for i := range sg.Targets {
		if target, ok := sg.Targets[i].(*vpcv1.SecurityGroupTargetReference); ok && *target.ID == *id {
			sg.Targets = slices.Delete(sg.Targets, i, i+1) // deleteTargetFromSG
			break
		}
	}
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
	nameToSGRemoteRef[name] = sgRemoteRef
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
	remote := sgRemote(nameToSGRemoteRef, rule)

	switch p := rule.Protocol.(type) {
	case ir.TCPUDP:
		data := tcpudp(p)
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp{
			Direction: direction,
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Local:     local,
			Remote:    remote,
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
			Remote:    remote,
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
			Remote:    remote,
			Protocol:  data.Protocol,
		}
	default:
		log.Fatalf("Impossible protocol type for sg rule %v: %T", i, p)
		return nil
	}
}

func makeSGRules(nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference,
	sg *ir.SG) []vpcv1.SecurityGroupRuleIntf {
	rules := sg.AllRules()
	ruleItems := make([]vpcv1.SecurityGroupRuleIntf, len(rules))
	for i := range rules {
		ruleItems[i] = makeSGRuleItem(nameToSGRemoteRef, &rules[i], i)
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
			ResourceType: utils.Ptr(ResourceTypeNif),
		}
		targets[i] = sgTargetRef
	}

	return targets
}

func updateSGInstances(model *configModel.ResourcesContainerModel, collection *ir.SGCollection,
	nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference, idToSGIndex map[string]int) {
	for _, instance := range model.InstanceList {
		vpc := instance.VPC
		sgName := ScopingString(*vpc.Name, *instance.Name)
		sgItemName := utils.Ptr(ir.ChangeScoping(sgName))
		sg := collection.SGs[*vpc.Name][ir.SGName(sgName)]
		ref := lookupOrCreate(nameToSGRemoteRef, *sgItemName)

		sgItem := configModel.NewSecurityGroup(&vpcv1.SecurityGroup{
			CRN:           ref.CRN,
			Href:          ref.Href,
			ID:            ref.ID,
			Name:          sgItemName,
			ResourceGroup: instance.ResourceGroup,
			Rules:         makeSGRules(nameToSGRemoteRef, sg),
			Targets:       parseTargetsSGInstance(instance),
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

		for j := range instance.NetworkInterfaces {
			for k := range instance.NetworkInterfaces[j].SecurityGroups {
				sgID := instance.NetworkInterfaces[j].SecurityGroups[k].ID
				nifID := instance.NetworkInterfaces[j].ID
				findAndDeleteTargetFromSG(model, idToSGIndex[*sgID], nifID)
			}
			instance.NetworkInterfaces[j].SecurityGroups = []vpcv1.SecurityGroupReference{sgRef}
		}
	}
}

func updateSGEndpointGW(model *configModel.ResourcesContainerModel, collection *ir.SGCollection,
	nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference, idToSGIndex map[string]int) {
	for _, endpointGW := range model.EndpointGWList {
		vpc := endpointGW.VPC
		sgName := ScopingString(*vpc.Name, *endpointGW.Name)
		sgItemName := utils.Ptr(ir.ChangeScoping(sgName))
		sg := collection.SGs[*vpc.Name][ir.SGName(sgName)]
		ref := lookupOrCreate(nameToSGRemoteRef, *sgItemName)
		target := &vpcv1.SecurityGroupTargetReference{
			Name:         endpointGW.Name,
			Href:         endpointGW.Href,
			ID:           endpointGW.ID,
			CRN:          endpointGW.CRN,
			ResourceType: utils.Ptr(ResourceTypeEndpointGateway),
		}

		sgItem := configModel.NewSecurityGroup(&vpcv1.SecurityGroup{
			CRN:           ref.CRN,
			Href:          ref.Href,
			ID:            ref.ID,
			Name:          sgItemName,
			ResourceGroup: endpointGW.ResourceGroup,
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

		for j := range endpointGW.SecurityGroups {
			sgID := endpointGW.SecurityGroups[j].ID
			endpointGatewayID := endpointGW.ID
			findAndDeleteTargetFromSG(model, idToSGIndex[*sgID], endpointGatewayID)
		}
		endpointGW.SecurityGroups = []vpcv1.SecurityGroupReference{sgRef}
	}
}

func updateSG(model *configModel.ResourcesContainerModel, collection *ir.SGCollection) {
	nameToSGRemoteRef := make(map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference)
	idToSGIndex := make(map[string]int, len(model.SecurityGroupList))
	for i := range model.SecurityGroupList {
		idToSGIndex[*model.SecurityGroupList[i].ID] = i
	}

	updateSGInstances(model, collection, nameToSGRemoteRef, idToSGIndex)
	updateSGEndpointGW(model, collection, nameToSGRemoteRef, idToSGIndex)

	globalIndex = 0 // making test results more predictable
}

func (w *Writer) WriteSG(collection *ir.SGCollection, _ string) error {
	updateSG(w.model, collection)
	return w.writeModel()
}
