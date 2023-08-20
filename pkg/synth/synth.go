// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"
)

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *spec.Spec) acl.Collection {
	result := acl.Collection{
		Items: []*acl.ACL{
			{
				Name:          "acl1",
				ResourceGroup: "var.resource_group_id",
				Vpc:           "var.vpc_id",
				Rules:         generateRules(s),
			},
		},
	}
	return result
}

func generateRules(s *spec.Spec) []*acl.Rule {
	var allowInternal []*acl.Rule
	var allowExternal []*acl.Rule
	for c, conn := range s.RequiredConnections {
		internalSrc := conn.Src.Type != spec.EndpointTypeExternal
		for i, src := range lookupEndpoint(s, *conn.Src) {
			internalDst := conn.Dst.Type != spec.EndpointTypeExternal
			for j, dst := range lookupEndpoint(s, *conn.Dst) {
				if src == dst {
					continue
				}
				for p := range conn.AllowedProtocols {
					prefix := fmt.Sprintf("c%v,p%v,[%v->%v],src%v,dst%v", c, p, conn.Src.Name, conn.Dst.Name, i, j)
					protocol := makeProtocol(conn.AllowedProtocols[p])

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
func makeDenyInternal() []*acl.Rule {
	localIPs := []string{ // https://datatracker.ietf.org/doc/html/rfc1918#section-3
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	var denyInternal []*acl.Rule
	for i, anyLocalIPSrc := range localIPs {
		for j, anyLocalIPDst := range localIPs {
			prefix := fmt.Sprintf("%vx%v", i, j)
			denyInternal = append(denyInternal, []*acl.Rule{
				packetRule(packet{anyLocalIPSrc, anyLocalIPDst, acl.AnyProtocol{}, prefix}, acl.Outbound, acl.Deny),
				packetRule(packet{anyLocalIPDst, anyLocalIPSrc, acl.AnyProtocol{}, prefix}, acl.Inbound, acl.Deny),
			}...)
		}
	}
	return denyInternal
}

type packet struct {
	src, dst string
	protocol acl.Protocol
	prefix   string
}

func allowDirectedConnection(src, dst string, internalSrc, internalDst bool, protocol acl.Protocol, prefix string) []*acl.Rule {
	var request, response *packet
	request = &packet{src, dst, protocol, prefix + ",request"}
	if inverseProtocol := protocol.InverseDirection(); inverseProtocol != nil {
		response = &packet{dst, src, inverseProtocol, prefix + ",response"}
	}

	var connection []*acl.Rule
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

func allowSend(packet packet) *acl.Rule {
	return packetRule(packet, acl.Outbound, acl.Allow)
}

func allowReceive(packet packet) *acl.Rule {
	return packetRule(packet, acl.Inbound, acl.Allow)
}

func packetRule(packet packet, direction acl.Direction, action acl.Action) *acl.Rule {
	return &acl.Rule{
		Name:        packet.prefix + fmt.Sprintf(",%v,%v", direction, action),
		Action:      action,
		Source:      packet.src,
		Destination: packet.dst,
		Direction:   direction,
		Protocol:    packet.protocol,
	}
}

func makeProtocol(protocol interface{}) acl.Protocol {
	switch p := protocol.(type) {
	case *spec.TcpUdp:
		pair := acl.PortRangePair{
			SrcPort: acl.PortRange{Min: acl.DefaultMinPort, Max: acl.DefaultMaxPort},
			DstPort: acl.PortRange{Min: p.MinDestinationPort, Max: p.MaxDestinationPort},
		}
		switch p.Protocol {
		case spec.TcpUdpProtocolUDP:
			return acl.UDP{PortRangePair: pair}
		case spec.TcpUdpProtocolTCP:
			return acl.TCP{PortRangePair: pair}
		}
	case *spec.Icmp:
		return acl.ICMP{Code: p.Type, Type: p.Code}
	case *spec.AnyProtocol:
		return acl.AnyProtocol{}
	default:
		log.Fatalf("Impossible protocol type: %v", p)
	}
	return nil
}

func lookupEndpoint(s *spec.Spec, endpoint spec.Endpoint) []string {
	name := endpoint.Name
	switch endpoint.Type {
	case spec.EndpointTypeCidr:
		return []string{name}
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
		if segment, ok := s.Segments[endpoint.Name]; ok {
			if segment.Type != spec.TypeSubnet {
				return []string{fmt.Sprintf("<Unsupported segment item type %v (%v)>", segment.Type, endpoint.Name)}
			}
			t := spec.EndpointType(segment.Type)
			var ips []string
			for _, subnetName := range segment.Items {
				subnet := spec.Endpoint{Name: subnetName, Type: t}
				ips = append(ips, lookupEndpoint(s, subnet)...)
			}
			return ips
		}
	case spec.EndpointTypeNif, spec.EndpointTypeInstance, spec.EndpointTypeVpe:
		return []string{fmt.Sprintf("<Unsupported %v %v>", endpoint.Type, name)}
	default:
		return []string{fmt.Sprintf("<Unknown type %v (%v)>", endpoint.Type, name)}
	}
	return []string{}
}
