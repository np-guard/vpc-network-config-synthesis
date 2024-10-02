/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"fmt"
	"log"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

// ReadSG translates SGs from a config_object file to map[ir.SGName]*SG
func ReadSGs(filename string) (*ir.SGCollection, error) {
	config, err := readModel(filename)
	if err != nil {
		return nil, err
	}

	result := ir.NewSGCollection()
	for i, sg := range config.SecurityGroupList {
		inbound, outbound, err := translateSGRules(&sg.SecurityGroup)
		if err != nil {
			return nil, err
		}
		if sg.Name == nil || sg.VPC == nil || sg.VPC.Name == nil {
			log.Printf("Warning: missing SG/VPC name in sg at index %d\n", i)
			continue
		}
		vpcName := *sg.VPC.Name
		if result.SGs[vpcName] == nil {
			result.SGs[vpcName] = make(map[ir.SGName]*ir.SG)
		}
		result.SGs[vpcName][ir.SGName(*sg.Name)] = &ir.SG{InboundRules: inbound, OutboundRules: outbound}
	}
	return result, nil
}

// parse security rules, splitted into ingress and egress rules
func translateSGRules(sg *vpcv1.SecurityGroup) (ingressRules, egressRules []*ir.SGRule, err error) {
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
func translateSGRule(sg *vpcv1.SecurityGroup, index int) (sgRule *ir.SGRule, err error) {
	switch r := sg.Rules[index].(type) {
	case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll:
		return translateSGRuleProtocolAll(r)
	case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp:
		return translateSGRuleProtocolTCPUDP(r)
	case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp:
		return translateSGRuleProtocolIcmp(r)
	}
	return nil, fmt.Errorf("error parsing rule number %d in sg %s in VPC %s", index, *sg.Name, *sg.VPC.Name)
}

func translateSGRuleProtocolAll(rule *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll) (sgRule *ir.SGRule, err error) {
	direction, err := translateDirection(*rule.Direction)
	if err != nil {
		return nil, err
	}
	remote, err := translateRemote(rule.Remote)
	if err != nil {
		return nil, err
	}
	local, err := translateLocal(rule.Local)
	if err != nil {
		return nil, err
	}
	return &ir.SGRule{Direction: direction, Remote: remote, Protocol: netp.AnyProtocol{}, Local: local}, nil
}

func translateSGRuleProtocolTCPUDP(rule *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp) (sgRule *ir.SGRule, err error) {
	direction, err := translateDirection(*rule.Direction)
	if err != nil {
		return nil, err
	}
	remote, err := translateRemote(rule.Remote)
	if err != nil {
		return nil, err
	}
	local, err := translateLocal(rule.Local)
	if err != nil {
		return nil, err
	}
	protocol, err := translateProtocolTCPUDP(rule)
	if err != nil {
		return nil, err
	}
	return &ir.SGRule{Direction: direction, Remote: remote, Protocol: protocol, Local: local}, nil
}

func translateSGRuleProtocolIcmp(rule *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp) (sgRule *ir.SGRule, err error) {
	direction, err := translateDirection(*rule.Direction)
	if err != nil {
		return nil, err
	}
	remote, err := translateRemote(rule.Remote)
	if err != nil {
		return nil, err
	}
	local, err := translateLocal(rule.Local)
	if err != nil {
		return nil, err
	}
	protocol, err := netp.ICMPFromTypeAndCode64(rule.Type, rule.Code)
	if err != nil {
		return nil, err
	}
	return &ir.SGRule{Direction: direction, Remote: remote, Protocol: protocol, Local: local}, nil
}

func translateDirection(direction string) (ir.Direction, error) {
	if direction == "inbound" {
		return ir.Inbound, nil
	} else if direction == "outbound" {
		return ir.Outbound, nil
	}
	return ir.Inbound, fmt.Errorf("SG rule direction must be either inbound or outbound")
}

func translateRemote(remote vpcv1.SecurityGroupRuleRemoteIntf) (ir.RemoteType, error) {
	if r, ok := remote.(*vpcv1.SecurityGroupRuleRemote); ok {
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

func translateLocal(local vpcv1.SecurityGroupRuleLocalIntf) (*netset.IPBlock, error) {
	var err error
	var ipAddrs *netset.IPBlock
	if l, ok := local.(*vpcv1.SecurityGroupRuleLocal); ok {
		if l.CIDRBlock != nil {
			ipAddrs, err = netset.IPBlockFromCidr(*l.CIDRBlock)
		}
		if l.Address != nil {
			ipAddrs, err = netset.IPBlockFromIPAddress(*l.CIDRBlock)
		}
		if err != nil {
			return nil, err
		}
		return verifyLocalValue(ipAddrs)
	}
	return nil, fmt.Errorf("error parsing Local field")
}

// temporary - first version of optimization requires local = 0.0.0.0/32
func verifyLocalValue(ipAddrs *netset.IPBlock) (*netset.IPBlock, error) {
	if !ipAddrs.Equal(netset.GetCidrAll()) {
		return nil, fmt.Errorf("only 0.0.0.0/32 CIDR block is supported for local values")
	}
	return ipAddrs, nil
}

func translateProtocolTCPUDP(rule *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp) (netp.Protocol, error) {
	isTCP := *rule.Protocol == vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudpProtocolTCPConst
	minDstPort := utils.GetProperty(rule.PortMin, netp.MinPort)
	maxDstPort := utils.GetProperty(rule.PortMax, netp.MaxPort)
	return netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(minDstPort), int(maxDstPort))
}
