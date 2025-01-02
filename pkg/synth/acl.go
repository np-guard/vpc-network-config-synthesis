/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package synth generates Network ACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type ACLSynthesizer struct {
	spec      *ir.Spec
	singleACL bool
	result    *ir.ACLCollection
}

const WarningUnspecifiedACL = "The following subnets do not have required connections; the generated ACL will block all traffic: "

// NewACLSynthesizer creates and returns a new ACLSynthesizer instance
func NewACLSynthesizer(s *ir.Spec, single bool) Synthesizer {
	return &ACLSynthesizer{spec: s, singleACL: single, result: ir.NewACLCollection()}
}

func (a *ACLSynthesizer) Synth() (collection ir.Collection, warning string) {
	return a.makeACL()
}

// makeACL translates Spec to a collection of nACLs
// 1. generate nACL rules for relevant subnets for each connection
// 2. generate nACL rules for blocked subnets (subnets that do not appear in Spec)
func (a *ACLSynthesizer) makeACL() (collection *ir.ACLCollection, warning string) {
	for _, conn := range a.spec.Connections {
		a.generateACLRulesFromConnection(conn, conn.Src, conn.Dst, a.allowConnectionSrc)
		a.generateACLRulesFromConnection(conn, conn.Dst, conn.Src, a.allowConnectionDst)
	}
	warning = a.generateACLRulesForBlockedSubnets()
	return a.result, warning
}

func (a *ACLSynthesizer) generateACLRulesFromConnection(conn *ir.Connection, thisResource, otherResource *ir.ConnectedResource,
	allowConnection func(*ir.Connection, *ir.TrackedProtocol, *ir.NamedAddrs, *netset.IPBlock)) {
	for _, thisSubnet := range thisResource.CidrsWhenLocal {
		for _, otherCidr := range otherResource.CidrsWhenRemote {
			if thisSubnet.IPAddrs.Equal(otherCidr.IPAddrs) {
				continue
			}
			for _, trackedProtocol := range conn.TrackedProtocols {
				allowConnection(conn, trackedProtocol, thisSubnet, otherCidr.IPAddrs)
			}
		}
	}
}

// if the src in internal, rule(s) will be created to allow traffic.
// if the protocol allows response, more rules will be created.
func (a *ACLSynthesizer) allowConnectionSrc(conn *ir.Connection, p *ir.TrackedProtocol, srcSubnet *ir.NamedAddrs, dstCidr *netset.IPBlock) {
	internalSrc, _, internal := internalConnection(conn)

	if !internalSrc {
		return
	}
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: p.Origin}
	request := &ir.Packet{Src: srcSubnet.IPAddrs, Dst: dstCidr, Protocol: p.Protocol, Explanation: reason.String()}
	a.addRuleToACL(ir.AllowSend(request), srcSubnet.Name, internal)
	if inverseProtocol := p.Protocol.InverseDirection(); inverseProtocol != nil {
		response := &ir.Packet{Src: dstCidr, Dst: srcSubnet.IPAddrs, Protocol: inverseProtocol, Explanation: reason.response().String()}
		a.addRuleToACL(ir.AllowReceive(response), srcSubnet.Name, internal)
	}
}

// if the dst in internal, rule(s) will be created to allow traffic.
// if the protocol allows response, more rules will be created.
func (a *ACLSynthesizer) allowConnectionDst(conn *ir.Connection, p *ir.TrackedProtocol, dstSubnet *ir.NamedAddrs, srcCidr *netset.IPBlock) {
	_, internalDst, internal := internalConnection(conn)

	if !internalDst {
		return
	}
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: p.Origin}
	request := &ir.Packet{Src: srcCidr, Dst: dstSubnet.IPAddrs, Protocol: p.Protocol, Explanation: reason.String()}
	a.addRuleToACL(ir.AllowReceive(request), dstSubnet.Name, internal)
	if inverseProtocol := p.Protocol.InverseDirection(); inverseProtocol != nil {
		response := &ir.Packet{Src: dstSubnet.IPAddrs, Dst: srcCidr, Protocol: inverseProtocol, Explanation: reason.response().String()}
		a.addRuleToACL(ir.AllowSend(response), dstSubnet.Name, internal)
	}
}

func (a *ACLSynthesizer) addRuleToACL(rule *ir.ACLRule, subnetName ir.ID, internal bool) {
	acl := a.result.LookupOrCreate(subnetName, a.singleACL)
	if internal {
		acl.AppendInternal(rule)
	} else {
		acl.AppendExternal(rule)
	}
}

// generate nACL rules for blocked subnets (subnets that do not appear in Spec)
func (a *ACLSynthesizer) generateACLRulesForBlockedSubnets() string {
	blockedSubnets := utils.TrueKeyValues(a.spec.BlockedSubnets)
	for _, subnet := range blockedSubnets {
		acl := a.result.LookupOrCreate(subnet, a.singleACL)
		cidr := a.spec.Defs.Subnets[subnet].Address()
		acl.AppendInternal(ir.DenyAllReceive(subnet, cidr))
		acl.AppendInternal(ir.DenyAllSend(subnet, cidr))
	}
	return setUnspecifiedWarning(WarningUnspecifiedACL, blockedSubnets)
}
