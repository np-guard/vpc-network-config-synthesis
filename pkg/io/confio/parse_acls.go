/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"errors"
	"fmt"
	"log"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// ReadACLs translates ACLs from a config_object file to ir.ACLCollection
func ReadACLs(filename string) (*ir.ACLCollection, error) {
	config, err := readModel(filename)
	if err != nil {
		return nil, err
	}

	result := ir.NewACLCollection()
	for i, acl := range config.NetworkACLList {
		if acl.Name == nil || acl.VPC == nil || acl.VPC.Name == nil {
			log.Printf("Warning: missing acl/VPC name in acl at index %d\n", i)
			continue
		}
		inbound, outbound, err := translateACLRules(&acl.NetworkACL)
		if err != nil {
			return nil, err
		}
		vpcName := *acl.VPC.Name
		if result.ACLs[vpcName] == nil {
			result.ACLs[vpcName] = make(map[string]*ir.ACL)
		}
		result.ACLs[vpcName][*acl.Name] = &ir.ACL{Name: *acl.Name,
			Subnets:  parseAttachedSubnets(&acl.NetworkACL),
			Inbound:  inbound,
			Outbound: outbound,
		}
	}
	return result, nil
}

func translateACLRules(acl *vpcv1.NetworkACL) (inbound, outbound []*ir.ACLRule, err error) {
	inbound = make([]*ir.ACLRule, 0)
	outbound = make([]*ir.ACLRule, 0)
	for index := range acl.Rules {
		rule, err := translateACLRule(acl, index)
		if err != nil {
			return nil, nil, err
		}
		if rule.Direction == ir.Inbound {
			inbound = append(inbound, rule)
		} else {
			outbound = append(outbound, rule)
		}
	}
	return inbound, outbound, nil
}

func translateACLRule(acl *vpcv1.NetworkACL, i int) (*ir.ACLRule, error) {
	switch r := acl.Rules[i].(type) {
	case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll:
		return translateACLRuleProtocolAll(r)
	case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp:
		return translateACLRuleProtocolTCPUDP(r)
	case *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp:
		return translateACLRuleProtocolIcmp(r)
	}
	return nil, fmt.Errorf("error parsing rule number %d in acl %s in VPC %s", i, *acl.Name, *acl.VPC.Name)
}

func translateACLRuleProtocolAll(rule *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll) (*ir.ACLRule, error) {
	action, err1 := translateAction(rule.Action)
	direction, err2 := translateDirection(*rule.Direction)
	src, err3 := translateResource(rule.Source)
	dst, err4 := translateResource(rule.Destination)
	if err := errors.Join(err1, err2, err3, err4); err != nil {
		return nil, err
	}
	return &ir.ACLRule{
		Action:      action,
		Direction:   direction,
		Source:      src,
		Destination: dst,
		Protocol:    netp.AnyProtocol{},
	}, nil
}

func translateACLRuleProtocolTCPUDP(rule *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp) (*ir.ACLRule, error) {
	action, err1 := translateAction(rule.Action)
	direction, err2 := translateDirection(*rule.Direction)
	src, err3 := translateResource(rule.Source)
	dst, err4 := translateResource(rule.Destination)
	protocol, err5 := translateProtocolTCPUDP(*rule.Protocol, rule.DestinationPortMin, rule.DestinationPortMax)
	if err := errors.Join(err1, err2, err3, err4, err5); err != nil {
		return nil, err
	}

	return &ir.ACLRule{
		Action:      action,
		Direction:   direction,
		Source:      src,
		Destination: dst,
		Protocol:    protocol,
	}, nil
}

func translateACLRuleProtocolIcmp(rule *vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp) (*ir.ACLRule, error) {
	action, err1 := translateAction(rule.Action)
	direction, err2 := translateDirection(*rule.Direction)
	src, err3 := translateResource(rule.Source)
	dst, err4 := translateResource(rule.Destination)
	protocol, err5 := netp.ICMPFromTypeAndCode64WithoutRFCValidation(rule.Type, rule.Code)
	if err := errors.Join(err1, err2, err3, err4, err5); err != nil {
		return nil, err
	}

	return &ir.ACLRule{
		Action:      action,
		Direction:   direction,
		Source:      src,
		Destination: dst,
		Protocol:    protocol,
	}, nil
}

func translateAction(action *string) (ir.Action, error) {
	if *action == vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllActionAllowConst {
		return ir.Allow, nil
	} else if *action == vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllActionDenyConst {
		return ir.Deny, nil
	}
	return ir.Deny, fmt.Errorf("an nACL rule action must be either allow or deny")
}

func translateResource(ipAddrs *string) (*netset.IPBlock, error) {
	return netset.IPBlockFromCidrOrAddress(*ipAddrs)
}

func parseAttachedSubnets(acl *vpcv1.NetworkACL) []string {
	if len(acl.Subnets) == 0 {
		log.Printf("Warning: nACL %s does not have attached subnets", *acl.Name)
	}
	res := make([]string, 0)
	for i, subnet := range acl.Subnets {
		if subnet.Name != nil {
			res = append(res, *subnet.Name)
		} else {
			log.Printf("Warning: error translating subnet %d in %s nACL", i, *acl.Name)
		}
	}
	return res
}
