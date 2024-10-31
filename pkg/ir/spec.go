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
	ID           = string
	NamedEntity  string
	ResourceType string

	Spec struct {
		// Required connections
		Connections []*Connection

		Defs *Definitions

		*BlockedResources
	}

	Connection struct {
		// Egress resource
		Src *LocalRemotePair

		// Ingress resource
		Dst *LocalRemotePair

		// Allowed protocols
		TrackedProtocols []*TrackedProtocol

		// Provenance information
		Origin fmt.Stringer
	}

	// LocalRemotePair holds a local resource and the remote CIDRs it should be connected to
	LocalRemotePair struct {
		// Symbolic name of resource, if available
		Name *string

		LocalCidrs []*NamedAddrs

		// Cidr list
		RemoteCidrs []*NamedAddrs

		// LocalType of resource
		LocalType ResourceType
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

	BlockedResources struct {
		BlockedSubnets   map[ID]bool
		BlockedInstances map[ID]bool
		BlockedVPEs      map[ID]bool
	}

	VPCDetails struct {
		AddressPrefixes *netset.IPBlock
	}

	SubnetDetails struct {
		NamedEntity
		CIDR   *netset.IPBlock
		VPC    ID
		LRPair *LocalRemotePair // caching lookup result
	}

	NifDetails struct {
		NamedEntity
		IP       *netset.IPBlock
		VPC      ID
		Instance ID
		Subnet   ID
		LRPair   *LocalRemotePair // caching lookup result
	}

	InstanceDetails struct {
		NamedEntity
		VPC    ID
		Nifs   []ID
		LRPair *LocalRemotePair // caching lookup result
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
		LRPair         *LocalRemotePair // caching lookup result
	}

	SegmentDetails struct {
		Elements []ID
		LRPair   *LocalRemotePair // caching lookup result
	}

	CidrSegmentDetails struct {
		Cidrs            *netset.IPBlock
		ContainedSubnets []ID
		LRPair           *LocalRemotePair // caching lookup result
	}

	ExternalDetails struct {
		ExternalAddrs *netset.IPBlock
		LRPair        *LocalRemotePair // caching lookup result
	}

	Reader interface {
		ReadSpec(filename string, defs *ConfigDefs) (*Spec, error)
	}

	Named interface {
		Name() string
	}

	NWResource interface {
		Address() *netset.IPBlock
		getLocalRemotePair() *LocalRemotePair
		setLocalRemotePair(l *LocalRemotePair)
	}

	// resources that are in a subnet. used for lookupContainerForACLSynth generic function
	SubSubnetResource interface {
		Address() *netset.IPBlock
		SubnetName() ID
	}

	EndpointProvider interface {
		endpointNames() []ID
		endpointMap(s *Definitions) map[ID]SubSubnetResource
		getLocalRemotePair() *LocalRemotePair
		setLocalRemotePair(l *LocalRemotePair)
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

func (s *SubnetDetails) getLocalRemotePair() *LocalRemotePair {
	return s.LRPair
}

func (s *SubnetDetails) setLocalRemotePair(l *LocalRemotePair) {
	s.LRPair = l
}

func (n *NifDetails) Address() *netset.IPBlock {
	return n.IP
}

func (n *NifDetails) SubnetName() ID {
	return n.Subnet
}

func (n *NifDetails) getLocalRemotePair() *LocalRemotePair {
	return n.LRPair
}

func (n *NifDetails) setLocalRemotePair(l *LocalRemotePair) {
	n.LRPair = l
}

func (v *VPEReservedIPsDetails) Address() *netset.IPBlock {
	return v.IP
}

func (v *VPEReservedIPsDetails) SubnetName() ID {
	return v.Subnet
}

func (e *ExternalDetails) Address() *netset.IPBlock {
	return e.ExternalAddrs
}

func (e *ExternalDetails) getLocalRemotePair() *LocalRemotePair {
	return e.LRPair
}

func (e *ExternalDetails) setLocalRemotePair(l *LocalRemotePair) {
	e.LRPair = l
}

func (i *InstanceDetails) endpointNames() []ID {
	return i.Nifs
}

func (i *InstanceDetails) endpointMap(s *Definitions) map[ID]SubSubnetResource {
	res := make(map[ID]SubSubnetResource, len(i.Nifs))
	for _, nifName := range i.Nifs {
		res[nifName] = s.NIFs[nifName]
	}
	return res
}

func (i *InstanceDetails) getLocalRemotePair() *LocalRemotePair {
	return i.LRPair
}

func (i *InstanceDetails) setLocalRemotePair(l *LocalRemotePair) {
	i.LRPair = l
}

func (v *VPEDetails) endpointNames() []ID {
	return v.VPEReservedIPs
}

func (v *VPEDetails) endpointMap(s *Definitions) map[ID]SubSubnetResource {
	res := make(map[ID]SubSubnetResource, len(v.VPEReservedIPs))
	for _, ripName := range v.VPEReservedIPs {
		res[ripName] = s.VPEReservedIPs[ripName]
	}
	return res
}

func (v *VPEDetails) getLocalRemotePair() *LocalRemotePair {
	return v.LRPair
}

func (v *VPEDetails) setLocalRemotePair(l *LocalRemotePair) {
	v.LRPair = l
}

// lookupSingle is called only when the resource type is ResourceTypeSubnet or ResourceTypeExternal
func lookupSingle[T NWResource](m map[ID]T, name string, t ResourceType) (*LocalRemotePair, error) {
	details, ok := m[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, name, t)
	}
	if details.getLocalRemotePair() != nil {
		return details.getLocalRemotePair(), nil
	}
	res := &LocalRemotePair{
		Name:        &name,
		LocalCidrs:  []*NamedAddrs{{Name: &name, IPAddrs: details.Address()}},
		RemoteCidrs: []*NamedAddrs{{Name: &name, IPAddrs: details.Address()}},
		LocalType:   t,
	}
	details.setLocalRemotePair(res)
	return res, nil
}

func (s *Definitions) lookupSegment(segment map[ID]*SegmentDetails, name string, t, elementType ResourceType,
	lookup func(ResourceType, string) (*LocalRemotePair, error)) (*LocalRemotePair, error) {
	segmentDetails, ok := segment[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, t)
	}
	if segmentDetails.LRPair != nil {
		return segmentDetails.LRPair, nil
	}

	res := &LocalRemotePair{Name: &name, LocalType: elementType}
	for _, elementName := range segmentDetails.Elements {
		element, err := lookup(elementType, elementName)
		if err != nil {
			return nil, err
		}
		res.LocalCidrs = append(res.LocalCidrs, element.LocalCidrs...)
		res.RemoteCidrs = append(res.RemoteCidrs, element.RemoteCidrs...)
	}
	segmentDetails.LRPair = res
	return res, nil
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
