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

const subnetNotFoundError = "ACL: src/dst of type network interface (or instance) is not supported."

type Options struct {
	SingleACL bool
}

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *ir.Spec, opt Options) *ir.ACLCollection {
	aclSelector := func(cidr *ipblock.IPBlock) string {
		result, ok := s.Defs.SubnetNameFromIP(cidr)
		if !ok {
			log.Fatalf(subnetNotFoundError)
		}
		return result
	}
	if opt.SingleACL {
		aclSelector = func(cidr *ipblock.IPBlock) string {
			result, ok := s.Defs.SubnetNameFromIP(cidr)
			if !ok {
				log.Fatalf(subnetNotFoundError)
			}
			return fmt.Sprintf("%s/singleACL", ir.VpcFromScopedResource(result))
		}
	}
	collections := []*ir.ACLCollection{}
	for c := range s.Connections {
		collection := GenerateACLCollectionFromConnection(s, &s.Connections[c], aclSelector)
		collections = append(collections, collection)
	}
	return ir.MergeACLCollections(collections...)
}

func GenerateACLCollectionFromConnection(s *ir.Spec, conn *ir.Connection,
	aclSelector func(target *ipblock.IPBlock) string) *ir.ACLCollection {
	internalSrc := conn.Src.Type != ir.ResourceTypeExternal
	internalDst := conn.Dst.Type != ir.ResourceTypeExternal
	internal := internalSrc && internalDst
	if !internalSrc && !internalDst {
		log.Fatalf("ACL: Both source and destination are external for connection %v", *conn)
	}
	result := ir.NewACLCollection()
	if !resourceRelevantToACL(conn.Src.Type) && !resourceRelevantToACL(conn.Dst.Type) {
		return result
	}
	var connectionRules []*ir.ACLRule
	for _, src := range conn.Src.Values {
		for _, dst := range conn.Dst.Values {
			if src == dst && conn.Src.Type != ir.ResourceTypeCidr && conn.Dst.Type != ir.ResourceTypeCidr {
				continue
			}
			for _, trackedProtocol := range conn.TrackedProtocols {
				reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}
				protocolRules := allowDirectedConnection(s, src, dst, conn, internalSrc, internalDst, trackedProtocol.Protocol, reason)
				connectionRules = append(connectionRules, protocolRules...)
			}
		}
	}
	for _, rule := range connectionRules {
		acl := result.LookupOrCreate(aclSelector(rule.Target()))
		if internal {
			acl.AppendInternal(rule)
		} else {
			acl.AppendExternal(rule)
		}
	}
	return result
}

func allowDirectedConnection(s *ir.Spec, srcCidr, dstCidr *ipblock.IPBlock, conn *ir.Connection, internalSrc, internalDst bool,
	protocol ir.Protocol, reason explanation) []*ir.ACLRule {
	var request, response *ir.Packet

	srcSubnetList := subnetsContainedInCidr(s, srcCidr, conn.Src)
	dstSubnetList := subnetsContainedInCidr(s, dstCidr, conn.Dst)

	var connection []*ir.ACLRule

	if internalSrc {
		for _, srcSubnetCidr := range srcSubnetList {
			if srcSubnetCidr == dstCidr {
				continue
			}
			request = &ir.Packet{Src: srcSubnetCidr, Dst: dstCidr, Protocol: protocol, Explanation: reason.String()}
			connection = append(connection, ir.AllowSend(*request))
			if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
				response = &ir.Packet{Src: dstCidr, Dst: srcSubnetCidr, Protocol: inverseProtocol, Explanation: reason.response().String()}
				connection = append(connection, ir.AllowReceive(*response))
			}
		}
	}

	if internalDst {
		for _, dstSubnetCidr := range dstSubnetList {
			if srcCidr == dstSubnetCidr {
				continue
			}
			request = &ir.Packet{Src: srcCidr, Dst: dstSubnetCidr, Protocol: protocol, Explanation: reason.String()}
			connection = append(connection, ir.AllowReceive(*request))
			if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
				response = &ir.Packet{Src: dstSubnetCidr, Dst: srcCidr, Protocol: inverseProtocol, Explanation: reason.response().String()}
				connection = append(connection, ir.AllowSend(*response))
			}
		}
	}

	return connection
}

func resourceRelevantToACL(e ir.ResourceType) bool {
	return e == ir.ResourceTypeSubnet || e == ir.ResourceTypeSegment || e == ir.ResourceTypeCidr
}

func subnetsContainedInCidr(s *ir.Spec, cidr *ipblock.IPBlock, resource ir.Resource) []*ipblock.IPBlock {
	if resource.Type != ir.ResourceTypeCidr {
		return []*ipblock.IPBlock{cidr}
	}
	cidrSegmentDetails := s.Defs.CidrSegments[resource.Name]
	cidrDetails := cidrSegmentDetails.Cidrs[cidr]
	result := make([]*ipblock.IPBlock, len(cidrDetails.ContainedSubnets))
	for i, subnet := range cidrDetails.ContainedSubnets {
		result[i] = s.Defs.Subnets[subnet].Address()
	}
	return result
}
