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
			aclRuleMaker, bidirectional := translateProtocol(protocol)
			for src, srcIP := range lookupEndpoint(s, *conn.Src) {
				for dst, dstIP := range lookupEndpoint(s, *conn.Dst) {
					prefix := fmt.Sprintf("rule-%v-%v-%v-%v", c, p, src, dst)
					egress := acl.NewRule(aclRuleMaker, prefix+"-outbound", true, srcIP, dstIP, true)
					rulePair := []*acl.Rule{egress}
					if bidirectional {
						ingress := acl.NewRule(aclRuleMaker, prefix+"-inbound", true, srcIP, dstIP, false)
						rulePair = []*acl.Rule{egress, ingress}
					}
					rules = append(rules, rulePair...)
				}
			}
		}
	}
	return rules
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
		log.Fatalf("External not found: %v", name)
	case spec.EndpointTypeSubnet:
		if ip, ok := s.Subnets[name]; ok {
			return []string{ip}
		}
		log.Fatalf("Subnet not found: %v", name)
	case spec.EndpointTypeSection:
		if section, ok := s.Sections[endpoint.Name]; ok {
			if section.Type != spec.TypeSubnet {
				log.Fatalf("Unsupported section item type %q", section.Type)
			}
			t := spec.EndpointType(section.Type)
			var ips []string
			for _, subnetName := range section.Items {
				subnet := spec.Endpoint{
					Name: subnetName,
					Type: t,
				}
				ips = append(ips, lookupEndpoint(s, subnet)...)
			}
			return ips
		}
	case spec.EndpointTypeNif:
	case spec.EndpointTypeInstance:
	case spec.EndpointTypeVpe:
		log.Fatalf("Unsupported endpoint type: %v", endpoint.Type)
	default:
		log.Fatalf("Unknown endpoint type: %v", endpoint.Type)
	}
	return []string{}
}

func translateProtocol(protocol interface{}) (ruleMaker acl.RuleMaker, bidirectional bool) {
	switch p := protocol.(type) {
	case *spec.TcpUdp:
		portRange := acl.PortRange{MinPort: p.MinDestinationPort, MaxPort: p.MaxDestinationPort}
		bidirectional = p.Bidirectional
		switch p.Protocol {
		case spec.TcpUdpProtocolTCP:
			ruleMaker = &acl.TCP{PortRange: portRange}
		case spec.TcpUdpProtocolUDP:
			ruleMaker = &acl.UDP{PortRange: portRange}
		default:
			log.Fatalf("Impossible protocol name: %q", p.Protocol)
		}
	case *spec.Icmp:
		bidirectional = p.Bidirectional
		ruleMaker = &acl.ICMP{Code: p.Code, Type: p.Type}
	case *spec.AnyProtocol:
		bidirectional = false
		ruleMaker = nil
	case map[string]interface{}:
		log.Fatalf("JSON unparsed: %v", p)
	default:
		log.Fatalf("Impossible protocol type: %v", p)
	}
	return
}
