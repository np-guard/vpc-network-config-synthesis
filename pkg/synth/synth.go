package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/tf"
)

// MakeACL translates Spec to a terraform resource
func (s *Spec) MakeACL(subnetToIP map[string]string) string {
	result := tf.ACLCollection{
		Items: []*tf.ACL{
			{
				Name:          "acl1",
				ResourceGroup: "var.resource_group_id",
				Vpc:           "var.vpc_id",
				Rules:         s.makeRules(lookupMap(s.Externals, subnetToIP)),
			},
		},
	}
	return result.Print()
}

func (s *Spec) makeRules(endpointToIP func(*Endpoint) string) []*tf.ACLRule {
	makeRulePair := func(aclRuleMaker tf.ACLRuleMaker, bidirectional bool, src, dst *Endpoint, name string) []*tf.ACLRule {
		srcIP := endpointToIP(src)
		dstIP := endpointToIP(dst)
		prefix := fmt.Sprintf("rule-%v-", name)
		egress := tf.NewACLRule(aclRuleMaker, prefix+"outbound", true, srcIP, dstIP, true)
		if bidirectional {
			ingress := tf.NewACLRule(aclRuleMaker, prefix+"inbound", true, srcIP, dstIP, false)
			return []*tf.ACLRule{egress, ingress}
		}
		return []*tf.ACLRule{egress}
	}

	var rules []*tf.ACLRule
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
					srcEndpoint := &Endpoint{src, EndpointType(section.Type)}
					dstEndpoint := &Endpoint{dst, EndpointType(section.Type)}
					rulePair := makeRulePair(aclRuleMaker, true, srcEndpoint, dstEndpoint,
						fmt.Sprintf("fc-section%v-src%v-dst%v-prot%v", i, j, k, m))
					rules = append(rules, rulePair...)
				}
			}
		}
	}
	return rules
}

func lookupMap(externals []SpecExternalsElem, subnetToIP map[string]string) func(*Endpoint) string {
	externalToIP := make(map[string]string)
	for _, ext := range externals {
		externalToIP[ext.Name] = ext.Cidr
	}
	return func(endpoint *Endpoint) string {
		switch endpoint.Type {
		case EndpointTypeCidr:
			return endpoint.Name
		case EndpointTypeExternal:
			ip, ok := externalToIP[endpoint.Name]
			if !ok {
				log.Fatalf("external not found: %v", endpoint.Name)
			}
			return ip
		case EndpointTypeSubnet:
			ip, ok := subnetToIP[endpoint.Name]
			if !ok {
				log.Fatalf("subnet not found: %v", endpoint.Name)
			}
			return ip
		case EndpointTypeNif:
		case EndpointTypeInstance:
		case EndpointTypeSection:
		case EndpointTypeVpe:
			log.Fatalf("Unsupported endpoint type: %v", endpoint.Type)
		default:
			log.Fatalf("Unknown endpoint type: %v", endpoint.Type)
		}
		return ""
	}
}

func translateProtocol(protocol interface{}) (ruleMaker tf.ACLRuleMaker, bidirectional bool) {
	switch p := protocol.(type) {
	case *TcpUdp:
		portRange := tf.PortRange{MinPort: p.MinPort, MaxPort: p.MaxPort}
		bidirectional = p.Bidirectional
		switch p.Protocol {
		case TcpUdpProtocolTCP:
			ruleMaker = &tf.TCP{PortRange: portRange}
		case TcpUdpProtocolUDP:
			ruleMaker = &tf.UDP{PortRange: portRange}
		default:
			log.Fatalf("Impossible protocol name: %q", p.Protocol)
		}
	case *Icmp:
		bidirectional = p.Bidirectional
		ruleMaker = &tf.ICMP{Code: p.Code, Type: p.Type}
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
