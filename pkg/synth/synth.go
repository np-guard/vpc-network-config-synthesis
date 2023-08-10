// Package synth generates NetworkACLs that collectively enable the connectivity described in a global specification.
package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"
)

// MakeACL translates Spec to a terraform resource
func MakeACL(s *spec.Spec, subnetToIP map[string]string) string {
	result := acl.Collection{
		Items: []*acl.ACL{
			{
				Name:          "acl1",
				ResourceGroup: "var.resource_group_id",
				Vpc:           "var.vpc_id",
				Rules:         makeRules(s, lookupMap(s.Externals, subnetToIP)),
			},
		},
	}
	return result.Print()
}

func makeRules(s *spec.Spec, endpointToIP func(*spec.Endpoint) string) []*acl.Rule {
	makeRulePair := func(aclRuleMaker acl.RuleMaker, bidirectional bool, src, dst *spec.Endpoint, name string) []*acl.Rule {
		srcIP := endpointToIP(src)
		dstIP := endpointToIP(dst)
		prefix := fmt.Sprintf("rule-%v-", name)
		egress := acl.NewRule(aclRuleMaker, prefix+"outbound", true, srcIP, dstIP, true)
		if bidirectional {
			ingress := acl.NewRule(aclRuleMaker, prefix+"inbound", true, srcIP, dstIP, false)
			return []*acl.Rule{egress, ingress}
		}
		return []*acl.Rule{egress}
	}

	var rules []*acl.Rule
	for i, conn := range s.RequiredConnections {
		for _, protocol := range conn.AllowedProtocols {
			aclRuleMaker, bidirectional := translateProtocol(protocol)
			rulePair := makeRulePair(aclRuleMaker, bidirectional, conn.Src, conn.Dst, fmt.Sprintf("%v", i))
			rules = append(rules, rulePair...)
		}
	}
	for i, section := range s.Sections {
		if !section.FullyConnected {
			continue
		}
		for j, src := range section.Items {
			for k, dst := range section.Items {
				if j == k {
					continue
				}
				for m, protocol := range section.FullyConnectedWithConnectionType {
					aclRuleMaker, _ := translateProtocol(protocol)
					srcEndpoint := &spec.Endpoint{Name: src, Type: spec.EndpointType(section.Type)}
					dstEndpoint := &spec.Endpoint{Name: dst, Type: spec.EndpointType(section.Type)}
					rulePair := makeRulePair(aclRuleMaker, true, srcEndpoint, dstEndpoint,
						fmt.Sprintf("fc-section%v-src%v-dst%v-prot%v", i, j, k, m))
					rules = append(rules, rulePair...)
				}
			}
		}
	}
	return rules
}

func lookupMap(externals []spec.SpecExternalsElem, subnetToIP map[string]string) func(*spec.Endpoint) string {
	externalToIP := make(map[string]string)
	for _, ext := range externals {
		externalToIP[ext.Name] = ext.Cidr
	}
	return func(endpoint *spec.Endpoint) string {
		switch endpoint.Type {
		case spec.EndpointTypeCidr:
			return endpoint.Name
		case spec.EndpointTypeExternal:
			ip, ok := externalToIP[endpoint.Name]
			if !ok {
				log.Fatalf("external not found: %v", endpoint.Name)
			}
			return ip
		case spec.EndpointTypeSubnet:
			ip, ok := subnetToIP[endpoint.Name]
			if !ok {
				log.Fatalf("subnet not found: %v", endpoint.Name)
			}
			return ip
		case spec.EndpointTypeNif:
		case spec.EndpointTypeInstance:
		case spec.EndpointTypeSection:
		case spec.EndpointTypeVpe:
			log.Fatalf("Unsupported endpoint type: %v", endpoint.Type)
		default:
			log.Fatalf("Unknown endpoint type: %v", endpoint.Type)
		}
		return ""
	}
}

func translateProtocol(protocol interface{}) (ruleMaker acl.RuleMaker, bidirectional bool) {
	switch p := protocol.(type) {
	case *spec.TcpUdp:
		portRange := acl.PortRange{MinPort: p.MinPort, MaxPort: p.MaxPort}
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
