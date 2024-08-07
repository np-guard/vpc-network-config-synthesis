/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/ipblock"
)

const MaximalIPv4PrefixLength = 32

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
		IPAddrs []*ipblock.IPBlock

		// Type of resource
		Type ResourceType
	}

	TrackedProtocol struct {
		Protocol
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
		AddressPrefixes *ipblock.IPBlock
		// tg
		// lb
	}

	SubnetDetails struct {
		NamedEntity
		CIDR *ipblock.IPBlock
		VPC  ID
	}

	NifDetails struct {
		NamedEntity
		IP       *ipblock.IPBlock
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
		IP      *ipblock.IPBlock
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
		Cidrs map[*ipblock.IPBlock]CIDRDetails
	}

	CIDRDetails struct {
		ContainedSubnets []ID
		OverlappingVPCs  []ID
	}

	ExternalDetails struct {
		ExternalAddrs *ipblock.IPBlock
	}

	Named interface {
		Name() string
	}

	NWResource interface {
		Address() *ipblock.IPBlock
	}

	ResourceVpc interface {
		getOverlappingVPCs() []ID
	}
)

func (n *NamedEntity) Name() string {
	return string(*n)
}

func (s *SubnetDetails) Address() *ipblock.IPBlock {
	return s.CIDR
}

func (n *NifDetails) Address() *ipblock.IPBlock {
	return n.IP
}

func (v *VPEReservedIPsDetails) Address() *ipblock.IPBlock {
	return v.IP
}

func (e *ExternalDetails) Address() *ipblock.IPBlock {
	return e.ExternalAddrs
}

func (s *SubnetDetails) getOverlappingVPCs() []ID {
	return []ID{s.VPC}
}

func (n *NifDetails) getOverlappingVPCs() []ID {
	return []ID{n.VPC}
}

func (i *InstanceDetails) getOverlappingVPCs() []ID {
	return []ID{i.VPC}
}

func (v *VPEDetails) getOverlappingVPCs() []ID {
	return []ID{v.VPC}
}

func (s *SubnetSegmentDetails) getOverlappingVPCs() []ID {
	return s.OverlappingVPCs
}

func (c *CidrSegmentDetails) getOverlappingVPCs() []ID {
	result := make([]ID, 0)
	for _, cidrDetails := range c.Cidrs {
		result = append(result, cidrDetails.OverlappingVPCs...)
	}
	return UniqueIDValues(result)
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

func getResourceVPCs[T ResourceVpc](m map[ID]T, name string) []ID {
	return m[name].getOverlappingVPCs()
}

func lookupSingle[T NWResource](m map[ID]T, name string, t ResourceType) (Resource, error) {
	if details, ok := m[name]; ok {
		return Resource{name, []*ipblock.IPBlock{details.Address()}, t}, nil
	}
	return Resource{}, resourceNotFoundError(name, t)
}

func (s *Definitions) lookupInstance(name string) (Resource, error) {
	if instanceDetails, ok := s.Instances[name]; ok {
		ips := make([]*ipblock.IPBlock, len(instanceDetails.Nifs))
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
	ips := make([]*ipblock.IPBlock, len(VPEDetails.VPEReservedIPs))
	for i, vpeEndPoint := range VPEDetails.VPEReservedIPs {
		ips[i] = s.VPEReservedIPs[vpeEndPoint].IP
	}
	return Resource{name, ips, ResourceTypeVPE}, nil
}

func (s *Definitions) lookupSubnetSegment(name string) (Resource, error) {
	if subnetSegmentDetails, ok := s.SubnetSegments[name]; ok {
		cidrs := make([]*ipblock.IPBlock, len(subnetSegmentDetails.Subnets))
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
	cidrs := make([]*ipblock.IPBlock, len(cidrSegmentDetails.Cidrs))
	i := 0
	for cidr := range cidrSegmentDetails.Cidrs {
		cidrs[i] = cidr
		i++
	}
	return Resource{name, cidrs, ResourceTypeCidr}, nil
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

func (s *Definitions) GetResourceOverlappingVPCs(t ResourceType, name string) []ID {
	switch t {
	case ResourceTypeExternal:
		return []ID{}
	case ResourceTypeSubnet:
		return getResourceVPCs(s.Subnets, name)
	case ResourceTypeNIF:
		return getResourceVPCs(s.NIFs, name)
	case ResourceTypeVPE:
		return getResourceVPCs(s.VPEs, name)
	case ResourceTypeInstance:
		return getResourceVPCs(s.Instances, name)
	case ResourceTypeSegment:
		if _, ok := s.SubnetSegments[name]; ok { // subnet segment
			return getResourceVPCs(s.SubnetSegments, name)
		}
		return getResourceVPCs(s.CidrSegments, name)
	default:
		return []ID{}
	}
}

func (s *Definitions) ValidateConnection(srcVPCs, dstVPCs []ID) error {
	if len(srcVPCs) == 0 || len(dstVPCs) == 0 { // src or dst is an external resource
		return nil
	}
	if len(srcVPCs) != 1 || len(dstVPCs) != 1 || srcVPCs[0] != dstVPCs[0] {
		return fmt.Errorf("only connections within same vpc are supported")
	}
	return nil
}

func inverseLookup[T NWResource](m map[ID]T, address *ipblock.IPBlock) (result string, ok bool) {
	for name, details := range m {
		if details.Address().Equal(address) {
			return name, true
		}
	}
	return "", false
}

func (s *ConfigDefs) NIFFromIP(ip *ipblock.IPBlock) (string, bool) {
	return inverseLookup(s.NIFs, ip)
}

type Reader interface {
	ReadSpec(filename string, defs *ConfigDefs) (*Spec, error)
}

func (s *ConfigDefs) SubnetsContainedInCidr(cidr ipblock.IPBlock) ([]ID, error) {
	var containedSubnets []string
	for subnet, subnetDetails := range s.Subnets {
		if subnetDetails.CIDR.ContainedIn(&cidr) {
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

func IsIPAddress(address *ipblock.IPBlock) bool {
	prefixLength, err := address.PrefixLength()
	if err != nil {
		log.Fatal(err)
	}
	return prefixLength == MaximalIPv4PrefixLength
}

func ChangeScoping(s string) string {
	return strings.ReplaceAll(s, "/", "--")
}
