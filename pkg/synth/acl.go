/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package synth generates Network ACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"

	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const ACLTypeNotSupported = "ACL: src/dst of type %s is not supported"

type ACLSynthesizer struct {
	spec      *ir.Spec
	singleACL bool
	result    *ir.ACLCollection
}

// NewACLSynthesizer creates and returns a new ACLSynthesizer instance
func NewACLSynthesizer(s *ir.Spec, single bool) Synthesizer {
	return &ACLSynthesizer{spec: s, singleACL: single, result: ir.NewACLCollection()}
}

func (a *ACLSynthesizer) Synth() ir.Collection {
	return a.makeACL()
}

// makeACL translates Spec to a collection of nACLs
// 1. generate nACL rules for relevant subnets for each connection
// 2. generate nACL rules for blocked subnets (subnets that do not appear in Spec)
func (a *ACLSynthesizer) makeACL() *ir.ACLCollection {
	for c := range a.spec.Connections {
		a.generateACLRulesFromConnection(a.spec.Connections[c])
	}
	a.generateACLRulesForBlockedSubnets()
	return a.result
}

func (a *ACLSynthesizer) generateACLRulesFromConnection(conn *ir.Connection) {
	for _, srcSubnet := range conn.Src.NamedAddrs {
		for _, dstCidr := range conn.Dst.Cidrs {
			if srcSubnet.IPAddrs.Equal(dstCidr.IPAddrs) && *conn.Src.Type != ir.ResourceTypeCidr && *conn.Dst.Type != ir.ResourceTypeCidr {
				continue
			}
			for _, trackedProtocol := range conn.TrackedProtocols {
				a.allowConnectionFromSrc(conn, trackedProtocol, srcSubnet, dstCidr.IPAddrs)
			}
		}
	}

	for _, srcCidr := range conn.Src.Cidrs {
		for _, dstSubnet := range conn.Dst.NamedAddrs {
			if dstSubnet.IPAddrs.Equal(srcCidr.IPAddrs) && *conn.Src.Type != ir.ResourceTypeCidr && *conn.Dst.Type != ir.ResourceTypeCidr {
				continue
			}
			for _, trackedProtocol := range conn.TrackedProtocols {
				a.allowConnectionToDst(conn, trackedProtocol, dstSubnet, srcCidr.IPAddrs)
			}
		}
	}
}

// if the src in internal, rule(s) will be created to allow traffic.
// if the protocol allows response, more rules will be created.
func (a *ACLSynthesizer) allowConnectionFromSrc(conn *ir.Connection, p *ir.TrackedProtocol, srcSubnet *ir.NamedAddrs, dstCidr *netset.IPBlock) {
	internalSrc, _, internal := internalConn(conn)

	if !internalSrc || srcSubnet.IPAddrs.Equal(dstCidr) {
		return
	}
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: p.Origin}
	request := &ir.Packet{Src: srcSubnet.IPAddrs, Dst: dstCidr, Protocol: p.Protocol, Explanation: reason.String()}
	a.addRuleToACL(ir.AllowSend(request), *srcSubnet.Name, internal, a.singleACL)
	if inverseProtocol := p.Protocol.InverseDirection(); inverseProtocol != nil {
		response := &ir.Packet{Src: dstCidr, Dst: srcSubnet.IPAddrs, Protocol: inverseProtocol, Explanation: reason.response().String()}
		a.addRuleToACL(ir.AllowReceive(response), *srcSubnet.Name, internal, a.singleACL)
	}
}

// if the dst in internal, rule(s) will be created to allow traffic.
// if the protocol allows response, more rules will be created.
func (a *ACLSynthesizer) allowConnectionToDst(conn *ir.Connection, p *ir.TrackedProtocol, dstSubnet *ir.NamedAddrs, srcCidr *netset.IPBlock) {
	_, internalDst, internal := internalConn(conn)

	if !internalDst || dstSubnet.IPAddrs.Equal(srcCidr) {
		return
	}
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: p.Origin}
	request := &ir.Packet{Src: srcCidr, Dst: dstSubnet.IPAddrs, Protocol: p.Protocol, Explanation: reason.String()}
	a.addRuleToACL(ir.AllowReceive(request), *dstSubnet.Name, internal, a.singleACL)
	if inverseProtocol := p.Protocol.InverseDirection(); inverseProtocol != nil {
		response := &ir.Packet{Src: dstSubnet.IPAddrs, Dst: srcCidr, Protocol: inverseProtocol, Explanation: reason.response().String()}
		a.addRuleToACL(ir.AllowSend(response), *dstSubnet.Name, internal, a.singleACL)
	}
}

func aclSelector(subnetName ir.ID, single bool) string {
	if single {
		return fmt.Sprintf("%s/singleACL", ir.VpcFromScopedResource(subnetName))
	}
	return subnetName
}

func (a *ACLSynthesizer) addRuleToACL(rule *ir.ACLRule, resourceName ir.ID, internal, single bool) {
	acl := a.result.LookupOrCreate(aclSelector(resourceName, single))
	if internal {
		acl.AppendInternal(rule)
	} else {
		acl.AppendExternal(rule)
	}
}

// generate nACL rules for blocked subnets (subnets that do not appear in Spec)
func (a *ACLSynthesizer) generateACLRulesForBlockedSubnets() {
	blockedSubnets := utils.TrueKeyValues(a.spec.Defs.BlockedSubnets)
	ir.PrintUnspecifiedWarning(ir.WarningUnspecifiedACL, blockedSubnets)
	for _, subnet := range blockedSubnets {
		acl := a.result.LookupOrCreate(aclSelector(subnet, a.singleACL))
		cidr := a.spec.Defs.Subnets[subnet].Address()
		acl.AppendInternal(ir.DenyAllReceive(subnet, cidr))
		acl.AppendInternal(ir.DenyAllSend(subnet, cidr))
	}
}
