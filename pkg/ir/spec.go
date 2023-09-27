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
		NIFToIP map[string]IP

		// Instance is a collection of NIFs
		InstanceToNIFs map[string][]string

		// VPEs have a single IP
		VPEToIP map[string]IP
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
	EndpointTypeNIF      EndpointType = "nif"
	EndpointTypeVPE      EndpointType = "vpe"
	EndpointTypeInstance EndpointType = "instance"
	EndpointTypeAny      EndpointType = "any"
)

func lookupSingle(m map[string]IP, name string, t EndpointType) (Endpoint, error) {
	if ip, ok := m[name]; ok {
		return Endpoint{name, []IP{ip}, t}, nil
	}
	return Endpoint{}, fmt.Errorf("%v %v not found", t, name)
}

func (s *Definitions) lookupMulti(m map[string][]string, name string, elemType, containerType EndpointType) (Endpoint, error) {
	if elems, ok := m[name]; ok {
		ips := []IP{}
		for _, elemName := range elems {
			nif, err := s.Lookup(elemType, elemName)
			if err != nil {
				return Endpoint{}, fmt.Errorf("%w while looking up %v %v for instance %v", err, elemType, elemName, name)
			}
			ips = append(ips, nif.Values...)
		}
		return Endpoint{name, ips, elemType}, nil
	}
	return Endpoint{}, fmt.Errorf("container %v %v not found", containerType, name)
}

func (s *Definitions) Lookup(t EndpointType, name string) (Endpoint, error) {
	switch t {
	case EndpointTypeExternal:
		return lookupSingle(s.Externals, name, t)
	case EndpointTypeSubnet:
		return lookupSingle(s.Subnets, name, t)
	case EndpointTypeNIF:
		return lookupSingle(s.NIFToIP, name, t)
	case EndpointTypeVPE:
		return lookupSingle(s.VPEToIP, name, t)
	case EndpointTypeInstance:
		return s.lookupMulti(s.InstanceToNIFs, name, EndpointTypeNIF, EndpointTypeInstance)
	case EndpointTypeSegment:
		return s.lookupMulti(s.SubnetSegments, name, EndpointTypeSubnet, EndpointTypeSegment)
	default:
		return Endpoint{}, fmt.Errorf("invalid type %v (endpoint %v)", t, name)
	}
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

func (s *ConfigDefs) NIFFromIP(ip IP) (string, bool) {
	return inverseLookup(s.NIFToIP, ip)
}

func (s *ConfigDefs) VPEFromIP(ip IP) (string, bool) {
	return inverseLookup(s.VPEToIP, ip)
}

func (s *ConfigDefs) InstanceFromNIF(nifName string) (string, bool) {
	return inverseLookupMulti(s.InstanceToNIFs, nifName)
}

func (s *ConfigDefs) RemoteFromIP(ip IP) RemoteType {
	if nif, okNIF := s.NIFFromIP(ip); okNIF {
		if instance, okInstance := s.InstanceFromNIF(nif); okInstance {
			return SGName(instance)
		}
		return SGName(fmt.Sprintf("<unknown instance %v>", nif))
	}
	if vpe, okVPE := s.VPEFromIP(ip); okVPE {
		return SGName(vpe)
	}
	return ip
}

type Reader interface {
	ReadSpec(filename string, defs *ConfigDefs) (*Spec, error)
}
