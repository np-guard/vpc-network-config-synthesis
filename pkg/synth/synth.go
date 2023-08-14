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
				Rules:         makeRules(s),
			},
		},
	}
	return result.Print()
}

func makeRules(s *spec.Spec) []*acl.Rule {
	var rules []*acl.Rule
	for c, conn := range s.RequiredConnections {
		for p, protocol := range conn.AllowedProtocols {
			aclConnection, bidirectional := translateProtocol(protocol)
			for src, srcIP := range lookupEndpoint(s, *conn.Src) {
				for dst, dstIP := range lookupEndpoint(s, *conn.Dst) {
					if srcIP == dstIP {
						continue
					}
					prefix := fmt.Sprintf("rule:c%v,p%v,[%v->%v],s%v,d%v", c, p, conn.Src.Name, conn.Dst.Name, src, dst)
					egress := newRule(prefix+"-outbound", srcIP, dstIP, true, aclConnection)
					rulePair := []*acl.Rule{egress}
					if bidirectional {
						ingress := newRule(prefix+"-inbound", dstIP, srcIP, false, aclConnection)
						rulePair = []*acl.Rule{egress, ingress}
					}
					rules = append(rules, rulePair...)
				}
			}
		}
	}
	return rules
}

func newRule(name, srcIP, dstIP string, outbound bool, aclConnection acl.Connection) *acl.Rule {
	return &acl.Rule{Name: name, Allow: true, Source: srcIP, Destination: dstIP, Outbound: outbound, Protocol: aclConnection}
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

func translateProtocol(protocol interface{}) (aclConnection acl.Connection, bidirectional bool) {
	switch p := protocol.(type) {
	case *spec.TcpUdp:
		portRange := acl.PortRange{MinPort: p.MinDestinationPort, MaxPort: p.MaxDestinationPort}
		bidirectional = p.Bidirectional
		switch p.Protocol {
		case spec.TcpUdpProtocolTCP:
			aclConnection = &acl.TCP{PortRange: portRange}
		case spec.TcpUdpProtocolUDP:
			aclConnection = &acl.UDP{PortRange: portRange}
		default:
			log.Fatalf("Impossible protocol name: %q", p.Protocol)
		}
	case *spec.Icmp:
		bidirectional = p.Bidirectional
		aclConnection = &acl.ICMP{Code: p.Code, Type: p.Type}
	case *spec.AnyProtocol:
		bidirectional = p.Bidirectional
		aclConnection = nil
	case map[string]interface{}:
		log.Fatalf("JSON unparsed: %v", p)
	default:
		log.Fatalf("Impossible protocol type: %v", p)
	}
	return
}
