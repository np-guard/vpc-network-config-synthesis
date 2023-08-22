// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *ir.Spec) ir.Collection {
	return ir.Collection{
		ACLs: map[string]ir.ACL{
			"acl1": {Rules: generateRules(s)},
		},
	}
}

func generateRules(s *ir.Spec) []*ir.Rule {
	var allowInternal []*ir.Rule
	var allowExternal []*ir.Rule
	for c := range s.Connections {
		conn := &s.Connections[c]
		internalSrc := conn.Src.Type != ir.EndpointTypeExternal
		for i, src := range conn.Src.Values {
			internalDst := conn.Dst.Type != ir.EndpointTypeExternal
			if !internalSrc && !internalDst {
				log.Fatalf("Both source and destination are external for connection #%v", c)
			}
			for j, dst := range conn.Dst.Values {
				if src == dst {
					continue
				}
				if len(conn.Protocols) == 0 {
					conn.Protocols = []ir.Protocol{ir.AnyProtocol{}}
				}
				for p, protocol := range conn.Protocols {
					prefix := fmt.Sprintf("c%v,p%v,[%v->%v],src%v,dst%v", c, p, conn.Src.Name, conn.Dst.Name, i, j)

					connection := allowDirectedConnection(src, dst, internalSrc, internalDst, protocol, prefix)

					if internalSrc && internalDst {
						allowInternal = append(allowInternal, connection...)
					} else {
						allowExternal = append(allowExternal, connection...)
					}
				}
			}
		}
	}
	result := allowInternal
	if len(allowExternal) != 0 {
		result = append(result, makeDenyInternal()...)
		result = append(result, allowExternal...)
	}
	return result
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
			prefix := fmt.Sprintf("%vx%v", i, j)
			denyInternal = append(denyInternal, []*ir.Rule{
				packetRule(packet{anyLocalIPSrc, anyLocalIPDst, ir.AnyProtocol{}, prefix}, ir.Outbound, ir.Deny),
				packetRule(packet{anyLocalIPDst, anyLocalIPSrc, ir.AnyProtocol{}, prefix}, ir.Inbound, ir.Deny),
			}...)
		}
	}
	return denyInternal
}

type packet struct {
	src, dst string
	protocol ir.Protocol
	prefix   string
}

func allowDirectedConnection(src, dst string, internalSrc, internalDst bool, protocol ir.Protocol, prefix string) []*ir.Rule {
	var request, response *packet
	request = &packet{src, dst, protocol, prefix + ",request"}
	if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
		response = &packet{dst, src, inverseProtocol, prefix + ",response"}
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
		Name:        packet.prefix + fmt.Sprintf(",%v,%v", direction, action),
		Action:      action,
		Source:      packet.src,
		Destination: packet.dst,
		Direction:   direction,
		Protocol:    packet.protocol,
	}
}
