/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// updateSGs updates the config object file with the optimized SG rules
func updateSGs(model *configModel.ResourcesContainerModel, collection *ir.SGCollection) error {
	sgRefMap := parseSGRefMap(model)
	for _, sg := range model.SecurityGroupList {
		if sg.Name == nil || sg.VPC == nil || sg.VPC.Name == nil {
			continue
		}
		if err := updateSG(&sg.SecurityGroup, collection.SGs[*sg.VPC.Name][ir.SGName(*sg.Name)], sgRefMap); err != nil {
			return err
		}
	}
	return nil
}

func updateSG(sg *vpcv1.SecurityGroup, optimizedSG *ir.SG, sgRefMap map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference) error {
	optimizedRules := optimizedSG.AllRules()
	if len(optimizedRules) == len(sg.Rules) {
		return nil
	}
	sg.Rules = make([]vpcv1.SecurityGroupRuleIntf, len(optimizedRules))
	for i, rule := range optimizedRules {
		r, err := makeSGRuleItem(sgRefMap, rule, i)
		if err != nil {
			return err
		}
		sg.Rules[i] = r
	}
	return nil
}

func parseSGRefMap(model *configModel.ResourcesContainerModel) map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference {
	res := make(map[string]*vpcv1.SecurityGroupRuleRemoteSecurityGroupReference)
	for _, sg := range model.SecurityGroupList {
		res[*sg.Name] = &vpcv1.SecurityGroupRuleRemoteSecurityGroupReference{
			ID:   sg.ID,
			CRN:  sg.CRN,
			Href: sg.Href,
			Name: sg.Name,
		}
	}
	return res
}
