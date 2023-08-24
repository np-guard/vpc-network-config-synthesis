// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"
	"reflect"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *ir.Spec) *ir.Collection {
	return &ir.Collection{
		ACLs: map[string]ir.ACL{
			"acl1": {Rules: generateRules(s)},
		},
	}
}

func generateRules(s *ir.Spec) []ir.Rule {
	var allowInternal []*ir.Rule
	var allowExternal []*ir.Rule
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
				if len(conn.Protocols) == 0 {
					conn.Protocols = []ir.Protocol{ir.AnyProtocol{}}
				}
				for p, protocol := range conn.Protocols {
					internal := internalSrc && internalDst
					explanation := ir.Explanation{ProtocolIndex: p, Internal: internal, Origin: conn.Reason}
					connection := allowDirectedConnection(src, dst, internalSrc, internalDst, protocol, explanation)

					if internal {
						allowInternal = append(allowInternal, connection...)
					} else {
						allowExternal = append(allowExternal, connection...)
					}
				}
			}
		}
	}

	rules := allowInternal
	if len(allowExternal) != 0 {
		rules = append(rules, makeDenyInternal()...)
		rules = append(rules, allowExternal...)
	}

	nullifyRedundant(rules)
	return copyNonNil(rules)
}

func copyNonNil(list []*ir.Rule) []ir.Rule {
	result := make([]ir.Rule, countNonNil(list))
	i := 0
	for _, maybeRule := range list {
		if maybeRule != nil {
			result[i] = *maybeRule
			i++
		}
	}
	return result
}

func countNonNil[T any](list []*T) int {
	result := 0
	for i := range list {
		if list[i] != nil {
			result++
		}
	}
	return result
}

func nullifyRedundant(rules []*ir.Rule) {
	for i, main := range rules {
		if main == nil {
			continue
		}
		for j := i + 1; j < len(rules); j++ {
			other := rules[j]
			if other == nil {
				continue
			}
			if mustSupersede(main, other) {
				rules[j] = nil
			}
		}
	}
}

func mustSupersede(main, other *ir.Rule) bool {
	otherExplanation := other.Explanation
	other.Explanation = main.Explanation
	res := reflect.DeepEqual(main, other)
	other.Explanation = otherExplanation
	return res
}

// makeDenyInternal prevents allowing external communications from accidentally allowing internal communications too
func makeDenyInternal() []*ir.Rule {
	localIPs := []string{ // https://datatracker.ietf.org/doc/html/rfc1918#section-3
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	var denyInternal []*ir.Rule
	for i, anyLocalIPSrc := range localIPs {
		for j, anyLocalIPDst := range localIPs {
			explanation := fmt.Sprintf("Deny other internal communication; see rfc1918#3; item %v,%v", i, j)
			denyInternal = append(denyInternal, []*ir.Rule{
				packetRule(packet{src: anyLocalIPSrc, dst: anyLocalIPDst, protocol: ir.AnyProtocol{}, explanation: explanation}, ir.Outbound, ir.Deny),
				packetRule(packet{src: anyLocalIPDst, dst: anyLocalIPSrc, protocol: ir.AnyProtocol{}, explanation: explanation}, ir.Inbound, ir.Deny),
			}...)
		}
	}
	return denyInternal
}

type packet struct {
	src, dst    string
	protocol    ir.Protocol
	explanation string
}

func allowDirectedConnection(src, dst string, internalSrc, internalDst bool, protocol ir.Protocol, explanation ir.Explanation) []*ir.Rule {
	var request, response *packet
	request = &packet{src: src, dst: dst, protocol: protocol, explanation: explanation.String()}
	if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
		responseExplanation := explanation
		responseExplanation.IsResponse = true
		response = &packet{src: dst, dst: src, protocol: inverseProtocol, explanation: explanation.String()}
	}

	var connection []*ir.Rule
	if internalSrc {
		connection = append(connection, allowSend(*request))
		if response != nil {
			connection = append(connection, allowReceive(*response))
		}
	}
	if internalDst {
		connection = append(connection, allowReceive(*request))
		if response != nil {
			connection = append(connection, allowSend(*response))
		}
	}
	return connection
}

func allowSend(packet packet) *ir.Rule {
	return packetRule(packet, ir.Outbound, ir.Allow)
}

func allowReceive(packet packet) *ir.Rule {
	return packetRule(packet, ir.Inbound, ir.Allow)
}

func packetRule(packet packet, direction ir.Direction, action ir.Action) *ir.Rule {
	return &ir.Rule{
		Action:      action,
		Source:      packet.src,
		Destination: packet.dst,
		Direction:   direction,
		Protocol:    packet.protocol,
		Explanation: packet.explanation,
	}
}
