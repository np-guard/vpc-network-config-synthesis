/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"fmt"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const (
	icmpConst     = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmpProtocolIcmpConst
	allConst      = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllProtocolAllConst
	ipv4Const     = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllIPVersionIpv4Const
	allowConst    = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllActionAllowConst
	denyConst     = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllActionDenyConst
	outboundConst = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllDirectionOutboundConst
	inboundConst  = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllDirectionInboundConst
)

type refData struct {
	ID   *string
	CRN  *string
	Href *string
}

var globalIndex int = 0

func allocateRef() refData {
	globalIndex++
	return refData{
		ID:   utils.Ptr(fmt.Sprintf("fake:id:%v", globalIndex)),
		CRN:  utils.Ptr(fmt.Sprintf("fake:crn:%v", globalIndex)),
		Href: utils.Ptr(fmt.Sprintf("fake:href:%v", globalIndex)),
	}
}

func action(a ir.Action) *string {
	switch a {
	case ir.Allow:
		return utils.Ptr(allowConst)
	case ir.Deny:
		return utils.Ptr(denyConst)
	}
	return nil
}

func direction(d ir.Direction) *string {
	switch d {
	case ir.Outbound:
		return utils.Ptr(outboundConst)
	case ir.Inbound:
		return utils.Ptr(inboundConst)
	}
	return nil
}

func minPort(r ir.PortRange) *int64 {
	return utils.Ptr(int64(r.Min))
}

func maxPort(r ir.PortRange) *int64 {
	return utils.Ptr(int64(r.Max))
}

type tcpudpData struct {
	Protocol   *string
	srcPortMin *int64
	srcPortMax *int64
	dstPortMin *int64
	dstPortMax *int64
}

type icmpData struct {
	Protocol *string
	Type     *int64
	Code     *int64
}

type allData struct {
	Protocol *string
}

func tcpudp(p ir.TCPUDP) tcpudpData {
	r := p.PortRangePair
	res := tcpudpData{
		Protocol:   utils.Ptr(strings.ToLower(string(p.Protocol))),
		srcPortMin: minPort(r.SrcPort),
		srcPortMax: maxPort(r.SrcPort),
		dstPortMin: minPort(r.DstPort),
		dstPortMax: maxPort(r.DstPort),
	}
	return res
}

func icmp(p ir.ICMP) icmpData {
	res := icmpData{
		Protocol: utils.Ptr(icmpConst),
		Type:     utils.Ptr(int64(p.Type)),
	}
	if p.Code != nil {
		res.Code = utils.Ptr(int64(*p.Code))
	}
	return res
}

func all() allData {
	return allData{
		Protocol: utils.Ptr(allConst),
	}
}
