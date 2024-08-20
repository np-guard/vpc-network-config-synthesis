/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"fmt"
	"reflect"

	vpc1 "github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

// ReadSG translates SGs from a config_object file to map[ir.SGName]*SG
func ReadSGs(filename string) (map[ir.SGName]*ir.SG, error) {
	config, err := readModel(filename)
	if err != nil {
		return nil, err
	}
	result := make(map[ir.SGName]*ir.SG, len(config.SecurityGroupList))

	for _, sg := range config.SecurityGroupList {
		inbound, outbound, err := translateSGRules(&sg.SecurityGroup)
		if err != nil {
			return nil, err
		}
		if sg.Name != nil {
			result[ir.SGName(*sg.Name)] = &ir.SG{InboundRules: inbound, OutboundRules: outbound}
		}
	}
	return result, nil
}

// parse security rules, splitted into ingress and egress rules
func translateSGRules(sg *vpc1.SecurityGroup) (ingressRules, egressRules []ir.SGRule, err error) {
	ingressRules = []ir.SGRule{}
	egressRules = []ir.SGRule{}
	for index := range sg.Rules {
		rule, err := translateSGRule(sg, index)
		if err != nil {
			return nil, nil, err
		}
		if rule.Direction == ir.Inbound {
			ingressRules = append(ingressRules, rule)
		} else {
			egressRules = append(egressRules, rule)
		}
	}
	return ingressRules, egressRules, nil
}

// translateSGRule translates a security group rule to ir.SGRule
func translateSGRule(sg *vpc1.SecurityGroup, index int) (sgRule ir.SGRule, err error) {
	fmt.Printf("Type of sgRule[%d]: %s \n", index, reflect.TypeOf(sg.Rules[index]))
	switch r := sg.Rules[index].(type) {
	case *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolAll:
		return translateSGRuleProtocolAll(r)
	case *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp:
		return translateSGRuleProtocolTCPUDP(r)
	case *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp:
		return translateSGRuleProtocolIcmp(r)
	}
	return ir.SGRule{}, fmt.Errorf("error parsing rule number %d in %s sg", index, *sg.Name)
}

func translateSGRuleProtocolAll(rule *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolAll) (sgRule ir.SGRule, err error) {
	direction, err := translateDirection(*rule.Direction)
	if err != nil {
		return ir.SGRule{}, err
	}
	remote, err := translateRemote(rule.Remote)
	if err != nil {
		return ir.SGRule{}, err
	}
	local, err := translateLocal(rule.Local)
	if err != nil {
		return ir.SGRule{}, err
	}
	return ir.SGRule{Direction: direction, Remote: remote, Protocol: netp.AnyProtocol{}, Local: local}, nil
}

func translateSGRuleProtocolTCPUDP(rule *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp) (sgRule ir.SGRule, err error) {
	direction, err := translateDirection(*rule.Direction)
	if err != nil {
		return ir.SGRule{}, err
	}
	remote, err := translateRemote(rule.Remote)
	if err != nil {
		return ir.SGRule{}, err
	}
	local, err := translateLocal(rule.Local)
	if err != nil {
		return ir.SGRule{}, err
	}
	protocol, err := translateProtocolTCPUDP(rule)
	if err != nil {
		return ir.SGRule{}, err
	}
	return ir.SGRule{Direction: direction, Remote: remote, Protocol: protocol, Local: local}, nil
}

func translateSGRuleProtocolIcmp(rule *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp) (sgRule ir.SGRule, err error) {
	direction, err := translateDirection(*rule.Direction)
	if err != nil {
		return ir.SGRule{}, err
	}
	remote, err := translateRemote(rule.Remote)
	if err != nil {
		return ir.SGRule{}, err
	}
	local, err := translateLocal(rule.Local)
	if err != nil {
		return ir.SGRule{}, err
	}
	protocol, err := netp.ICMPFromTypeAndCode64(rule.Type, rule.Code)
	if err != nil {
		return ir.SGRule{}, err
	}
	return ir.SGRule{Direction: direction, Remote: remote, Protocol: protocol, Local: local}, nil
}

func translateDirection(direction string) (ir.Direction, error) {
	if direction == "inbound" {
		return ir.Inbound, nil
	} else if direction == "outbound" {
		return ir.Outbound, nil
	}
	return ir.Inbound, fmt.Errorf("a SG rule direction must be either inbound or outbound")
}

func translateRemote(remote vpc1.SecurityGroupRuleRemoteIntf) (ir.RemoteType, error) {
	if r, ok := remote.(*vpc1.SecurityGroupRuleRemote); ok {
		switch {
		case r.CIDRBlock != nil:
			return netset.IPBlockFromCidr(*r.CIDRBlock)
		case r.Address != nil:
			return netset.IPBlockFromIPAddress(*r.Address)
		case r.Name != nil:
			return ir.SGName(*r.Name), nil
		}
	}
	return nil, fmt.Errorf("unexpected SG rule remote")
}

func translateLocal(local vpc1.SecurityGroupRuleLocalIntf) (*netset.IPBlock, error) {
	if l, ok := local.(*vpc1.SecurityGroupRuleLocal); ok {
		if l.CIDRBlock != nil {
			return netset.IPBlockFromCidr(*l.CIDRBlock)
		}
		if l.Address != nil {
			return netset.IPBlockFromIPAddress(*l.CIDRBlock)
		}
	}
	return nil, fmt.Errorf("error parsing Local field")
}

func translateProtocolTCPUDP(rule *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp) (netp.Protocol, error) {
	isTCP := *rule.Protocol == string(netp.ProtocolStringTCP)
	minDstPort := utils.GetProperty(rule.PortMin, netp.MinPort)
	maxDstPort := utils.GetProperty(rule.PortMax, netp.MaxPort)
	return netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(minDstPort), int(maxDstPort))
}
