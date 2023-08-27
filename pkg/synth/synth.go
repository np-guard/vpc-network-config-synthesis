// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"
	"reflect"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

type Options struct {
	Single bool
}

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *ir.Spec, opt Options) *ir.Collection {
	if opt.Single {
		return generateCollection(s, func(target string) string { return "1" })
	}
	return generateCollection(s, func(target string) string {
		name, ok := s.SubnetNames[target]
		if !ok {
			return fmt.Sprintf("Unknown subnet %v", target)
		}
		return name
	})
}

func generateCollection(s *ir.Spec, aclSelector func(target string) string) *ir.Collection {
	result := ir.Collection{ACLs: map[string]*ir.ACL{}}
	for c := range s.Connections {
		conn := &s.Connections[c]
		internalSrc := conn.Src.Type != ir.EndpointTypeExternal
		for _, src := range conn.Src.Values {
			internalDst := conn.Dst.Type != ir.EndpointTypeExternal
			if !internalSrc && !internalDst {
				log.Fatalf("Both source and destination are external for connection #%v", c)
			}
			for _, dst := range conn.Dst.Values {
				if src == dst {
					continue
				}
				for _, trackedProtocol := range conn.TrackedProtocols {
					internal := internalSrc && internalDst
					reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}
					connection := allowDirectedConnection(src, dst, internalSrc, internalDst, trackedProtocol.Protocol, reason)

					for _, rule := range connection {
						acl := result.LookupOrCreate(aclSelector(rule.Target()))
						if internal {
							if !redundant(rule, acl.Internal) {
								acl.AppendInternal(rule)
							}
						} else {
							if !redundant(rule, acl.External) {
								acl.AppendExternal(rule)
							}
						}
					}
				}
			}
		}
	}

	return &result
}

func redundant(rule *ir.Rule, rules []ir.Rule) bool {
	for i := range rules {
		if mustSupersede(&rules[i], rule) {
			return true
		}
	}
	return false
}

func mustSupersede(main, other *ir.Rule) bool {
	otherExplanation := other.Explanation
	other.Explanation = main.Explanation
	res := reflect.DeepEqual(main, other)
	other.Explanation = otherExplanation
	return res
}

func allowDirectedConnection(src, dst string, internalSrc, internalDst bool, protocol ir.Protocol, reason explanation) []*ir.Rule {
	var request, response *ir.Packet
	request = &ir.Packet{Src: src, Dst: dst, Protocol: protocol, Explanation: reason.String()}
	if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
		responseReason := reason
		responseReason.isResponse = true
		response = &ir.Packet{Src: dst, Dst: src, Protocol: inverseProtocol, Explanation: responseReason.String()}
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
