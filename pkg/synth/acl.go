/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package synth generates Network ACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const ACLTypeNotSupported = "ACL: src/dst of type %s is not supported."

type ACLSynthesizer struct {
	spec      *ir.Spec
	singleACL bool
	result    *ir.ACLCollection
}

// NewACLSynthesizer creates and returns a new ACLSynthesizer instance
func NewACLSynthesizer(s *ir.Spec, single bool) Synthesizer {
	return &ACLSynthesizer{spec: s, singleACL: single, result: ir.NewACLCollection()}
}

func (a *ACLSynthesizer) Synth() ir.SynthCollection {
	return a.makeACL()
}

// makeACL translates Spec to a collection of nACLs
// 1. generate nACL rules for relevant subnets for each connection
// 2. generate nACL rules for blocked subnets (subnets that do not appear in Spec)
func (a *ACLSynthesizer) makeACL() *ir.ACLCollection {
	for c := range a.spec.Connections {
		a.generateACLRulesFromConnection(&a.spec.Connections[c])
	}
	a.generateACLRulesForBlockedSubnets()
	return a.result
}

//  1. check that both resources are supported in nACL generation.
//  2. check that at least one resource is internal.
//  3. convert src and dst resources to namedAddrs slices to make it more convenient to go through the addrs
//     and add the rule to the relevant acl. Note: in case where the resource in a nif, src/dst will be
//     updated to be its subnet.
//  4. generate rules and add them to relevant ACL to allow traffic for all pairs of IPAddrs of both resources.
func (a *ACLSynthesizer) generateACLRulesFromConnection(conn *ir.Connection) {
	if !resourceRelevantToACL(conn.Src.Type) {
		log.Fatalf(ACLTypeNotSupported, string(conn.Src.Type))
	}
	if !resourceRelevantToACL(conn.Dst.Type) {
		log.Fatalf(ACLTypeNotSupported, string(conn.Dst.Type))
	}
	internalSrc, internalDst, _ := internalConn(conn)
	if !internalSrc && !internalDst {
		log.Fatalf("ACL: Both source and destination are external for connection %v", *conn)
	}
	for _, src := range conn.Src.IPAddrs {
		srcSubnets, srcCidr := adjustResource(&a.spec.Defs, src, conn.Src)
		for _, dst := range conn.Dst.IPAddrs {
			dstSubnets, dstCidr := adjustResource(&a.spec.Defs, dst, conn.Dst)
			if src == dst && conn.Src.Type != ir.ResourceTypeCidr && conn.Dst.Type != ir.ResourceTypeCidr {
				continue
			}
			for _, trackedProtocol := range conn.TrackedProtocols {
				a.allowConnectionFromSrc(conn, trackedProtocol, srcSubnets, dstCidr)
				a.allowConnectionToDst(conn, trackedProtocol, dstSubnets, srcCidr)
			}
		}
	}
}

// if the src in internal, rule(s) will be created to allow traffic.
// if the protocol allows response, more rules will be created.
func (a *ACLSynthesizer) allowConnectionFromSrc(conn *ir.Connection, trackedProtocol ir.TrackedProtocol,
	srcSubnets []*namedAddrs, dstCidr *netset.IPBlock) {
	internalSrc, _, internal := internalConn(conn)

	if !internalSrc {
		return
	}
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}
	for _, srcSubnet := range srcSubnets {
		if srcSubnet.Addrs.Equal(dstCidr) { // srcSubnet and dstCidr are the same subnet
			continue
		}
		request := &ir.Packet{Src: srcSubnet.Addrs, Dst: dstCidr, Protocol: trackedProtocol.Protocol, Explanation: reason.String()}
		a.addRuleToACL(ir.AllowSend(request), srcSubnet, internal, a.singleACL)
		if inverseProtocol := trackedProtocol.Protocol.InverseDirection(); inverseProtocol != nil {
			response := &ir.Packet{Src: dstCidr, Dst: srcSubnet.Addrs, Protocol: inverseProtocol, Explanation: reason.response().String()}
			a.addRuleToACL(ir.AllowReceive(response), srcSubnet, internal, a.singleACL)
		}
	}
}

