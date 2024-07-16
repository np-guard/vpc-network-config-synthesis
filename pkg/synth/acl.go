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
// (1) generate a nACL collection for relevant subnets for each connection
// (2) generate a nACL collection for blocked subnets (subnets that does not appear in Spec)
// (3) merge and return the nACL collection
func MakeACL(s *ir.Spec, opt Options) *ir.ACLCollection {
	collections := []*ir.ACLCollection{}
	for c := range s.Connections {
		collection := generateACLCollectionFromConnection(&s.Defs, &s.Connections[c], opt.SingleACL)
		collections = append(collections, collection)
	}
	collections = append(collections, generateACLCollectionForBlockedSubnets(s, opt.SingleACL))
	return ir.MergeACLCollections(collections...)
}

// (1) check that both resources are supported in nACL generation.
// (2) check that at least one resource is internal
// (3) for all pairs of IPAddrs of both resources, rules that allow communication will be created
func generateACLCollectionFromConnection(s *ir.Definitions, conn *ir.Connection, single bool) *ir.ACLCollection {
	if !resourceRelevantToACL(conn.Src.Type) {
		log.Fatalf(fmt.Sprintf(ACLTypeNotSupported, string(conn.Src.Type)))
	}
	if !resourceRelevantToACL(conn.Dst.Type) {
		log.Fatalf(fmt.Sprintf(ACLTypeNotSupported, string(conn.Dst.Type)))
	}
	internalSrc := conn.Src.Type != ir.ResourceTypeExternal
	internalDst := conn.Dst.Type != ir.ResourceTypeExternal
	internal := internalSrc && internalDst
	if !internalSrc && !internalDst {
		log.Fatalf("ACL: Both source and destination are external for connection %v", *conn)
	}
	result := ir.NewACLCollection()
	for _, src := range conn.Src.IPAddrs {
		srcCidrs, src := adjustResource(s, src, conn.Src)
		for _, dst := range conn.Dst.IPAddrs {
			dstCidrs, dst := adjustResource(s, dst, conn.Dst)
			if src == dst && conn.Src.Type != ir.ResourceTypeCidr && conn.Dst.Type != ir.ResourceTypeCidr {
				continue
			}
			for _, trackedProtocol := range conn.TrackedProtocols {
				reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}
				allowDirectedConnection(src, dst, srcCidrs, dstCidrs, internalSrc, internalDst, trackedProtocol.Protocol, reason, result, single)
			}
		}
	}
	return result
}

// create relevant rule(s) to allow traffic
func allowDirectedConnection(srcCidr, dstCidr *ipblock.IPBlock, srcSubnets, dstSubnets []*ir.ConnResource, internalSrc, internalDst bool,
	protocol ir.Protocol, reason explanation, result *ir.ACLCollection, singleACL bool) {
	var request, response *ir.Packet
	internal := internalSrc && internalDst

	if internalSrc {
		for _, srcSubnet := range srcSubnets {
			if srcSubnet.Addrs.Equal(dstCidr) { // srcSubnet and dstCidr are the same subnet
				continue
			}
			request = &ir.Packet{Src: srcSubnet.Addrs, Dst: dstCidr, Protocol: protocol, Explanation: reason.String()}
			addRuleToACL(result, ir.AllowSend(*request), srcSubnet, internal, singleACL)
			if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
				response = &ir.Packet{Src: dstCidr, Dst: srcSubnet.Addrs, Protocol: inverseProtocol, Explanation: reason.response().String()}
				addRuleToACL(result, ir.AllowReceive(*response), srcSubnet, internal, singleACL)
			}
		}
	}

	if internalDst {
		for _, dstSubnet := range dstSubnets {
			if dstSubnet.Addrs.Equal(srcCidr) { // dstSubnet and srcCidr are the same subnet
				continue
			}
			request = &ir.Packet{Src: srcCidr, Dst: dstSubnet.Addrs, Protocol: protocol, Explanation: reason.String()}
			addRuleToACL(result, ir.AllowReceive(*request), dstSubnet, internal, singleACL)
			if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
				response = &ir.Packet{Src: dstSubnet.Addrs, Dst: srcCidr, Protocol: inverseProtocol, Explanation: reason.response().String()}
				addRuleToACL(result, ir.AllowSend(*response), dstSubnet, internal, singleACL)
			}
		}
	}
}

