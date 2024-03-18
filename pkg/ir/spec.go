// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
	"sort"

	"github.com/np-guard/models/pkg/ipblock"
)

type (
	Spec struct {
		// Required connections
		Connections []Connection

		Defs Definitions
	}

	Connection struct {
		// Egress resource
		Src Resource

		// Ingress resource
		Dst Resource

		// Allowed protocols
		TrackedProtocols []TrackedProtocol

		// Provenance information
		Origin fmt.Stringer
	}

	Resource struct {
		// Symbolic name of resource, if available
		Name string

		// list of CIDR / Ip addresses.
		Values []IP

		// Type of resource
		Type ResourceType
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

		// list of VPC's cidrs
		AddressPrefixes []CIDR
	}

	// Definitions adds to ConfigDefs the spec-specific definitions
	Definitions struct {
		ConfigDefs

		// Segments are a way for users to create aggregations.
		SubnetSegments map[string][]string

		// key = name of the segment. value = map where its key is the cidr and its value is the contained subnets
		CidrSegments map[string]map[string][]string

		// Externals are a way for users to name IP addresses or ranges external to the VPC.
		Externals map[string]IP
	}
)

type ResourceType string

const (
	ResourceTypeExternal ResourceType = "external"
	ResourceTypeSegment  ResourceType = "segment"
	ResourceTypeCidr     ResourceType = "cidr"
	ResourceTypeSubnet   ResourceType = "subnet"
	ResourceTypeNIF      ResourceType = "nif"
	ResourceTypeVPE      ResourceType = "vpe"
	ResourceTypeInstance ResourceType = "instance"
	ResourceTypeAny      ResourceType = "any"
)

func lookupSingle(m map[string]IP, name string, t ResourceType) (Resource, error) {
	if ip, ok := m[name]; ok {
		return Resource{name, []IP{ip}, t}, nil
	}
	return Resource{}, fmt.Errorf("%v %v not found", t, name)
}

func (s *Definitions) lookupMulti(m map[string][]string, name string, elemType, containerType ResourceType) (Resource, error) {
	if elems, ok := m[name]; ok {
		ips := []IP{}
		for _, elemName := range elems {
			nif, err := s.Lookup(elemType, elemName)
			if err != nil {
				return Resource{}, fmt.Errorf("%w while looking up %v %v for instance %v", err, elemType, elemName, name)
			}
			ips = append(ips, nif.Values...)
		}
		return Resource{name, ips, elemType}, nil
	}
	return Resource{}, fmt.Errorf("container %v %v not found", containerType, name)
}

func (s *Definitions) Lookup(t ResourceType, name string) (Resource, error) {
	err := fmt.Errorf("invalid type %v (resource %v)", t, name)
	switch t {
	case ResourceTypeExternal:
		return lookupSingle(s.Externals, name, t)
	case ResourceTypeSubnet:
		return lookupSingle(s.Subnets, name, t)
	case ResourceTypeCidr:
		return lookupSingle(s.Subnets, name, t)
	case ResourceTypeNIF:
		return lookupSingle(s.NIFToIP, name, t)
	case ResourceTypeVPE:
		return lookupSingle(s.VPEToIP, name, t)
	case ResourceTypeInstance:
		return s.lookupMulti(s.InstanceToNIFs, name, ResourceTypeNIF, ResourceTypeInstance)
	case ResourceTypeSegment:
		if _, ok := s.SubnetSegments[name]; ok { // subnet segment
			return s.lookupMulti(s.SubnetSegments, name, ResourceTypeSubnet, ResourceTypeSegment)
		} else if _, ok := s.CidrSegments[name]; ok { // cidr segment
			return Resource{name, cidrsAsIPs(s.CidrSegments, name), ResourceTypeCidr}, nil
		} else {
			return Resource{}, err
		}
	default:
		return Resource{}, err
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

func cidrsAsIPs(cidrSegments map[string]map[string][]string, segmentName string) []IP {
	retVal := make([]IP, 0)
	for cidr := range cidrSegments[segmentName] {
		retVal = append(retVal, IPFromString(cidr))
	}
	return retVal
}

func (s *ConfigDefs) SubnetsContainedInCidr(cidr ipblock.IPBlock) ([]string, error) {
	var containedSubnets []string
	for subnet, ip := range s.Subnets {
		subnetIPBlock, err := ipblock.FromCidr(ip.String())
		if err != nil {
			return nil, err
		}
		if subnetIPBlock.ContainedIn(&cidr) {
			containedSubnets = append(containedSubnets, subnet)
		}
	}
	sort.Strings(containedSubnets)
	return containedSubnets, nil
}
