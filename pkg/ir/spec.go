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
		Src *FirewallResource

		// Ingress resource
		Dst *FirewallResource

		// Allowed protocols
		TrackedProtocols []*TrackedProtocol

		// Provenance information
		Origin fmt.Stringer
	}

	FirewallResource struct {
		// Symbolic name of resource, if available
		Name *string

		AppliedTo []*NamedAddrs

		// Cidr list
		RemoteCidrs []*NamedAddrs

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

		BlockedResources

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

	INWResource interface {
		NWResource
		SubnetName() ID
	}

	EndpointProvider interface {
		endpointNames() []ID
		endpointMap(s *Definitions) map[ID]INWResource
		endpointType() ResourceType
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

	WarningUnspecifiedACL = "The following subnets do not have required connections; the generated ACL will block all traffic: "
	WarningUnspecifiedSG  = "The following endpoints do not have required connections; the generated SGs will block all traffic: "
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

func (i *InstanceDetails) endpointNames() []ID {
	return i.Nifs
}

func (i *InstanceDetails) endpointMap(s *Definitions) map[ID]INWResource {
	res := make(map[ID]INWResource, len(s.NIFs))
	for k, v := range s.NIFs {
		res[k] = v
	}
	return res
}

func (i *InstanceDetails) endpointType() ResourceType {
	return ResourceTypeNIF
}

func (v *VPEDetails) endpointNames() []ID {
	return v.VPEReservedIPs
}

func (v *VPEDetails) endpointMap(s *Definitions) map[ID]INWResource {
	res := make(map[ID]INWResource, len(s.VPEReservedIPs))
	for k, v := range s.VPEReservedIPs {
		res[k] = v
	}
	return res
}

func (v *VPEDetails) endpointType() ResourceType {
	return ResourceTypeVPE
}

func lookupSingle[T NWResource](m map[ID]T, name string, t ResourceType) (*FirewallResource, error) {
	if details, ok := m[name]; ok {
		return &FirewallResource{
			Name:        &name,
			AppliedTo:   []*NamedAddrs{{Name: &name, IPAddrs: details.Address()}},
			RemoteCidrs: []*NamedAddrs{{Name: &name, IPAddrs: details.Address()}},
			Type:        utils.Ptr(t),
		}, nil
	}
	return nil, fmt.Errorf(resourceNotFound, name, t)
}

func (s *Definitions) lookupSegment(segment map[ID]*SegmentDetails, name string, t, elementType ResourceType,
	lookup func(ResourceType, string) (*FirewallResource, error)) (*FirewallResource, error) {
	segmentDetails, ok := segment[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, t)
	}

	res := &FirewallResource{Name: &name, AppliedTo: []*NamedAddrs{}, RemoteCidrs: []*NamedAddrs{}, Type: utils.Ptr(elementType)}
	for _, elementName := range segmentDetails.Elements {
		subnet, err := lookup(elementType, elementName)
		if err != nil {
			return nil, err
		}
		res.AppliedTo = append(res.AppliedTo, subnet.AppliedTo...)
		res.RemoteCidrs = append(res.RemoteCidrs, subnet.RemoteCidrs...)
	}
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

func SetUnspecifiedWarning(warningPrefix string, blockedResources []ID) string {
	warning := ""
	if len(blockedResources) > 0 {
		warning = warningPrefix + strings.Join(blockedResources, ", ")
		log.Println(warning)
	}
	return warning
}
