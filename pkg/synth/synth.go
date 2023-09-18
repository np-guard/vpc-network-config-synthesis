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
	aclSelector := func(target ir.IP) string { return "1" }
	if !opt.SingleACL {
		aclSelector = s.Defs.SubnetNameFromIP
	}
	collections := []*ir.ACLCollection{}
	for c := range s.Connections {
		collection := GenerateACLCollectionFromConnection(&s.Connections[c], aclSelector)
		collections = append(collections, collection)
	}
	return ir.MergeACLCollections(collections...)
}

// MakeSG translates Spec to a collection of security groups
func MakeSG(s *ir.Spec, opt Options) *ir.SGCollection {
	collections := []*ir.SGCollection{}
	for c := range s.Connections {
		collection := GenerateSGCollectionFromConnection(&s.Connections[c], s.Defs.RemoteFromIP)
		collections = append(collections, collection)
	}
	return ir.MergeSGCollections(collections...)
}

func GenerateSGCollectionFromConnection(conn *ir.Connection, sgSelector func(target ir.IP) ir.RemoteType) *ir.SGCollection {
	internalSrc := conn.Src.Type != ir.EndpointTypeExternal
	internalDst := conn.Dst.Type != ir.EndpointTypeExternal
	if !internalSrc && !internalDst {
		log.Fatalf("SG: Both source and destination are external for connection %v", *conn)
	}

	result := ir.NewSGCollection()

	if !sgTrigger(conn.Src.Type) && !sgTrigger(conn.Dst.Type) {
		return result
	}

	for _, src := range conn.Src.Values {
		for _, dst := range conn.Dst.Values {
			if src == dst {
				continue
			}

			for _, trackedProtocol := range conn.TrackedProtocols {
				reason := explanation{
					internal:         internalSrc && internalDst,
					connectionOrigin: conn.Origin,
					protocolOrigin:   trackedProtocol.Origin,
				}.String()

				if internalSrc {
					sgSrcName, ok := sgSelector(src).(ir.SGName)
					if !ok {
						log.Panicf("Source is not security group name: %v", src)
					}
					sgSrc := result.LookupOrCreate(sgSrcName)
					sgSrc.Attached = []ir.SGName{sgSrcName}
					sgSrc.Rules = append(sgSrc.Rules, ir.SGRule{
						Remote:      sgSelector(dst),
						Direction:   ir.Outbound,
						Protocol:    trackedProtocol.Protocol,
						Explanation: reason,
					})
				}
				if internalDst {
					sgDstName, ok := sgSelector(dst).(ir.SGName)
					if !ok {
						log.Panicf("Dst is not security group name: %v", dst)
					}
					sgDst := result.LookupOrCreate(sgDstName)
					sgDst.Attached = []ir.SGName{sgDstName}
					sgDst.Rules = append(sgDst.Rules, ir.SGRule{
						Remote:      sgSelector(src),
						Direction:   ir.Inbound,
						Protocol:    trackedProtocol.Protocol.InverseDirection(),
						Explanation: reason,
					})
				}
			}
		}
	}

	return result
}

func GenerateACLCollectionFromConnection(conn *ir.Connection, aclSelector func(target ir.IP) string) *ir.ACLCollection {
	internalSrc := conn.Src.Type != ir.EndpointTypeExternal
	internalDst := conn.Dst.Type != ir.EndpointTypeExternal
	internal := internalSrc && internalDst
	if !internalSrc && !internalDst {
		log.Fatalf("ACL: Both source and destination are external for connection %v", *conn)
	}
	result := ir.NewACLCollection()
	if !aclTrigger(conn.Src.Type) && !aclTrigger(conn.Dst.Type) {
		return result
	}
	var connectionRules []*ir.ACLRule
	for _, src := range conn.Src.Values {
		for _, dst := range conn.Dst.Values {
			if src == dst {
				continue
			}
			for _, trackedProtocol := range conn.TrackedProtocols {
				reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}
				protocolRules := allowDirectedConnection(src, dst, internalSrc, internalDst, trackedProtocol.Protocol, reason)
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

func allowDirectedConnection(src, dst ir.IP, internalSrc, internalDst bool, protocol ir.Protocol, reason explanation) []*ir.ACLRule {
	var request, response *ir.Packet
	request = &ir.Packet{Src: src, Dst: dst, Protocol: protocol, Explanation: reason.String()}
	if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
		response = &ir.Packet{Src: dst, Dst: src, Protocol: inverseProtocol, Explanation: reason.response().String()}
	}

	var connection []*ir.ACLRule
	if internalSrc {
		connection = append(connection, ir.AllowSend(*request))
		if response != nil {
			connection = append(connection, ir.AllowReceive(*response))
		}
	}
	if internalDst {
		connection = append(connection, ir.AllowReceive(*request))
		if response != nil {
			connection = append(connection, ir.AllowSend(*response))
		}
	}
	return connection
}

func aclTrigger(e ir.EndpointType) bool {
	return e == ir.EndpointTypeSubnet || e == ir.EndpointTypeSegment
}

func sgTrigger(e ir.EndpointType) bool {
	return e == ir.EndpointTypeNif || e == ir.EndpointTypeInstance
}
