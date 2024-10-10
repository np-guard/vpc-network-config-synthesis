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
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type (
	ID           = string
	NamedEntity  string
	ResourceType string

	Spec struct {
		// Required connections
		Connections []*Connection

		Defs *Definitions
	}

	Connection struct {
		// Egress resource
		Src *Resource

		// Ingress resource
		Dst *Resource

		// Allowed protocols
		TrackedProtocols []*TrackedProtocol

		// Provenance information
		Origin fmt.Stringer
	}

	Resource struct {
		// Symbolic name of resource, if available
		Name *string

		// list of CIDR / Ip addresses.
		NamedAddrs []*NamedAddrs

		// Cidr list (in case of CIDR segment)
		Cidrs []*NamedAddrs

		// Type of resource
		Type *ResourceType
	}

	NamedAddrs struct {
		IPAddrs *netset.IPBlock
		Name    *string
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

		SubnetSegments map[ID]*SegmentDetails

		CidrSegments map[ID]*CidrSegmentDetails

		NifSegments map[ID]*SegmentDetails

		InstanceSegments map[ID]*SegmentDetails

		VpeSegments map[ID]*SegmentDetails

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

	SegmentDetails struct {
		Elements []ID
	}

	CidrSegmentDetails struct {
		Cidrs            *netset.IPBlock
		ContainedSubnets []ID
	}

	ExternalDetails struct {
		ExternalAddrs *netset.IPBlock
	}

	Reader interface {
		ReadSpec(filename string, defs *ConfigDefs) (*Spec, error)
	}

	Named interface {
		Name() string
	}

	NWResource interface {
		Address() *netset.IPBlock
	}
)

const (
	ResourceTypeExternal        ResourceType = "external"
	ResourceTypeCidr            ResourceType = "cidr"
	ResourceTypeSubnet          ResourceType = "subnet"
	ResourceTypeNIF             ResourceType = "nif"
	ResourceTypeVPE             ResourceType = "vpe"
	ResourceTypeInstance        ResourceType = "instance"
	ResourceTypeSubnetSegment   ResourceType = "subnetSegment"
	ResourceTypeCidrSegment     ResourceType = "cidrSegment"
	ResourceTypeNifSegment      ResourceType = "nifSegment"
	ResourceTypeInstanceSegment ResourceType = "instanceSegment"
	ResourceTypeVpeSegment      ResourceType = "vpeSegment"

	resourceNotFound  = "%v %v not found"
	containerNotFound = "container %v %v not found"
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

func lookupSingle[T NWResource](m map[ID]T, name string, t ResourceType) (*Resource, error) {
	if details, ok := m[name]; ok {
		return &Resource{
			Name:       &name,
			NamedAddrs: []*NamedAddrs{{Name: &name, IPAddrs: details.Address()}},
			Cidrs:      []*NamedAddrs{{Name: &name, IPAddrs: details.Address()}},
			Type:       utils.Ptr(ResourceTypeSubnet),
		}, nil
	}
	return nil, fmt.Errorf(resourceNotFound, name, t)
}

func (s *Definitions) lookupNifACL(name string) (*Resource, error) {
	if nifDetails, ok := s.NIFs[name]; ok {

	}
	return nil, fmt.Errorf(resourceNotFound, name, ResourceTypeNIF)
}

func (s *Definitions) lookupSubnetSegmentACL(name string) (*Resource, error) {
	if segmentDetails, ok := s.SubnetSegments[name]; ok {
		res := &Resource{Name: &name, NamedAddrs: []*NamedAddrs{}, Cidrs: []*NamedAddrs{}, Type: utils.Ptr(ResourceTypeSubnet)}
		for _, subnetName := range segmentDetails.Elements {
			subnet, err := lookupSingle(s.Subnets, subnetName, ResourceTypeSubnet)
			if err != nil {
				return nil, fmt.Errorf("%w while looking up %v %v for subnet segment %v", err, ResourceTypeSubnet, subnetName, name)
			}
			res.NamedAddrs = append(res.NamedAddrs, subnet.NamedAddrs...)
			res.Cidrs = append(res.Cidrs, subnet.Cidrs...)
		}
		return res, nil
	}
	return nil, fmt.Errorf(containerNotFound, name, ResourceTypeSubnetSegment)
}

func (s *Definitions) lookupCidrSegmentACL(name string) (*Resource, error) {
	if segmentDetails, ok := s.CidrSegments[name]; ok {
		res := &Resource{Name: &name, NamedAddrs: []*NamedAddrs{}, Cidrs: []*NamedAddrs{}, Type: utils.Ptr(ResourceTypeSubnet)}
		for _, subnetName := range segmentDetails.ContainedSubnets {
			subnet, err := lookupSingle(s.Subnets, subnetName, ResourceTypeSubnet)
			if err != nil {
				return nil, fmt.Errorf("%w while looking up %v %v for cidr segment %v", err, ResourceTypeSubnet, subnetName, name)
			}
			res.NamedAddrs = append(res.NamedAddrs, subnet.NamedAddrs...)
		}
		for _, cidr := range segmentDetails.Cidrs.SplitToCidrs() {
			res.Cidrs = append(res.Cidrs, &NamedAddrs{Name: &name, IPAddrs: cidr})
		}
		return res, nil
	}
	return nil, fmt.Errorf(containerNotFound, name, ResourceTypeSubnet)
}

func (s *Definitions) lookupInstanceACL(name string) (*Resource, error) {
	if instnaceDetails, ok := s.Instances[name]; ok {

	}
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

func (s *Definitions) LookupACL(t ResourceType, name string) (*Resource, error) {
	switch t {
	case ResourceTypeExternal:
		return lookupSingle(s.Externals, name, t)
	case ResourceTypeSubnet:
		return lookupSingle(s.Subnets, name, t)
	case ResourceTypeNIF:
		return s.lookupNifACL(s.NIFs, name, t)
	case ResourceTypeInstance:
		return s.lookupInstanceACL(name)
	case ResourceTypeVPE:
		return s.lookupVpeACL(name)
	case ResourceTypeSubnetSegment:
		return s.lookupSubnetSegmentACL(name)
	case ResourceTypeCidrSegment:
		return s.lookupCidrSegmentACL(name)
	case ResourceTypeNifSegment:
		return s.lookupNifSegmentACL(name)
	case ResourceTypeInstanceSegment:
		return s.lookupInstanceSegmentACL(name)
	case ResourceTypeVpeSegment:
		return s.lookupVpeSegmentACL(name)
	}
	return nil, nil // should not get here
}

func (s *Definitions) LookupSG(t ResourceType, name string) (*Resource, error) {
	switch t {
	case ResourceTypeExternal:
		return lookupSingle(s.Externals, name, t)
	case ResourceTypeSubnet:
		return lookupSubnetSG(s.Subnets, name, t)
	case ResourceTypeNIF:
		return s.lookupNifSG(s.NIFs, name, t)
	case ResourceTypeInstance:
		return s.lookupInstanceSG(name)
	case ResourceTypeVPE:
		return s.lookupVpeSG(name)
	case ResourceTypeSubnetSegment:
		return s.lookupSubnetSegmentSG(name)
	case ResourceTypeCidrSegment:
		return s.lookupCidrSegmentSG(name)
	case ResourceTypeNifSegment:
		return s.lookupNifSegmentSG(name)
	case ResourceTypeInstanceSegment:
		return s.lookupInstanceSegmentSG(name)
	case ResourceTypeVpeSegment:
		return s.lookupVpeSegmentSG(name)
	}
	return nil, nil // should not get here
}

func (s *ConfigDefs) SubnetsContainedInCidr(cidr *netset.IPBlock) ([]ID, error) {
	var containedSubnets []string
	for subnet, subnetDetails := range s.Subnets {
		if subnetDetails.CIDR.IsSubset(cidr) {
			containedSubnets = append(containedSubnets, subnet)
		}
	}
	sort.Strings(containedSubnets)
	return containedSubnets, nil
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
