// Package synth generates Network ACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

type Options struct {
	SingleACL bool
}

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *ir.Spec, opt Options) *ir.ACLCollection {
	aclSelector := func(ip ir.IP) string {
		result, ok := s.Defs.SubnetNameFromIP(ip)
		if !ok {
			log.Fatalf("ACL: src/dst of type network interface (or instance) is not supported.")
		}
		return result
	}
	if opt.SingleACL {
		aclSelector = func(target ir.IP) string { return "1" }
	}
	collections := []*ir.ACLCollection{}
	for c := range s.Connections {
		collection := GenerateACLCollectionFromConnection(s, &s.Connections[c], aclSelector)
		collections = append(collections, collection)
	}
	return ir.MergeACLCollections(collections...)
}

func GenerateACLCollectionFromConnection(s *ir.Spec, conn *ir.Connection, aclSelector func(target ir.IP) string) *ir.ACLCollection {
	internalSrc := conn.Src.Type != ir.EndpointTypeExternal
	internalDst := conn.Dst.Type != ir.EndpointTypeExternal
	internal := internalSrc && internalDst
	if !internalSrc && !internalDst {
		log.Fatalf("ACL: Both source and destination are external for connection %v", *conn)
	}
	result := ir.NewACLCollection()
	if !endpointRelevantToACL(conn.Src.Type) && !endpointRelevantToACL(conn.Dst.Type) {
		return result
	}
	var connectionRules []*ir.ACLRule
	for _, src := range conn.Src.Values {
		for _, dst := range conn.Dst.Values {
			if src == dst && conn.Src.Type != ir.EndpointTypeCidr && conn.Dst.Type != ir.EndpointTypeCidr {
				continue
			}
			for _, trackedProtocol := range conn.TrackedProtocols {
				reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}
				protocolRules := allowDirectedConnection(s, src, dst, conn.Src, conn.Dst, internalSrc, internalDst, trackedProtocol.Protocol, reason)
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

func allowDirectedConnection(s *ir.Spec, src, dst ir.IP, srcEp, dstEp ir.Endpoint, internalSrc, internalDst bool,
	protocol ir.Protocol, reason explanation) []*ir.ACLRule {
	var request, response *ir.Packet

	srcList := endPointsContainedInCidr(s, src, srcEp)
	dstList := endPointsContainedInCidr(s, dst, dstEp)

	var connection []*ir.ACLRule

	if internalSrc {
		for _, srcIP := range srcList {
			if srcIP == dst {
				continue
			}
			request = &ir.Packet{Src: srcIP, Dst: dst, Protocol: protocol, Explanation: reason.String()}
			connection = append(connection, ir.AllowSend(*request))
			if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
				response = &ir.Packet{Src: dst, Dst: srcIP, Protocol: inverseProtocol, Explanation: reason.response().String()}
				connection = append(connection, ir.AllowReceive(*response))
			}
		}
	}

	if internalDst {
		for _, dstIP := range dstList {
			if src == dstIP {
				continue
			}
			request = &ir.Packet{Src: src, Dst: dstIP, Protocol: protocol, Explanation: reason.String()}
			connection = append(connection, ir.AllowReceive(*request))
			if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
				response = &ir.Packet{Src: dstIP, Dst: src, Protocol: inverseProtocol, Explanation: reason.response().String()}
				connection = append(connection, ir.AllowSend(*response))
			}
		}
	}

	return connection
}

func endpointRelevantToACL(e ir.EndpointType) bool {
	return e == ir.EndpointTypeSubnet || e == ir.EndpointTypeSegment || e == ir.EndpointTypeCidr
}

func endPointsContainedInCidr(s *ir.Spec, epIP ir.IP, ep ir.Endpoint) []ir.IP {
	if ep.Type != ir.EndpointTypeCidr {
		return []ir.IP{epIP}
	}
	retVal := make([]ir.IP, 0)                             // list of subnet IPs contained in the cidr
	subnets := s.Defs.CidrSegments[ep.Name][epIP.String()] // list of subnets contained in the cidr
	for _, subnet := range subnets {
		retVal = append(retVal, s.Defs.Subnets[subnet])
	}

	return retVal
}
