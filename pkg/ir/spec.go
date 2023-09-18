// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
	"strings"
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

	ConfigDefs struct {
		Subnets map[string]IP

		// Network interface name to IP address
		NifToIP map[string]IP

		// Instance is a collection of NIFs
		InstanceToNifs map[string][]string
	}

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

//nolint:gocyclo  // Case by case basis
func (s *Definitions) Lookup(name string, expectedType EndpointType) (Endpoint, error) {
	var result []Endpoint
	if ip, ok := s.Externals[name]; ok {
		actualType := EndpointTypeExternal
		if expectedType != EndpointTypeAny && expectedType != actualType {
			return Endpoint{}, fmt.Errorf("%v is external, not %v", name, expectedType)
		}
		result = append(result, Endpoint{name, []IP{ip}, actualType})
	}
	if ip, ok := s.Subnets[name]; ok {
		actualType := EndpointTypeSubnet
		if expectedType != EndpointTypeAny && expectedType != actualType {
			return Endpoint{}, fmt.Errorf("%v is subnet, not %v", name, expectedType)
		}
		result = append(result, Endpoint{name, []IP{ip}, actualType})
	}
	if ip, ok := s.NifToIP[name]; ok {
		actualType := EndpointTypeNif
		if expectedType != EndpointTypeAny && expectedType != actualType {
			return Endpoint{}, fmt.Errorf("%v is nif, not %v", name, expectedType)
		}
		result = append(result, Endpoint{name, []IP{ip}, actualType})
	}
	if nifs, ok := s.InstanceToNifs[name]; ok {
		actualType := EndpointTypeInstance
		if expectedType != EndpointTypeAny && expectedType != actualType {
			return Endpoint{}, fmt.Errorf("%v is instance, not %v", name, expectedType)
		}
		ips := []IP{}
		for _, nif := range nifs {
			ips = append(ips, s.NifToIP[nif])
		}
		result = append(result, Endpoint{name, ips, EndpointTypeNif})
	}
	if segment, ok := s.SubnetSegments[name]; ok {
		actualType := EndpointTypeSegment
		if expectedType != EndpointTypeAny && expectedType != actualType {
			return Endpoint{}, fmt.Errorf("%v is segment, not %v", name, expectedType)
		}
		actualType = EndpointTypeSubnet
		var ips []IP
		for _, subnetName := range segment {
			subnet, err := s.Lookup(subnetName, actualType)
			if err != nil {
				return Endpoint{}, err
			}
			ips = append(ips, subnet.Values...)
		}
		result = append(result, Endpoint{name, ips, actualType})
	}
	if len(result) > 1 {
		var possibleTypes []string
		for i := range result {
			possibleTypes[i] = string(result[i].Type)
		}
		return Endpoint{}, fmt.Errorf("%v is ambiguous: may be either one of %v", name, strings.Join(possibleTypes, ", "))
	}
	if len(result) == 0 {
		return Endpoint{}, fmt.Errorf("%v is not a valid %v", name, expectedType)
	}
	return result[0], nil
}

func inverseLookup[K, V comparable](m map[K]V, x V, notFound K) K {
	for k, v := range m {
		if v == x {
			return k
		}
	}
	return notFound
}

func inverseLookupMulti[K, V comparable](m map[K][]V, x V, notFound K) K {
	for k, vs := range m {
		for _, v := range vs {
			if v == x {
				return k
			}
		}
	}
	return notFound
}

func (s *ConfigDefs) SubnetNameFromIP(ip IP) string {
	return inverseLookup(s.Subnets, ip, fmt.Sprintf("<unknown subnet %v>", ip))
}

func (s *ConfigDefs) NifFromIP(ip IP) string {
	return inverseLookup(s.NifToIP, ip, fmt.Sprintf("<unknown nif %v>", ip))
}

func (s *ConfigDefs) InstanceFromNif(nifName string) string {
	return inverseLookupMulti(s.InstanceToNifs, nifName, fmt.Sprintf("<unknown instance %v>", nifName))
}

func (s *ConfigDefs) RemoteFromIP(ip IP) RemoteType {
	if ip.String() == "0.0.0.0" {
		return ip
	}
	if ip.String() == "0.0.0.0/0" {
		return ip
	}
	return SGName(s.InstanceFromNif(s.NifFromIP(ip)))
}

type Reader interface {
	ReadSpec(filename string, defs *ConfigDefs) (*Spec, error)
}
