// Package spec describes the input-format-agnostic specification of the required connectivity
package spec

import (
	"fmt"
)

type (
	Spec struct {
		// A list of required connections
		Connections []Connection
	}

	Connection struct {
		// In unidirectional connection, this is the egress endpoint
		Src Endpoint

		// In unidirectional connection, this is the ingress endpoint
		Dst Endpoint

		// List of allowed transport-layer connections
		Protocols []Protocol
	}

	Endpoint struct {
		// Symbolic name of endpoint, if available
		Name string

		// list of CIDR / Ip addresses.
		Values []string

		// Type of endpoint
		Type EndpointType
	}

	Definitions struct {
		Subnets map[string]string

		// Segments are a way for users to create aggregations.
		SubnetSegments map[string][]string

		// Externals are a way for users to name IP addresses or ranges external to the VPC.
		Externals map[string]string
	}
)

func MakeConnection(src, dst Endpoint, protocols []Protocol, bidirectional bool) []Connection {
	out := Connection{
		Src:       src,
		Dst:       dst,
		Protocols: protocols,
	}
	if bidirectional {
		in := Connection{Src: dst, Dst: src, Protocols: protocols}
		return []Connection{out, in}
	}
	return []Connection{out}
}

type EndpointType string

const (
	EndpointTypeExternal EndpointType = "external"
	EndpointTypeSegment  EndpointType = "segment"
	EndpointTypeSubnet   EndpointType = "subnet"
)

func (s *Definitions) Lookup(name string, expectedType EndpointType) (Endpoint, error) {
	if ip, ok := s.Externals[name]; ok {
		if expectedType != EndpointTypeExternal {
			return Endpoint{}, fmt.Errorf("<%v is external, not %v>", name, expectedType)
		}
		return Endpoint{name, []string{ip}, expectedType}, nil
	} else if ip, ok := s.Subnets[name]; ok {
		if expectedType != EndpointTypeSubnet {
			return Endpoint{}, fmt.Errorf("<%v is subnet, not %v>", name, expectedType)
		}
		return Endpoint{name, []string{ip}, expectedType}, nil
	} else {
		if segment, ok := s.SubnetSegments[name]; ok {
			if expectedType != EndpointTypeSegment {
				return Endpoint{}, fmt.Errorf("<%v is segment, not %v>", name, expectedType)
			}
			var ips []string
			for _, subnetName := range segment {
				subnet, err := s.Lookup(subnetName, EndpointTypeSubnet)
				if err != nil {
					return Endpoint{}, err
				}
				ips = append(ips, subnet.Values...)
			}
			return Endpoint{name, ips, expectedType}, nil
		}
	}
	return Endpoint{}, fmt.Errorf("<%v is not a valid %v>", name, expectedType)
}

type Reader interface {
	ReadSpec(filename string, subnetMap map[string]string) (*Spec, error)
}
