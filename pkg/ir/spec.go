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

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

type (
	ID          = string
	NamedEntity string

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
		IPAddrs []*netset.IPBlock

		// Type of resource
		Type ResourceType
	}

	TrackedProtocol struct {
		netp.Protocol
		Origin fmt.Stringer
	}

	// ConfigDefs holds definitions that are part of the network architecture
	ConfigDefs struct {
		VPCs map[ID]*VPCDetails

		Subnets map[ID]*SubnetDetails

		NIFs map[ID]*NifDetails

		Instances map[ID]*InstanceDetails

		VPEReservedIPs map[ID]*VPEReservedIPsDetails

		VPEs map[ID]*VPEDetails
	}

	// Definitions adds to ConfigDefs the spec-specific definitions
	Definitions struct {
		ConfigDefs

		// Segments are a way for users to create aggregations.
		SubnetSegments map[ID]*SubnetSegmentDetails

		// Cidr segment might contain some cidrs
		CidrSegments map[ID]*CidrSegmentDetails

		// Externals are a way for users to name IP addresses or ranges external to the VPC.
		Externals map[ID]*ExternalDetails
	}

	VPCDetails struct {
		AddressPrefixes *netset.IPBlock
		// tg
		// lb
	}

	SubnetDetails struct {
		NamedEntity
		CIDR *netset.IPBlock
		VPC  ID
	}

	NifDetails struct {
		NamedEntity
		IP       *netset.IPBlock
		VPC      ID
		Instance ID
		Subnet   ID
	}

	InstanceDetails struct {
		NamedEntity
		VPC  ID
		Nifs []ID
	}

	VPEReservedIPsDetails struct {
		NamedEntity
		IP      *netset.IPBlock
		VPEName ID
		Subnet  ID
		VPC     ID
	}

	VPEDetails struct {
		NamedEntity
		VPEReservedIPs []ID
		VPC            ID
	}

	SubnetSegmentDetails struct {
		Subnets         []ID
		OverlappingVPCs []ID
	}

	CidrSegmentDetails struct {
		Cidrs            *netset.IPBlock
		ContainedSubnets []ID
		OverlappingVPCs  []ID
	}

	ExternalDetails struct {
		ExternalAddrs *netset.IPBlock
	}

	Named interface {
		Name() string
	}

	NWResource interface {
		Address() *netset.IPBlock
	}
)

func (n *NamedEntity) Name() string {
	return string(*n)
}

func (s *SubnetDetails) Address() *netset.IPBlock {
	return s.CIDR
}

func (n *NifDetails) Address() *netset.IPBlock {
	return n.IP
}

func (v *VPEReservedIPsDetails) Address() *netset.IPBlock {
	return v.IP
}

func (e *ExternalDetails) Address() *netset.IPBlock {
	return e.ExternalAddrs
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
	if details, ok := m[name]; ok {
		return Resource{name, []*netset.IPBlock{details.Address()}, t}, nil
	}
	return Resource{}, resourceNotFoundError(name, t)
}

func (s *Definitions) lookupInstance(name string) (Resource, error) {
	if instanceDetails, ok := s.Instances[name]; ok {
		ips := make([]*netset.IPBlock, len(instanceDetails.Nifs))
		for i, elemName := range instanceDetails.Nifs {
			nif, err := s.Lookup(ResourceTypeNIF, elemName)
			if err != nil {
				return Resource{}, fmt.Errorf("%w while looking up %v %v for instance %v", err, ResourceTypeNIF, elemName, name)
			}
			// each nif has only one IP address
			ips[i] = nif.IPAddrs[0]
		}
		return Resource{name, ips, ResourceTypeNIF}, nil
	}
	return Resource{}, containerNotFoundError(name, ResourceTypeInstance)
}

func (s *Definitions) lookupVPE(name string) (Resource, error) {
	VPEDetails, ok := s.VPEs[name]
	if !ok {
		return Resource{}, resourceNotFoundError(name, ResourceTypeVPE)
	}
	ips := make([]*netset.IPBlock, len(VPEDetails.VPEReservedIPs))
	for i, vpeEndPoint := range VPEDetails.VPEReservedIPs {
		ips[i] = s.VPEReservedIPs[vpeEndPoint].IP
	}
	return Resource{name, ips, ResourceTypeVPE}, nil
}

func (s *Definitions) lookupSubnetSegment(name string) (Resource, error) {
	if subnetSegmentDetails, ok := s.SubnetSegments[name]; ok {
		cidrs := make([]*netset.IPBlock, len(subnetSegmentDetails.Subnets))
		for i, subnetName := range subnetSegmentDetails.Subnets {
			subnet, err := s.Lookup(ResourceTypeSubnet, subnetName)
			if err != nil {
				return Resource{}, fmt.Errorf("%w while looking up %v %v for subnet %v", err, ResourceTypeSubnet, subnetName, name)
			}
			// each subnet has only one CIDR block.
			cidrs[i] = subnet.IPAddrs[0]
		}
		return Resource{name, cidrs, ResourceTypeSubnet}, nil
	}
	return Resource{}, containerNotFoundError(name, ResourceTypeSegment)
}

func (s *Definitions) lookupCidrSegment(name string) (Resource, error) {
	cidrSegmentDetails, ok := s.CidrSegments[name]
	if !ok {
		return Resource{}, containerNotFoundError(name, ResourceTypeSegment)
	}
	return Resource{name, cidrSegmentDetails.Cidrs.SplitToCidrs(), ResourceTypeCidr}, nil
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
		return s.lookupVPE(name)
	case ResourceTypeInstance:
		return s.lookupInstance(name)
	case ResourceTypeSegment:
		if _, ok := s.SubnetSegments[name]; ok { // subnet segment
			return s.lookupSubnetSegment(name)
		} else if _, ok := s.CidrSegments[name]; ok { // cidr segment
			return s.lookupCidrSegment(name)
		} else {
			return Resource{}, err
		}
	default:
		return Resource{}, err
	}
}

func inverseLookup[T NWResource](m map[ID]T, address *netset.IPBlock) (result string, ok bool) {
	for name, details := range m {
		if details.Address().Equal(address) {
			return name, true
		}
	}
	return "", false
}

func (s *ConfigDefs) NIFFromIP(ip *netset.IPBlock) (string, bool) {
	return inverseLookup(s.NIFs, ip)
}

type Reader interface {
	ReadSpec(filename string, defs *ConfigDefs) (*Spec, error)
}

func (s *ConfigDefs) SubnetsContainedInCidr(cidr netset.IPBlock) ([]ID, error) {
	var containedSubnets []string
	for subnet, subnetDetails := range s.Subnets {
		if subnetDetails.CIDR.IsSubset(&cidr) {
			containedSubnets = append(containedSubnets, subnet)
		}
	}
	sort.Strings(containedSubnets)
	return containedSubnets, nil
}

func resourceNotFoundError(name string, resource ResourceType) error {
	return fmt.Errorf("%v %v not found", resource, name)
}

func containerNotFoundError(name string, resource ResourceType) error {
	return fmt.Errorf("container %v %v not found", resource, name)
}

func UniqueIDValues(s []ID) []ID {
	result := make([]ID, 0)
	seenValues := make(map[ID]struct{})

	for _, item := range s {
		if _, seen := seenValues[item]; !seen {
			seenValues[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func ScopingComponents(s string) []string {
	return strings.Split(s, "/")
}

func VpcFromScopedResource(resource ID) ID {
	return ScopingComponents(resource)[0]
}

func ChangeScoping(s string) string {
	return strings.ReplaceAll(s, "/", "--")
}
