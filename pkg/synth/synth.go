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
	var result []*acl.Rule
	for c, conn := range s.RequiredConnections {
		for i, src := range lookupEndpoint(s, *conn.Src) {
			for j, dst := range lookupEndpoint(s, *conn.Dst) {
				if src == dst {
					continue
				}
				for p := range conn.AllowedProtocols {
					prefix := fmt.Sprintf("c%v,p%v,[%v->%v],src%v,dst%v", c, p, conn.Src.Name, conn.Dst.Name, i, j)
					protocol := conn.AllowedProtocols[p].(spec.ProtocolInfo)
					result = append(result, allowDirectedConnection(src, dst, protocol, prefix)...)
					if protocol.Bidi() {
						result = append(result, allowDirectedConnection(dst, src, protocol, prefix)...)
					}
				}
			}
		}
	}
	return result
}

type packet struct {
	src, dst string
	protocol acl.Protocol
	prefix   string
}

func allowDirectedConnection(src, dst string, protocol spec.ProtocolInfo, prefix string) []*acl.Rule {
	inout := makeProtocol(protocol)
	request := packet{src, dst, inout, prefix + ",request"}
	result := []*acl.Rule{
		allowSend(request),
		allowReceive(request),
	}
	if inverseProtocol := inout.InverseDirection(); inverseProtocol != nil {
		response := packet{dst, src, inverseProtocol, prefix + ",response"}
		result = append(result, []*acl.Rule{
			allowSend(response),
			allowReceive(response),
		}...)
	}
	return result
}

func allowSend(packet packet) *acl.Rule {
	return allowPacket(packet, acl.Outbound)
}

func allowReceive(packet packet) *acl.Rule {
	return allowPacket(packet, acl.Inbound)
}

func allowPacket(packet packet, direction acl.Direction) *acl.Rule {
	return &acl.Rule{
		Name:        packet.prefix + fmt.Sprintf(",%v", direction),
		Action:      acl.Allow,
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
	case spec.EndpointTypeSection:
		if section, ok := s.Sections[endpoint.Name]; ok {
			if section.Type != spec.TypeSubnet {
				return []string{fmt.Sprintf("<Unsupported section item type %v (%v)>", section.Type, endpoint.Name)}
			}
			t := spec.EndpointType(section.Type)
			var ips []string
			for _, subnetName := range section.Items {
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
