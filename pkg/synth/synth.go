// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"
)

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *spec.Spec) string {
	result := acl.Collection{
		Items: []*acl.ACL{
			{
				Name:          "acl1",
				ResourceGroup: "var.resource_group_id",
				Vpc:           "var.vpc_id",
				Rules:         generateConstraints(s),
			},
		},
	}
	return result.Print()
}

func generateConstraints(s *spec.Spec) []*acl.Rule {
	var result []*acl.Rule
	for c, conn := range s.RequiredConnections {
		for src, srcIP := range lookupEndpoint(s, *conn.Src) {
			for dst, dstIP := range lookupEndpoint(s, *conn.Dst) {
				if srcIP == dstIP {
					continue
				}
				for p, protocol := range conn.AllowedProtocols {
					originPrefix := fmt.Sprintf("c%v,p%v,[%v->%v],s%v,d%v", c, p, conn.Src.Name, conn.Dst.Name, src, dst)
					flows := makeFlows(originPrefix, srcIP, dstIP, protocol.(spec.ProtocolInfo).Bidi())
					protocols := makeProtocols(protocol)
					for i := range protocols {
						result = append(result, protocols[i].Rule(flows[i]))
					}
				}
			}
		}
	}
	return result
}

func makeFlows(originPrefix, srcIP, dstIP string, bidirectional bool) []acl.Flow {
	flows := []acl.Flow{
		{Name: originPrefix + "-outbound-src", Outbound: true, Source: srcIP, Destination: dstIP, Allow: true},
		{Name: originPrefix + "-inbound-src", Outbound: false, Source: srcIP, Destination: dstIP, Allow: true},
	}
	if bidirectional {
		flows = append(flows, []acl.Flow{
			{Name: originPrefix + "-outbound-dst", Outbound: true, Source: dstIP, Destination: srcIP, Allow: true},
			{Name: originPrefix + "-inbound-dst", Outbound: false, Source: dstIP, Destination: srcIP, Allow: true},
		}...)
	}
	return flows
}

func makeProtocols(protocol interface{}) []acl.Protocol {
	var protocols []acl.Protocol
	switch p := protocol.(type) {
	case *spec.TcpUdp:
		pair := acl.PortRangePair{
			SrcPort: acl.PortRange{Min: acl.DefaultMinPort, Max: acl.DefaultMaxPort},
			DstPort: acl.PortRange{Min: p.MinDestinationPort, Max: p.MaxDestinationPort},
		}
		pairs := []acl.PortRangePair{pair, pair}
		if p.Bidirectional {
			pair = acl.Swap(pair)
			pairs = append(pairs, []acl.PortRangePair{pair, pair}...)
		}
		switch p.Protocol {
		case spec.TcpUdpProtocolUDP:
			for f := range pairs {
				protocols = append(protocols, acl.UDP{PortRangePair: pairs[f]})
			}
		case spec.TcpUdpProtocolTCP:
			for f := range pairs {
				protocols = append(protocols, acl.TCP{PortRangePair: pairs[f]})
			}
		}
	case *spec.Icmp:
		aclProtocol := acl.ICMP{Code: p.Type, Type: p.Code}
		direction := []acl.Protocol{aclProtocol, aclProtocol}
		protocols = append(protocols, direction...)
		if p.Bidi() {
			protocols = append(protocols, direction...)
		}
	case *spec.AnyProtocol:
		direction := []acl.Protocol{acl.AnyProtocol{}, acl.AnyProtocol{}}
		protocols = append(protocols, direction...)
		if p.Bidi() {
			protocols = append(protocols, direction...)
		}
	default:
		log.Fatalf("Impossible protocol type: %v", p)
	}
	return protocols
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
