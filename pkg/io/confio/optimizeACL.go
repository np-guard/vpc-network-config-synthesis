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

func updateACLs(model *configModel.ResourcesContainerModel, collection *ir.ACLCollection) error {
	for _, acl := range model.NetworkACLList {
		if acl.Name == nil || acl.VPC == nil || acl.VPC.Name == nil {
			continue
		}
		if err := updateACL(&acl.NetworkACL, collection.ACLs[*acl.VPC.Name][*acl.VPC.Name]); err != nil {
			return err
		}
	}
	return nil
}

func updateACL(acl *vpcv1.NetworkACL, optimizedACL *ir.ACL) error {
	optimizedRules := optimizedACL.Rules()
	if len(optimizedRules) >= len(acl.Rules) {
		return nil
	}
	rules, err := aclRules(optimizedACL)
	acl.Rules = rules
	return err
}