// generate a nACL collection, which include nACL(s) rules that block all traffic to/from the blocked subnets
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

// Change ir.Resource to a slice of ir.ConnResource (each ir.ConnResource is either a pair of a of a subnet with its CIDR or an
// external with its IPAddrs), to make it convenient to go through all pairs of subnets in allowDirectedConnection function.
func adjustResource(s *ir.Definitions, addrs *ipblock.IPBlock, resource ir.Resource) ([]*ir.ConnResource, *ipblock.IPBlock) {
	switch resource.Type {
	case ir.ResourceTypeSubnet:
		return adjustSubnet(s, addrs, resource.Name), addrs
	case ir.ResourceTypeExternal:
		connResource := ir.ConnResource{Name: resource.Name, Addrs: addrs}
		return []*ir.ConnResource{&connResource}, addrs
	case ir.ResourceTypeNIF:
		result := expandNifToSubnet(s, addrs)
		return result, result[0].Addrs // return nif's subnet, not its IP address
	case ir.ResourceTypeCidr:
		return adjustCidrSegment(s, addrs, resource.Name), addrs
	}
	return []*ir.ConnResource{}, nil // dead code
}

func adjustSubnet(s *ir.Definitions, addrs *ipblock.IPBlock, resourceName string) []*ir.ConnResource {
	// Todo: Handle the case where there is a subnet and a subnetSegment with the same name
	if subnetDetails, ok := s.Subnets[resourceName]; ok { // resource is a subnet
		connResource := ir.ConnResource{Name: resourceName, Addrs: subnetDetails.Address()}
		return []*ir.ConnResource{&connResource}
	}
	// resource is a subnet segment
	for _, subnetName := range s.SubnetSegments[resourceName].Subnets {
		if s.Subnets[subnetName].Address().Equal(addrs) {
			connResource := ir.ConnResource{Name: subnetName, Addrs: s.Subnets[subnetName].Address()}
			return []*ir.ConnResource{&connResource}
		}
	}
	return []*ir.ConnResource{} // dead code
}

func adjustCidrSegment(s *ir.Definitions, cidr *ipblock.IPBlock, resourceName string) []*ir.ConnResource {
	cidrSegmentDetails := s.CidrSegments[resourceName]
	cidrDetails := cidrSegmentDetails.Cidrs[cidr]
	result := make([]*ir.ConnResource, len(cidrDetails.ContainedSubnets))
	for i, subnet := range cidrDetails.ContainedSubnets {
		connResource := ir.ConnResource{Name: subnet, Addrs: s.Subnets[subnet].Address()}
		result[i] = &connResource
	}
	return result
}

func expandNifToSubnet(s *ir.Definitions, addr *ipblock.IPBlock) []*ir.ConnResource {
	nifName, _ := s.NIFFromIP(addr) // already checked before (Lookup function) that the NIF exists
	subnetName := s.NIFs[nifName].Subnet
	subnetCidr := s.Subnets[subnetName].Address()

	connResource := ir.ConnResource{Name: subnetName, Addrs: subnetCidr}
	return []*ir.ConnResource{&connResource}
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

func addRuleToACL(result *ir.ACLCollection, rule *ir.ACLRule, resource *ir.ConnResource, internal, single bool) {
	acl := result.LookupOrCreate(aclSelector(resource.Name, single))
	if internal {
		acl.AppendInternal(rule)
	} else {
		acl.AppendExternal(rule)
	}
}
