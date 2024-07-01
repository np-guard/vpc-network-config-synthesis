/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package synth generates Network ACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const ACLTypeNotSupported = "ACL: src/dst of type %s is not supported."

type Options struct {
	SingleACL bool
}

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *ir.Spec, opt Options) *ir.ACLCollection {
	collections := []*ir.ACLCollection{}
	for c := range s.Connections {
		collection := generateACLCollectionFromConnection(&s.Defs, &s.Connections[c], opt.SingleACL)
		collections = append(collections, collection)
	}
	collections = append(collections, generateACLCollectionForBlockedSubnets(s, opt.SingleACL))
	return ir.MergeACLCollections(collections...)
}

func generateACLCollectionFromConnection(s *ir.Definitions, conn *ir.Connection, single bool) *ir.ACLCollection {
	internalSrc := conn.Src.Type != ir.ResourceTypeExternal
	internalDst := conn.Dst.Type != ir.ResourceTypeExternal
	internal := internalSrc && internalDst
	if !internalSrc && !internalDst {
		log.Fatalf("ACL: Both source and destination are external for connection %v", *conn)
	}
	result := ir.NewACLCollection()
	if !resourceRelevantToACL(conn.Src.Type) {
		log.Fatalf(fmt.Sprintf(ACLTypeNotSupported, string(conn.Src.Type)))
	}
	if !resourceRelevantToACL(conn.Dst.Type) {
		log.Fatalf(fmt.Sprintf(ACLTypeNotSupported, string(conn.Dst.Type)))
	}
	for _, src := range conn.Src.IPAddrs {
		for _, dst := range conn.Dst.IPAddrs {
			if src == dst && conn.Src.Type != ir.ResourceTypeCidr && conn.Dst.Type != ir.ResourceTypeCidr {
				continue
			}
			srcEndpoints, src := updateResource(s, src, conn.Src, single)
			dstEndpoints, dst := updateResource(s, dst, conn.Dst, single)
			for _, trackedProtocol := range conn.TrackedProtocols {
				reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}
				allowDirectedConnection(src, dst, srcEndpoints, dstEndpoints, internalSrc, internalDst, trackedProtocol.Protocol, reason, result)
			}
		}
	}
	return result
}

func allowDirectedConnection(srcCidr, dstCidr *ipblock.IPBlock, srcEndpoints, dstEndpoints []ir.ConnResource, internalSrc, internalDst bool,
	protocol ir.Protocol, reason explanation, result *ir.ACLCollection) {
	var request, response *ir.Packet
	internal := internalSrc && internalDst

	if internalSrc {
		for _, srcSubnet := range srcEndpoints {
			if srcSubnet.Addrs.Equal(dstCidr) { // srcSubnet and dstCidr are the same subnet
				continue
			}
			request = &ir.Packet{Src: srcSubnet.Addrs, Dst: dstCidr, Protocol: protocol, Explanation: reason.String()}
			addRuleToACL(result, ir.AllowSend(*request), srcSubnet, internal)
			if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
				response = &ir.Packet{Src: dstCidr, Dst: srcSubnet.Addrs, Protocol: inverseProtocol, Explanation: reason.response().String()}
				addRuleToACL(result, ir.AllowReceive(*response), srcSubnet, internal)
			}
		}
	}

	if internalDst {
		for _, dstSubnet := range dstEndpoints {
			if dstSubnet.Addrs.Equal(srcCidr) {
				continue
			}
			request = &ir.Packet{Src: srcCidr, Dst: dstSubnet.Addrs, Protocol: protocol, Explanation: reason.String()}
			addRuleToACL(result, ir.AllowReceive(*request), dstSubnet, internal)
			if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
				response = &ir.Packet{Src: dstSubnet.Addrs, Dst: srcCidr, Protocol: inverseProtocol, Explanation: reason.response().String()}
				addRuleToACL(result, ir.AllowSend(*response), dstSubnet, internal)
			}
		}
	}
}

func generateACLCollectionForBlockedSubnets(s *ir.Spec, single bool) *ir.ACLCollection {
	blockedSubnets := s.ComputeBlockedSubnets()
	result := ir.NewACLCollection()
	for _, subnet := range blockedSubnets {
		acl := result.LookupOrCreate(aclSelector(subnet, single))
		cidr := s.Defs.Subnets[subnet].Address()
		acl.AppendInternal(ir.DenyAllReceive(subnet, cidr))
		acl.AppendInternal(ir.DenyAllSend(subnet, cidr))
	}
	return result
}

func updateResource(s *ir.Definitions, addrs *ipblock.IPBlock, resource ir.Resource, single bool) ([]ir.ConnResource, *ipblock.IPBlock) {
	switch resource.Type {
	case ir.ResourceTypeSubnet:
		return updateSubnet(s, addrs, resource.Name, single), addrs
	case ir.ResourceTypeExternal:
		connResource := ir.ConnResource{Name: resource.Name, Addrs: addrs}
		return []ir.ConnResource{connResource}, addrs
	case ir.ResourceTypeNIF:
		result := updateNif(s, addrs, single)
		return result, result[0].Addrs // should return nif's subnet
	case ir.ResourceTypeCidr:
		return updateCidrSegment(s, addrs, resource.Name, single), addrs
	}
	return []ir.ConnResource{}, nil // dead code
}

func updateCidrSegment(s *ir.Definitions, cidr *ipblock.IPBlock, resourceName string, single bool) []ir.ConnResource {
	cidrSegmentDetails := s.CidrSegments[resourceName]
	cidrDetails := cidrSegmentDetails.Cidrs[cidr]
	result := make([]ir.ConnResource, len(cidrDetails.ContainedSubnets))
	for i, subnet := range cidrDetails.ContainedSubnets {
		connResource := ir.ConnResource{Name: aclSelector(subnet, single), Addrs: s.Subnets[subnet].Address()}
		result[i] = connResource
	}
	return result
}

func updateNif(s *ir.Definitions, addr *ipblock.IPBlock, single bool) []ir.ConnResource {
	nifName, _ := s.NIFFromIP(addr)
	subnetName := s.NIFs[nifName].Subnet
	subnetCidr := s.Subnets[subnetName].Address()

	connResource := ir.ConnResource{Name: aclSelector(subnetName, single), Addrs: subnetCidr}
	return []ir.ConnResource{connResource}
}

func updateSubnet(s *ir.Definitions, addrs *ipblock.IPBlock, resourceName string, single bool) []ir.ConnResource {
	if subnetDetails, ok := s.Subnets[resourceName]; ok { // resource is a subnet
		connResource := ir.ConnResource{Name: aclSelector(resourceName, single), Addrs: subnetDetails.Address()}
		return []ir.ConnResource{connResource}
	}
	// resource is a subnet segment
	for _, subnetName := range s.SubnetSegments[resourceName].Subnets {
		if s.Subnets[subnetName].Address().Equal(addrs) {
			connResource := ir.ConnResource{Name: aclSelector(subnetName, single), Addrs: s.Subnets[subnetName].Address()}
			return []ir.ConnResource{connResource}
		}
	}
	return []ir.ConnResource{} // dead code
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

func addRuleToACL(result *ir.ACLCollection, rule *ir.ACLRule, resource ir.ConnResource, internal bool) {
	acl := result.LookupOrCreate(resource.Name)
	if internal {
		acl.AppendInternal(rule)
	} else {
		acl.AppendExternal(rule)
	}
}
