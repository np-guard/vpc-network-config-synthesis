/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

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
		Externals map[ID]ExternalDetails
	}

	VPCDetails struct {
		AddressPrefixes []CIDR
		// tg
		// lb
	}

	SubnetDetails struct {
		NamedEntity
		CIDR IP
		VPC  ID
	}

	NifDetails struct {
		NamedEntity
		IP       IP
		Instance ID
	}

	InstanceDetails struct {
		NamedEntity
		VPC  ID
		Nifs []ID
	}

	VPEEndpointDetails struct {
		NamedEntity
		IP      IP
		VPEName ID
		Subnet  ID
	}

	VPEDetails struct {
		NamedEntity
		VPEEndpoint []ID
		VPC         ID
	}

	CIDRDetails struct {
		ContainedSubnets []ID
		OverlappingVPCs  []ID
	}

	ExternalDetails struct {
		IP IP
	}

	Named interface {
		Name() string
	}

	NWResource interface {
		Address() IP
	}
)

func (n NamedEntity) Name() string {
	return string(n)
}

func (s SubnetDetails) Address() IP {
	return s.CIDR
}

func (n NifDetails) Address() IP {
	return n.IP
}

func (v VPEEndpointDetails) Address() IP {
	return v.IP
}

func (e ExternalDetails) Address() IP {
	return e.IP
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

func lookupSingle[T NWResource](m map[ID]T, name string, t ResourceType) (Resource, error) {
	if details, ok := m[ID(name)]; ok {
		return Resource{name, []IP{details.Address()}, t}, nil
	}
	return Resource{}, fmt.Errorf("%v %v not found", t, name)
}

func (s *Definitions) lookupInstance(name string) (Resource, error) {
	if instanceDetails, ok := s.Instances[ID(name)]; ok {
		ips := []IP{}
		for _, elemName := range instanceDetails.Nifs {
			nif, err := s.Lookup(ResourceTypeNIF, string(elemName))
			if err != nil {
				return Resource{}, fmt.Errorf("%w while looking up %v %v for instance %v", err, ResourceTypeNIF, elemName, name)
			}
			ips = append(ips, nif.Values...)
		}
		return Resource{name, ips, ResourceTypeNIF}, nil
	}
	return Resource{}, containerNotFoundError(name, ResourceTypeInstance)
}

func (s *Definitions) lookupSubnetSegment(name string) (Resource, error) {
	ips := []IP{}
	if subnets, ok := s.SubnetSegments[name]; ok {
		for _, subnetName := range subnets {
			subnet, err := s.Lookup(ResourceTypeSubnet, string(subnetName))
			if err != nil {
				return Resource{}, fmt.Errorf("%w while looking up %v %v for subnet %v", err, ResourceTypeSubnet, subnetName, name)
			}
			ips = append(ips, subnet.Values...)
		}
		return Resource{name, ips, ResourceTypeSubnet}, nil
	}
	return Resource{}, containerNotFoundError(name, ResourceTypeSegment)
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
		return lookupSingle(s.NIFs, name, t)
	case ResourceTypeVPE:
		return lookupSingle(s.VPEEndpoints, name, t)
	case ResourceTypeInstance:
		return s.lookupInstance(name)
	case ResourceTypeSegment:
		if _, ok := s.SubnetSegments[name]; ok { // subnet segment
			return s.lookupSubnetSegment(name)
		} else if _, ok := s.CidrSegments[name]; ok { // cidr segment
			return Resource{name, cidrsAsIPs(s.CidrSegments, name), ResourceTypeCidr}, nil
		} else {
			return Resource{}, err
		}
	default:
		return Resource{}, err
	}
}

func (s *Definitions) ValidateConnection(src, dst Resource) error {
	if src.Type == ResourceTypeExternal || dst.Type == ResourceTypeExternal {
		return nil
	}
	srcVPCs, err := s.getVPCs(src)
	if err != nil {
		return err
	}
	dstVPCs, err := s.getVPCs(dst)
	if err != nil {
		return err
	}
	if len(srcVPCs) != 1 || len(dstVPCs) != 1 || srcVPCs[0] != dstVPCs[0] {
		return fmt.Errorf("only connections within same vpc are supported")
	}
	return nil
}

func inverseLookup[T NWResource](m map[ID]T, ip IP) (result string, ok bool) {
	for name, details := range m {
		if details.Address() == ip {
			return string(name), true
		}
	}
	return "", false
}

func inverseLookupInstance(m map[ID]InstanceDetails, name string) (result string, ok bool) {
	for instanceName, instanceDetails := range m {
		for _, nif := range instanceDetails.Nifs {
			if string(nif) == name {
				return string(instanceName), true
			}
		}
	}
	return "", false
}

func (s *ConfigDefs) SubnetNameFromIP(ip IP) (string, bool) {
	return inverseLookup(s.Subnets, ip)
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

func cidrsAsIPs(cidrSegments map[string]map[CIDR]CIDRDetails, segmentName string) []IP {
	retVal := make([]IP, 0)
	for cidr := range cidrSegments[segmentName] {
		retVal = append(retVal, IPFromCidr(cidr))
	}
	return retVal
}

func (s *ConfigDefs) SubnetsContainedInCidr(cidr ipblock.IPBlock) ([]ID, error) {
	var containedSubnets []string
	for subnet, details := range s.Subnets {
		subnetIPBlock, err := ipblock.FromCidrOrAddress(details.CIDR.String())
		if err != nil {
			return nil, err
		}
		if subnetIPBlock.ContainedIn(&cidr) {
			containedSubnets = append(containedSubnets, string(subnet))
		}
	}
	sort.Strings(containedSubnets)
	return ConvertStringToIDSlice(containedSubnets), nil
}

func ConvertStringToIDSlice(s []string) []ID {
	result := make([]ID, len(s))
	for i, val := range s {
		result[i] = ID(val)
	}
	return result
}

func ScopingComponents(s string) []string {
	return strings.Split(s, "_")
}

func containerNotFoundError(name string, resource ResourceType) error {
	return fmt.Errorf("container %v %v not found", ResourceTypeSegment, name)
}

func (s ConfigDefs) getVPCs(r Resource) ([]ID, error) {
	vpcs := make([]ID, 0)

	// convert address to IPBlock
	if len(r.Values) == 0 {
		return vpcs, nil
	}
	addresses, err := ipblock.FromCidrOrAddress(r.Values[0].String())
	if err != nil {
		return nil, err
	}
	for _, address := range r.Values {
		currAddress, err := ipblock.FromCidrOrAddress(address.String())
		if err != nil {
			return nil, err
		}
		addresses = addresses.Union(currAddress)
	}

	for vpcName, vpcDetails := range s.VPCs {
		for _, addressPrefix := range vpcDetails.AddressPrefixes {
			iAddressPrefix, err := ipblock.FromCidrOrAddress(addressPrefix.String())
			if err != nil {
				return nil, err
			}
			if !addresses.Intersect(iAddressPrefix).IsEmpty() {
				vpcs = append(vpcs, vpcName)
				break
			}
		}
	}

	return vpcs, nil
}
