// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/ipblock"
)

type (
	NamedEntity string
	ID          string

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

		IPs []IP

		CIDRs []CIDR

		// Type of resource
		Type ResourceType
	}

	TrackedProtocol struct {
		Protocol
		Origin fmt.Stringer
	}

	// ConfigDefs holds definitions that are part of the network architecture
	ConfigDefs struct {
		VPCs map[ID]VPCDetails

		Subnets map[ID]SubnetDetails

		NIFs map[ID]NifDetails

		Instances map[ID]InstanceDetails

		VPEEndpoints map[ID]VPEEndpointDetails

		VPEs map[ID]VPEDetails
	}

	// Definitions adds to ConfigDefs the spec-specific definitions
	Definitions struct {
		ConfigDefs

		// Segments are a way for users to create aggregations.
		SubnetSegments map[string][]ID

		// Cidr segment might contain some cidrs
		CidrSegments map[string]map[CIDR]CIDRDetails

		// Externals are a way for users to name IP addresses or ranges external to the VPC.
		Externals map[string]IP
	}

	VPCDetails struct {
		AddressPrefixes []CIDR
		// tg
		// lb
	}

	SubnetDetails struct {
		NamedEntity
		CIDR CIDR
		VPC  ID
	}

	NifDetails struct {
		NamedEntity
		IP
		Instance ID
	}

	InstanceDetails struct {
		NamedEntity
		VPC  ID
		Nifs []ID
	}

	VPEEndpointDetails struct {
		NamedEntity
		IP
		VPEName ID
		Subnet  ID
	}

	VPEDetails struct {
		NamedEntity
		ReservedIPs []ID
		VPC         ID
	}

	CIDRDetails struct {
		ContainedSubnets []ID
		OverlappingVPCs  []ID
	}

	Named interface {
		Name() string
	}

	Endpoint interface {
		getIP() IP
	}
)

func (n NamedEntity) Name() string {
	return string(n)
}

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

func lookupExternal(m map[string]IP, name string) (Resource, error) {
	if ip, ok := m[name]; ok {
		return Resource{name, []IP{ip}, []CIDR{}, ResourceTypeExternal}, nil
	}
	return Resource{}, fmt.Errorf("%v %v not found", ResourceTypeExternal, name)
}

func lookupEndpoint[T Endpoint](m map[ID]T, name string, t ResourceType) (Resource, error) {
	if details, ok := m[ID(name)]; ok {
		return Resource{name, []IP{details.getIP()}, []CIDR{}, t}, nil
	}
	return Resource{}, fmt.Errorf("%v %v not found", t, name)
}

func lookupSubnet(m map[ID]SubnetDetails, name string) (Resource, error) {
	if details, ok := m[ID(name)]; ok {
		return Resource{name, []IP{}, []CIDR{details.CIDR}, ResourceTypeSubnet}, nil
	}
	return Resource{}, fmt.Errorf("%v %v not found", ResourceTypeExternal, name)
}

func (s *Definitions) lookupInstance(name string) (Resource, error) {
	if details, ok := s.Instances[ID(name)]; ok {
		ips := []IP{}
		for _, elemName := range details.Nifs {
			nif, err := lookupEndpoint(s.NIFs, string(elemName), ResourceTypeNIF)
			if err != nil {
				return Resource{}, fmt.Errorf("%w while looking up %v %v for instance %v", err, ResourceTypeNIF, elemName, name)
			}
			ips = append(ips, nif.IPs...)
		}
		return Resource{name, ips, []CIDR{}, ResourceTypeNIF}, nil
	}
	return Resource{}, fmt.Errorf("container %v %v not found", ResourceTypeInstance, name)
}

func (s *Definitions) lookupSubnetSegment(name string) (Resource, error) {
	cidrs := make([]CIDR, 0)
	for _, subnet := range s.SubnetSegments[name] {
		cidrs = append(cidrs, s.Subnets[subnet].CIDR)
	}
	return Resource{name, []IP{}, cidrs, ResourceTypeSubnet}, nil
}

func (s *Definitions) lookupCidrSegment(name string) (Resource, error) {
	cidrs := make([]CIDR, 0)
	for _, segment := range s.CidrSegments {
		for _, cidrDetails := range segment {
			for _, subnet := range cidrDetails.ContainedSubnets {
				cidrs = append(cidrs, s.Subnets[subnet].CIDR)
			}
		}
	}
	return Resource{name, []IP{}, cidrs, ResourceTypeCidr}, nil
}

func (s *Definitions) Lookup(t ResourceType, name string) (Resource, error) {
	err := fmt.Errorf("invalid type %v (resource %v)", t, name)
	switch t {
	case ResourceTypeExternal:
		return lookupExternal(s.Externals, name)
	case ResourceTypeSubnet:
		return lookupSubnet(s.Subnets, name)
	case ResourceTypeNIF:
		return lookupEndpoint(s.NIFs, name, t)
	case ResourceTypeVPE:
		return lookupEndpoint(s.VPEEndpoints, name, t)
	case ResourceTypeInstance:
		return s.lookupInstance(name)
	case ResourceTypeSegment:
		if _, ok := s.SubnetSegments[name]; ok {
			return s.lookupSubnetSegment(name)
		} else if _, ok := s.CidrSegments[name]; ok {
			return s.lookupCidrSegment(name)
		} else {
			return Resource{}, err
		}
	default:
		return Resource{}, err
	}
}

func inverseLookup[T Endpoint](m map[ID]T, ip IP) (result string, ok bool) {
	for id, details := range m {
		if details.getIP() == ip {
			return string(id), true
		}
	}
	return "", false
}

func inverseLookupSubnet(m map[ID]SubnetDetails, cidr CIDR) (result string, ok bool) {
	for id, details := range m {
		if details.CIDR == cidr {
			return string(id), true
		}
	}
	return "", false
}

func inverseLookupInstance(m map[ID]InstanceDetails, nifName string) (result string, ok bool) {
	for instanceName, instanceDetails := range m {
		for _, nif := range instanceDetails.Nifs {
			if nif == ID(nifName) {
				return string(instanceName), true
			}
		}
	}
	return "", false
}

func (s *ConfigDefs) SubnetNameFromCidr(cidr CIDR) (string, bool) {
	return inverseLookupSubnet(s.Subnets, cidr)
}

func (s *ConfigDefs) NIFFromIP(ip IP) (string, bool) {
	return inverseLookup(s.NIFs, ip)
}

func (s *ConfigDefs) VPEFromIP(ip IP) (string, bool) {
	return inverseLookup(s.VPEEndpoints, ip)
}

func (s *ConfigDefs) InstanceFromNIF(nifName string) (string, bool) {
	return inverseLookupInstance(s.Instances, nifName)
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

func (s *ConfigDefs) SubnetsContainedInCidr(cidr ipblock.IPBlock) ([]ID, error) {
	var containedSubnets []string
	for subnetName, subnetDetails := range s.Subnets {
		subnetsCidr := subnetDetails.CIDR
		subnetIPBlock, err := ipblock.FromCidr(subnetsCidr.String())
		if err != nil {
			return nil, err
		}
		if subnetIPBlock.ContainedIn(&cidr) {
			containedSubnets = append(containedSubnets, string(subnetName))
		}
	}
	sort.Strings(containedSubnets)

	result := make([]ID, len(containedSubnets))
	for i := range containedSubnets {
		result[i] = ID(containedSubnets[i])
	}

	return result, nil
}

func ScopingComponents(s string) []string {
	return strings.Split(s, "/")
}
