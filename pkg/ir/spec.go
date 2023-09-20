// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
)

type (
	Spec struct {
		// Required connections
		Connections []Connection

		Defs Definitions
	}

	Connection struct {
		// Egress endpoint
		Src Endpoint

		// Ingress endpoint
		Dst Endpoint

		// Allowed protocols
		TrackedProtocols []TrackedProtocol

		// Provenance information
		Origin fmt.Stringer
	}

	Endpoint struct {
		// Symbolic name of endpoint, if available
		Name string

		// list of CIDR / Ip addresses.
		Values []IP

		// Type of endpoint
		Type EndpointType
	}

	TrackedProtocol struct {
		Protocol
		Origin fmt.Stringer
	}

	// ConfigDefs holds definitions that are part of the network architecture
	ConfigDefs struct {
		Subnets map[string]IP

		// Network interface name to IP address
		NifToIP map[string]IP

		// Instance is a collection of NIFs
		InstanceToNifs map[string][]string
	}

	// Definitions adds to ConfigDefs the spec-specific definitions
	Definitions struct {
		ConfigDefs

		// Segments are a way for users to create aggregations.
		SubnetSegments map[string][]string

		// Externals are a way for users to name IP addresses or ranges external to the VPC.
		Externals map[string]IP
	}
)

type EndpointType string

const (
	EndpointTypeExternal EndpointType = "external"
	EndpointTypeSegment  EndpointType = "segment"
	EndpointTypeSubnet   EndpointType = "subnet"
	EndpointTypeNif      EndpointType = "nif"
	EndpointTypeInstance EndpointType = "instance"
	EndpointTypeAny      EndpointType = "any"
)

func (s *Definitions) Lookup(t EndpointType, name string) (Endpoint, error) {
	switch t {
	case EndpointTypeExternal:
		if ip, ok := s.Externals[name]; ok {
			return Endpoint{name, []IP{ip}, t}, nil
		}
	case EndpointTypeSubnet:
		if ip, ok := s.Subnets[name]; ok {
			return Endpoint{name, []IP{ip}, t}, nil
		}
	case EndpointTypeNif:
		if ip, ok := s.NifToIP[name]; ok {
			return Endpoint{name, []IP{ip}, t}, nil
		}
	case EndpointTypeInstance:
		if nifs, ok := s.InstanceToNifs[name]; ok {
			ips := []IP{}
			for _, nifName := range nifs {
				nif, err := s.Lookup(EndpointTypeNif, nifName)
				if err != nil {
					return Endpoint{}, fmt.Errorf("%w while looking up nif %v for instance %v", err, nifName, name)
				}
				ips = append(ips, nif.Values...)
			}
			return Endpoint{name, ips, EndpointTypeNif}, nil
		}
	case EndpointTypeSegment:
		if segment, ok := s.SubnetSegments[name]; ok {
			var ips []IP
			for _, subnetName := range segment {
				subnet, err := s.Lookup(EndpointTypeSubnet, subnetName)
				if err != nil {
					return Endpoint{}, fmt.Errorf("%w while looking up subnet %v for segment %v", err, subnet, name)
				}
				ips = append(ips, subnet.Values...)
			}
			return Endpoint{name, ips, EndpointTypeSubnet}, nil
		}
	default:
		return Endpoint{}, fmt.Errorf("invalid type %v (endpoint %v)", t, name)
	}
	return Endpoint{}, fmt.Errorf("%v %v not found", t, name)
}

func inverseLookup[K, V comparable](m map[K]V, x V) (result K, ok bool) {
	for k, v := range m {
		if v == x {
			return k, true
		}
	}
	return
}

func inverseLookupMulti[K, V comparable](m map[K][]V, x V) (result K, ok bool) {
	for k, vs := range m {
		for _, v := range vs {
			if v == x {
				return k, true
			}
		}
	}
	return
}

func (s *ConfigDefs) SubnetNameFromIP(ip IP) (string, bool) {
	return inverseLookup(s.Subnets, ip)
}

func (s *ConfigDefs) NifFromIP(ip IP) (string, bool) {
	return inverseLookup(s.NifToIP, ip)
}

func (s *ConfigDefs) InstanceFromNif(nifName string) (string, bool) {
	return inverseLookupMulti(s.InstanceToNifs, nifName)
}

func (s *ConfigDefs) RemoteFromIP(ip IP) RemoteType {
	nif, ok := s.NifFromIP(ip)
	if !ok {
		return ip
	}
	instance, ok := s.InstanceFromNif(nif)
	if !ok {
		return SGName(fmt.Sprintf("<unknown instance %v>", nif))
	}
	return SGName(instance)
}

type Reader interface {
	ReadSpec(filename string, defs *ConfigDefs) (*Spec, error)
}
