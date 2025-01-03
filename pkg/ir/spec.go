/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
	"slices"
	"strings"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

type (
	ID           = string
	ResourceType string

	Spec struct {
		// Required connections
		Connections []*Connection

		Defs *Definitions

		// resources that does not appear in the Spec file
		*BlockedResources
	}

	Connection struct {
		// Egress resource
		Src *ConnectedResource

		// Ingress resource
		Dst *ConnectedResource

		// Allowed protocols
		TrackedProtocols []*TrackedProtocol

		// Provenance information
		Origin fmt.Stringer
	}

	ConnectedResource struct {
		Name string

		// CidrsWhenLocal is a list of endpoints (IP ranges or names) representing
		// this resource, to which firewall rules should be applied when allowing
		// a required connection
		CidrsWhenLocal []*NamedAddrs

		// CidrsWhenRemote is a list of endpoints for this resource that should be
		// specified as the remotes of firewall rules
		CidrsWhenRemote []*NamedAddrs

		ResourceType ResourceType
	}

	NamedAddrs struct {
		IPAddrs *netset.IPBlock
		Name    string
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

	// ConnectedResource is for caching lookup results
	SubnetDetails struct {
		CIDR              *netset.IPBlock
		ConnectedResource *ConnectedResource
	}

	NifDetails struct {
		IP                *netset.IPBlock
		Instance          ID
		Subnet            ID
		ConnectedResource *ConnectedResource
	}

	InstanceDetails struct {
		Nifs              []ID
		ConnectedResource *ConnectedResource
	}

	VPEReservedIPsDetails struct {
		IP      *netset.IPBlock
		VPEName ID
		Subnet  ID
	}

	VPEDetails struct {
		VPEReservedIPs    []ID
		ConnectedResource *ConnectedResource
	}

	SegmentDetails struct {
		Elements          []ID
		ConnectedResource *ConnectedResource
	}

	CidrSegmentDetails struct {
		Cidrs             *netset.IPBlock
		ConnectedResource *ConnectedResource
	}

	ExternalDetails struct {
		ExternalAddrs     *netset.IPBlock
		ConnectedResource *ConnectedResource
	}

	Reader interface {
		ReadSpec(filename string, defs *ConfigDefs) (*Spec, error)
	}

	// generalizes subnet and external resource types
	NWResource interface {
		Address() *netset.IPBlock
		getConnectedResource() *ConnectedResource
		setConnectedResource(r *ConnectedResource)
	}

	// resources that are in a subnet. used for lookupContainerForACLSynth generic function
	SubSubnetResource interface {
		Address() *netset.IPBlock
		SubnetName() ID
	}

	EndpointProvider interface {
		endpointNames() []ID
		endpointMap(s *Definitions) map[ID]SubSubnetResource
		getConnectedResource() *ConnectedResource
		setConnectedResource(r *ConnectedResource)
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

func (s *SubnetDetails) Address() *netset.IPBlock {
	return s.CIDR
}

func (s *SubnetDetails) getConnectedResource() *ConnectedResource {
	return s.ConnectedResource
}

func (s *SubnetDetails) setConnectedResource(r *ConnectedResource) {
	s.ConnectedResource = r
}

func (n *NifDetails) Address() *netset.IPBlock {
	return n.IP
}

func (n *NifDetails) SubnetName() ID {
	return n.Subnet
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

func (e *ExternalDetails) getConnectedResource() *ConnectedResource {
	return e.ConnectedResource
}

func (e *ExternalDetails) setConnectedResource(r *ConnectedResource) {
	e.ConnectedResource = r
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

func (i *InstanceDetails) getConnectedResource() *ConnectedResource {
	return i.ConnectedResource
}

func (i *InstanceDetails) setConnectedResource(r *ConnectedResource) {
	i.ConnectedResource = r
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

func (v *VPEDetails) getConnectedResource() *ConnectedResource {
	return v.ConnectedResource
}

func (v *VPEDetails) setConnectedResource(r *ConnectedResource) {
	v.ConnectedResource = r
}

// lookupSingle is called only when the resource type is ResourceTypeSubnet or ResourceTypeExternal
func lookupSingle[T NWResource](m map[ID]T, name string, t ResourceType) (*ConnectedResource, error) {
	details, ok := m[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, name, t)
	}
	if details.getConnectedResource() != nil {
		return details.getConnectedResource(), nil
	}
	res := &ConnectedResource{
		Name:            name,
		CidrsWhenLocal:  []*NamedAddrs{{Name: name, IPAddrs: details.Address()}},
		CidrsWhenRemote: []*NamedAddrs{{Name: name, IPAddrs: details.Address()}},
		ResourceType:    t,
	}
	details.setConnectedResource(res)
	return res, nil
}

func (s *Definitions) lookupSegment(segment map[ID]*SegmentDetails, name string, t, elementType ResourceType,
	lookup func(ResourceType, string) (*ConnectedResource, error)) (*ConnectedResource, error) {
	segmentDetails, ok := segment[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, t)
	}
	if segmentDetails.ConnectedResource != nil {
		return segmentDetails.ConnectedResource, nil
	}

	res := &ConnectedResource{Name: name, ResourceType: elementType}
	for _, elementName := range segmentDetails.Elements {
		element, err := lookup(elementType, elementName)
		if err != nil {
			return nil, err
		}
		res.CidrsWhenLocal = slices.Concat(res.CidrsWhenLocal, element.CidrsWhenLocal)
		res.CidrsWhenRemote = slices.Concat(res.CidrsWhenRemote, element.CidrsWhenRemote)
	}
	segmentDetails.ConnectedResource = res
	return res, nil
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
