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
		Values []string

		// Type of endpoint
		Type EndpointType
	}

	TrackedProtocol struct {
		Protocol
		Origin fmt.Stringer
	}

	Definitions struct {
		Subnets map[string]string

		// Segments are a way for users to create aggregations.
		SubnetSegments map[string][]string

		// Externals are a way for users to name IP addresses or ranges external to the VPC.
		Externals map[string]string
	}
)

type EndpointType string

const (
	EndpointTypeExternal EndpointType = "external"
	EndpointTypeSegment  EndpointType = "segment"
	EndpointTypeSubnet   EndpointType = "subnet"
	EndpointTypeAny      EndpointType = "any"
)

func (s *Definitions) Lookup(name string, expectedType EndpointType) (Endpoint, error) {
	var result []Endpoint
	if ip, ok := s.Externals[name]; ok {
		actualType := EndpointTypeExternal
		if expectedType != EndpointTypeAny && expectedType != actualType {
			return Endpoint{}, fmt.Errorf("%v is external, not %v", name, expectedType)
		}
		result = append(result, Endpoint{name, []string{ip}, actualType})
	}
	if ip, ok := s.Subnets[name]; ok {
		actualType := EndpointTypeSubnet
		if expectedType != EndpointTypeAny && expectedType != EndpointTypeSubnet {
			return Endpoint{}, fmt.Errorf("%v is subnet, not %v", name, expectedType)
		}
		result = append(result, Endpoint{name, []string{ip}, actualType})
	}
	if segment, ok := s.SubnetSegments[name]; ok {
		actualType := EndpointTypeSegment
		if expectedType != EndpointTypeAny && expectedType != actualType {
			return Endpoint{}, fmt.Errorf("%v is segment, not %v", name, expectedType)
		}
		actualType = EndpointTypeSubnet
		var ips []string
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

type Reader interface {
	ReadSpec(filename string, subnetMap map[string]string) (*Spec, error)
}
