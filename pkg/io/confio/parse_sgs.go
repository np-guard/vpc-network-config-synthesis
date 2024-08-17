/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"fmt"

	vpc1 "github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

// get a slice of ir.SG from a config_object file
func ReadSGs(filename string) ([]ir.SG, error) {
	config, err := readModel(filename)
	if err != nil {
		return nil, err
	}
	result := make([]ir.SG, len(config.SecurityGroupList))

	for i, sg := range config.SecurityGroupList {
		inbound, outbound, err := translateSGRules(&sg.SecurityGroup)
		if err != nil {
			return nil, err
		}
		result[i] = ir.SG{InboundRules: inbound, OutboundRules: outbound}
	}
	return result, nil
}

// parse ingress and egress rules of a security group
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

// translate the rule to ir.SGRule
func translateSGRule(sg *vpc1.SecurityGroup, index int) (sgRule ir.SGRule, err error) {
	if r, ok := sg.Rules[index].(*vpc1.SecurityGroupRule); ok {
		direction, err := translateDirection(*r.Direction)
		if err != nil {
			return ir.SGRule{}, err
		}
		remote, err := translateRemote(r.Remote)
		if err != nil {
			return ir.SGRule{}, err
		}
		protocol, err := translateProtocol(sg.Rules[index])
		if err != nil {
			return ir.SGRule{}, err
		}
		local, err := translateLocal(r.Local)
		if err != nil {
			return ir.SGRule{}, err
		}

		return ir.SGRule{Direction: direction, Remote: remote, Protocol: protocol, Local: local}, nil
	}
	return ir.SGRule{}, fmt.Errorf("error parsing rule number %d in %s sg", index, *sg.Name)
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
	switch r := remote.(type) {
	case *vpc1.SecurityGroupRuleRemoteSecurityGroupReference:
		return ir.SGName(*r.Name), nil
	case *vpc1.SecurityGroupRuleRemoteCIDR:
		return netset.IPBlockFromCidr(*r.CIDRBlock)
	case *vpc1.SecurityGroupRuleRemoteIP:
		return netset.IPBlockFromIPAddress(*r.Address)
	}

	return nil, fmt.Errorf("unexpected SG rule remote")
}

func translateProtocol(rule vpc1.SecurityGroupRuleIntf) (netp.Protocol, error) {
	switch r := rule.(type) {
	case *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolAll:
		return netp.AnyProtocol{}, nil
	case *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp:
		return translateProtocolTCPUDP(r)
	case *vpc1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp:
		return netp.ICMPFromTypeAndCode64(r.Type, r.Code)
	default:
		return nil, fmt.Errorf("unsupported rule protocol type")
	}
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
	return newTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(minDstPort), int(maxDstPort)) // Todo: replace with netp.NewTCPUDP
}

// Todo: remove func when netp.NewTCPUDP is added
func newTCPUDP(isTCP bool, minSrcPort, maxSrcPort, minDstPort, maxDstPort int) (netp.TCPUDP, error) {
	if min(minSrcPort, minDstPort) < netp.MinPort || max(maxSrcPort, maxDstPort) > netp.MaxPort {
		return netp.TCPUDP{}, fmt.Errorf("TCPUDP ports are in range %d-%d", netp.MinPort, netp.MaxPort)
	}
	srcPorts := interval.New(int64(minSrcPort), int64(maxSrcPort))
	dstPorts := interval.New(int64(minDstPort), int64(maxDstPort))
	return netp.TCPUDP{IsTCP: isTCP, PortRangePair: netp.PortRangePair{SrcPort: srcPorts, DstPort: dstPorts}}, nil
}
