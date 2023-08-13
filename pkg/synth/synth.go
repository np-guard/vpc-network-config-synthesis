// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"
)

// MakeACL translates Spec to a collection of ACLs
func MakeACL(s *spec.Spec, subnetToIP map[string]string) string {
	result := acl.Collection{
		Items: []*acl.ACL{
			{
				Name:          "acl1",
				ResourceGroup: "var.resource_group_id",
				Vpc:           "var.vpc_id",
				Rules:         makeRules(s, lookupMap(s, subnetToIP)),
			},
		},
	}
	return result.Print()
}

func makeRules(s *spec.Spec, endpointToIPs func(spec.Endpoint) []string) []*acl.Rule {
	var rules []*acl.Rule
	for i, conn := range s.RequiredConnections {
		for m, protocol := range conn.AllowedProtocols {
			aclRuleMaker, bidirectional := translateProtocol(protocol)
			for j, srcIP := range endpointToIPs(*conn.Src) {
				for k, dstIP := range endpointToIPs(*conn.Dst) {
					prefix := fmt.Sprintf("rule-%v-%v-%v-%v", i, m, j, k)
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

func lookupMap(s *spec.Spec, subnetToIP map[string]string) func(spec.Endpoint) []string {
	return func(endpoint spec.Endpoint) []string {
		return lookup(endpoint, s, subnetToIP)
	}
}

func newSection(name string) *spec.Endpoint {
	return &spec.Endpoint{Name: name, Type: spec.EndpointTypeSection}
}

func lookup(endpoint spec.Endpoint, s *spec.Spec, subnetToIP map[string]string) []string {
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
		if ip, ok := subnetToIP[name]; ok {
			return []string{ip}
		}
		log.Fatalf("Subnet not found: %v", name)
	case spec.EndpointTypeSection:
		if section, ok := s.Sections[endpoint.Name]; ok {
			var ips []string
			for _, item := range section.Items {
				ips = append(ips, lookup(*newSection(item), s, subnetToIP)...)
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
	case nil:
		bidirectional = false
		ruleMaker = nil
	case map[string]interface{}:
		log.Fatalf("JSON unparsed: %v", p)
	default:
		log.Fatalf("Impossible protocol type: %v", p)
	}
	return
}
