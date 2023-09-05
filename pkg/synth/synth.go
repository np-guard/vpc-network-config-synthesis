// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

type Options struct {
	Single bool
}

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *ir.Spec, opt Options) *ir.Collection {
	if opt.Single {
		return generateCollection(s, func(target ir.IP) string { return "1" })
	}
	return generateCollection(s, func(target ir.IP) string {
		name, ok := s.SubnetNames[target]
		if !ok {
			return fmt.Sprintf("Unknown subnet %v", target)
		}
		return name
	})
}

func GenerateCollectionFromConnection(conn *ir.Connection, aclSelector func(target ir.IP) string) *ir.Collection {
	internalSrc := conn.Src.Type != ir.EndpointTypeExternal
	internalDst := conn.Dst.Type != ir.EndpointTypeExternal
	internal := internalSrc && internalDst
	if !internalSrc && !internalDst {
		log.Fatalf("Both source and destination are external for connection %v", *conn)
	}
	result := ir.NewCollection()
	var connectionRules []*ir.Rule
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

func generateCollection(s *ir.Spec, aclSelector func(target ir.IP) string) *ir.Collection {
	collections := []*ir.Collection{}
	for c := range s.Connections {
		conn := &s.Connections[c]
		collections = append(collections, GenerateCollectionFromConnection(conn, aclSelector))
	}
	return ir.MergeCollections(collections...)
}

// func redundant(rule *ir.Rule, rules []ir.Rule) bool {
// 	for i := range rules {
// 		if mustSupersede(&rules[i], rule) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func mustSupersede(main, other *ir.Rule) bool {
// 	otherExplanation := other.Explanation
// 	other.Explanation = main.Explanation
// 	res := reflect.DeepEqual(main, other)
// 	other.Explanation = otherExplanation
// 	return res
// }

func allowDirectedConnection(src, dst ir.IP, internalSrc, internalDst bool, protocol ir.Protocol, reason explanation) []*ir.Rule {
	var request, response *ir.Packet
	request = &ir.Packet{Src: src, Dst: dst, Protocol: protocol, Explanation: reason.String()}
	if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
		response = &ir.Packet{Src: dst, Dst: src, Protocol: inverseProtocol, Explanation: reason.response().String()}
	}

	var connection []*ir.Rule
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