// if the dst in internal, rule(s) will be created to allow traffic.
// if the protocol allows response, more rules will be created.
func (a *ACLSynthesizer) allowConnectionToDst(conn *ir.Connection, trackedProtocol ir.TrackedProtocol,
	dstSubnets []*namedAddrs, srcCidr *netset.IPBlock) {
	_, internalDst, internal := internalConn(conn)

	if !internalDst {
		return
	}
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}
	for _, dstSubnet := range dstSubnets {
		if dstSubnet.Addrs.Equal(srcCidr) { // dstSubnet and srcCidr are the same subnet
			continue
		}
		request := &ir.Packet{Src: srcCidr, Dst: dstSubnet.Addrs, Protocol: trackedProtocol.Protocol, Explanation: reason.String()}
		a.addRuleToACL(ir.AllowReceive(request), dstSubnet, internal, a.singleACL)
		if inverseProtocol := trackedProtocol.Protocol.InverseDirection(); inverseProtocol != nil {
			response := &ir.Packet{Src: dstSubnet.Addrs, Dst: srcCidr, Protocol: inverseProtocol, Explanation: reason.response().String()}
			a.addRuleToACL(ir.AllowSend(response), dstSubnet, internal, a.singleACL)
		}
	}
}

// generate nACL rules for blocked subnets (subnets that do not appear in Spec)
func (a *ACLSynthesizer) generateACLRulesForBlockedSubnets() {
	blockedSubnets := a.spec.ComputeBlockedSubnets()
	for _, subnet := range blockedSubnets {
		acl := a.result.LookupOrCreate(aclSelector(subnet, a.singleACL))
		cidr := a.spec.Defs.Subnets[subnet].Address()
		acl.AppendInternal(ir.DenyAllReceive(subnet, cidr))
		acl.AppendInternal(ir.DenyAllSend(subnet, cidr))
	}
}

// convert src and dst resources to namedAddrs slices to make it more convenient to go through the addrs and add
// the rule to the relevant acl. Note: in case where the resource in a nif, src/dst will be updated to be its subnet.
func adjustResource(s *ir.Definitions, addrs *netset.IPBlock, resource ir.Resource) ([]*namedAddrs, *netset.IPBlock) {
	switch resource.Type {
	case ir.ResourceTypeSubnet:
		return adjustSubnet(s, addrs, resource.Name), addrs
	case ir.ResourceTypeExternal:
		return []*namedAddrs{{Name: resource.Name, Addrs: addrs}}, addrs
	case ir.ResourceTypeNIF:
		result := expandNifToSubnet(s, addrs)
		return result, result[0].Addrs // return nif's subnet, not its IP address
	case ir.ResourceTypeCidr:
		return adjustCidrSegment(s, addrs, resource.Name), addrs
	}
	return []*namedAddrs{}, nil // shouldn't happen
}

func adjustSubnet(s *ir.Definitions, addrs *netset.IPBlock, resourceName string) []*namedAddrs {
	// Todo: Handle the case where there is a subnet and a subnetSegment with the same name
	if subnetDetails, ok := s.Subnets[resourceName]; ok { // resource is a subnet
		return []*namedAddrs{{Name: resourceName, Addrs: subnetDetails.Address()}}
	}
	// resource is a subnet segment
	for _, subnetName := range s.SubnetSegments[resourceName].Subnets {
		if s.Subnets[subnetName].Address().Equal(addrs) {
			return []*namedAddrs{{Name: subnetName, Addrs: s.Subnets[subnetName].Address()}}
		}
	}
	return []*namedAddrs{} // shouldn't happen
}

func adjustCidrSegment(s *ir.Definitions, cidr *netset.IPBlock, resourceName string) []*namedAddrs {
	cidrSegmentDetails := s.CidrSegments[resourceName]
	cidrDetails := cidrSegmentDetails.Cidrs[cidr]
	result := make([]*namedAddrs, len(cidrDetails.ContainedSubnets))
	for i, subnet := range cidrDetails.ContainedSubnets {
		result[i] = &namedAddrs{Name: subnet, Addrs: s.Subnets[subnet].Address()}
	}
	return result
}

func expandNifToSubnet(s *ir.Definitions, addr *netset.IPBlock) []*namedAddrs {
	nifName, _ := s.NIFFromIP(addr) // already checked before (Lookup function) that the NIF exists
	subnetName := s.NIFs[nifName].Subnet
	subnetCidr := s.Subnets[subnetName].Address()

	return []*namedAddrs{{Name: subnetName, Addrs: subnetCidr}}
}

func aclSelector(subnetName ir.ID, single bool) string {
	if single {
		return fmt.Sprintf("%s/singleACL", ir.VpcFromScopedResource(subnetName))
	}
	return subnetName
}

func resourceRelevantToACL(e ir.ResourceType) bool {
	return e == ir.ResourceTypeSubnet || e == ir.ResourceTypeCidr || e == ir.ResourceTypeNIF || e == ir.ResourceTypeExternal
}

func (a *ACLSynthesizer) addRuleToACL(rule *ir.ACLRule, resource *namedAddrs, internal, single bool) {
	acl := a.result.LookupOrCreate(aclSelector(resource.Name, single))
	if internal {
		acl.AppendInternal(rule)
	} else {
		acl.AppendExternal(rule)
	}
}
