/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"errors"
	"fmt"
	"slices"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

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
	case *netset.IPBlock:
		if t.IsSingleIPAddress() { // single IP address
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
	rule *ir.SGRule, i int) (vpcv1.SecurityGroupRuleIntf, error) {
	iPVersion := utils.Ptr(ipv4Const)
	direction := direction(rule.Direction)
	ref := allocateRef()
	remote := sgRemote(nameToSGRemoteRef, rule)

	var local vpcv1.SecurityGroupRuleLocalIntf
	if rule.Local.IsSingleIPAddress() {
		local = &vpcv1.SecurityGroupRuleLocalIP{
			Address: utils.Ptr(rule.Local.FirstIPAddress()),
		}
	} else {
		local = &vpcv1.SecurityGroupRuleLocalCIDR{
			CIDRBlock: utils.Ptr(rule.Local.ToCidrList()[0]),
		}
	}

	switch p := rule.Protocol.(type) {
	case netp.TCPUDP:
		data := tcpudp(p)
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp{
			Direction: direction,
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Local:     local,
			Remote:    remote,
			Protocol:  data.Protocol,
			PortMin:   data.dstPortMin,
			PortMax:   data.dstPortMax,
		}, nil
	case netp.ICMP:
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
		}, nil
	case netp.AnyProtocol:
		data := all()
		return &vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll{
			Direction: direction,
			Href:      ref.Href,
			ID:        ref.ID,
			IPVersion: iPVersion,
			Local:     local,
			Remote:    remote,
			Protocol:  data.Protocol,
		}, nil
	}
	return nil, fmt.Errorf("impossible protocol type for sg rule %v: %T", i, rule.Protocol)
}

func makeSGRules(nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference,
	sg *ir.SG) ([]vpcv1.SecurityGroupRuleIntf, error) {
	rules := sg.AllRules()
	ruleItems := make([]vpcv1.SecurityGroupRuleIntf, len(rules))
	for i, rule := range rules {
		rule, err := makeSGRuleItem(nameToSGRemoteRef, rule, i)
		if err != nil {
			return nil, err
		}
		ruleItems[i] = rule
	}
	return ruleItems, nil
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
	nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference, idToSGIndex map[string]int) error {
	for _, instance := range model.InstanceList {
		vpc := instance.VPC
		sgName := ScopingString(*vpc.Name, *instance.Name)
		sgItemName := utils.Ptr(ir.ChangeScoping(sgName))
		sg := collection.SGs[*vpc.Name][ir.SGName(sgName)]
		sgRules, err := makeSGRules(nameToSGRemoteRef, sg)
		if err != nil {
			return err
		}
		ref := lookupOrCreate(nameToSGRemoteRef, *sgItemName)

		sgItem := configModel.NewSecurityGroup(&vpcv1.SecurityGroup{
			CRN:           ref.CRN,
			Href:          ref.Href,
			ID:            ref.ID,
			Name:          sgItemName,
			ResourceGroup: instance.ResourceGroup,
			Rules:         sgRules,
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
	return nil
}

func updateSGEndpointGW(model *configModel.ResourcesContainerModel, collection *ir.SGCollection,
	nameToSGRemoteRef map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference, idToSGIndex map[string]int) error {
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

		sgRules, err := makeSGRules(nameToSGRemoteRef, sg)
		if err != nil {
			return err
		}
		sgItem := configModel.NewSecurityGroup(&vpcv1.SecurityGroup{
			CRN:           ref.CRN,
			Href:          ref.Href,
			ID:            ref.ID,
			Name:          sgItemName,
			ResourceGroup: endpointGW.ResourceGroup,
			Rules:         sgRules,
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
	return nil
}

func writeSGs(model *configModel.ResourcesContainerModel, collection *ir.SGCollection) error {
	nameToSGRemoteRef := make(map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference)
	idToSGIndex := make(map[string]int, len(model.SecurityGroupList))
	for i := range model.SecurityGroupList {
		idToSGIndex[*model.SecurityGroupList[i].ID] = i
	}

	err1 := updateSGInstances(model, collection, nameToSGRemoteRef, idToSGIndex)
	err2 := updateSGEndpointGW(model, collection, nameToSGRemoteRef, idToSGIndex)
	return errors.Join(err1, err2)
}

func (w *Writer) WriteSG(collection *ir.SGCollection, _ string, isSynth bool) error {
	var err error
	if isSynth {
		err = writeSGs(w.model, collection)
	} else {
		err = updateSGs(w.model, collection)
	}
	if err != nil {
		return err
	}
	globalIndex = 0 // making test results more predictable
	return w.writeModel()
}
