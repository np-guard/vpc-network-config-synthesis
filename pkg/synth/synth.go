// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"
)

type ACL struct {
	Name          string
	ResourceGroup string
	Vpc           string
	Rules         []*spec.Rule
}

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *spec.Spec) spec.Collection {
	return spec.Collection{
		ACLs: map[string]spec.ACL{
			"acl1": {Rules: generateRules(s)},
		},
	}
}

func generateRules(s *spec.Spec) []*spec.Rule {
	var allowInternal []*spec.Rule
	var allowExternal []*spec.Rule
	for c, conn := range s.Connections {
		internalSrc := conn.Src.Type != spec.EndpointTypeExternal
		for i, src := range lookupEndpoint(s, *conn.Src) {
			internalDst := conn.Dst.Type != spec.EndpointTypeExternal
			for j, dst := range lookupEndpoint(s, *conn.Dst) {
				if src == dst {
					continue
				}
				if len(conn.Protocols) == 0 {
					conn.Protocols = []spec.Protocol{spec.AnyProtocol{}}
				}
				for p, protocol := range conn.Protocols {
					prefix := fmt.Sprintf("c%v,p%v,[%v->%v],src%v,dst%v", c, p, conn.Src.Name, conn.Dst.Name, i, j)

					connection := allowDirectedConnection(src, dst, internalSrc, internalDst, protocol, prefix)
					if conn.Bidirectional {
						connection = append(connection, allowDirectedConnection(dst, src, internalDst, internalSrc, protocol, prefix)...)
					}

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
func makeDenyInternal() []*spec.Rule {
	localIPs := []string{ // https://datatracker.ietf.org/doc/html/rfc1918#section-3
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	var denyInternal []*spec.Rule
	for i, anyLocalIPSrc := range localIPs {
		for j, anyLocalIPDst := range localIPs {
			prefix := fmt.Sprintf("%vx%v", i, j)
			denyInternal = append(denyInternal, []*spec.Rule{
				packetRule(packet{anyLocalIPSrc, anyLocalIPDst, spec.AnyProtocol{}, prefix}, spec.Outbound, spec.Deny),
				packetRule(packet{anyLocalIPDst, anyLocalIPSrc, spec.AnyProtocol{}, prefix}, spec.Inbound, spec.Deny),
			}...)
		}
	}
	return denyInternal
}

type packet struct {
	src, dst string
	protocol spec.Protocol
	prefix   string
}

func allowDirectedConnection(src, dst string, internalSrc, internalDst bool, protocol spec.Protocol, prefix string) []*spec.Rule {
	var request, response *packet
	request = &packet{src, dst, protocol, prefix + ",request"}
	if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
		response = &packet{dst, src, inverseProtocol, prefix + ",response"}
	}

	var connection []*spec.Rule
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

func allowSend(packet packet) *spec.Rule {
	return packetRule(packet, spec.Outbound, spec.Allow)
}

func allowReceive(packet packet) *spec.Rule {
	return packetRule(packet, spec.Inbound, spec.Allow)
}

func packetRule(packet packet, direction spec.Direction, action spec.Action) *spec.Rule {
	return &spec.Rule{
		Name:        packet.prefix + fmt.Sprintf(",%v,%v", direction, action),
		Action:      action,
		Source:      packet.src,
		Destination: packet.dst,
		Direction:   direction,
		Protocol:    packet.protocol,
	}
}

func lookupEndpoint(s *spec.Spec, endpoint spec.Endpoint) []string {
	name := endpoint.Name
	switch endpoint.Type {
	case spec.EndpointTypeExternal:
		if ip, ok := s.Externals[name]; ok {
			return []string{ip}
		}
		return []string{fmt.Sprintf("<Unknown external %v>", name)}
	case spec.EndpointTypeSubnet:
		if ip, ok := s.Subnets[name]; ok {
			return []string{ip}
		}
		return []string{fmt.Sprintf("<Unknown subnet %v>", name)}
	case spec.EndpointTypeSegment:
		segment, ok := s.SubnetSegments[endpoint.Name]
		if ok {
			var ips []string
			for _, subnetName := range segment {
				subnet := spec.Endpoint{Name: subnetName, Type: spec.EndpointTypeSubnet}
				ips = append(ips, lookupEndpoint(s, subnet)...)
			}
			return ips
		}
		return []string{fmt.Sprintf("<Unknown segment %v>", name)}
	default:
		return []string{fmt.Sprintf("<Unknown type %v (%v)>", endpoint.Type, name)}
	}
}
