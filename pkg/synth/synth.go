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
	for c := range s.RequiredConnections {
		conn := s.RequiredConnections[c]
		for p := range conn.AllowedProtocols {
			protocol := conn.AllowedProtocols[p]
			aclRuleMaker, bidirectional := translateProtocol(protocol)
			srcEndpoints := lookupEndpoint(s, *conn.Src)
			dstEndpoints := lookupEndpoint(s, *conn.Dst)
			for src := range srcEndpoints {
				for dst := range dstEndpoints {
					prefix := fmt.Sprintf("rule-%v-%v-%v-%v", c, p, src, dst)
					egress := acl.NewRule(aclRuleMaker, prefix+"-outbound", true, srcEndpoints[src], dstEndpoints[dst], true)
					rulePair := []*acl.Rule{egress}
					if bidirectional {
						ingress := acl.NewRule(aclRuleMaker, prefix+"-inbound", true, srcEndpoints[src], dstEndpoints[dst], false)
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
			ips := make([]string, 0)
			for i := range section.Items {
				subnet := spec.Endpoint{
					Name: section.Items[i],
					Type: spec.EndpointType(section.Type),
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
