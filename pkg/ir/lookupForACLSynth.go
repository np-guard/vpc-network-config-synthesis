/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

func (s *Definitions) LookupForACLSynth(t ResourceType, name string) (*Resource, error) {
	switch t {
	case ResourceTypeExternal:
		return lookupSingle(s.Externals, name, t)
	case ResourceTypeSubnet:
		return lookupSingle(s.Subnets, name, t)
	case ResourceTypeNIF:
		return lookupSingleForACLSynth(s.NIFs, s.Subnets, name, t)
	case ResourceTypeInstance:
		return lookupContainerForACLSynth(s.Instances, s, name, ResourceTypeInstance)
	case ResourceTypeVPE:
		return lookupContainerForACLSynth(s.VPEs, s, name, ResourceTypeVPE)
	case ResourceTypeSubnetSegment:
		return s.lookupSegment(s.SubnetSegments, name, t, ResourceTypeSubnet, s.LookupForACLSynth)
	case ResourceTypeCidrSegment:
		return s.lookupCidrSegmentACL(name)
	case ResourceTypeNifSegment:
		return s.lookupSegment(s.NifSegments, name, t, ResourceTypeNIF, s.LookupForACLSynth)
	case ResourceTypeInstanceSegment:
		return s.lookupSegment(s.InstanceSegments, name, t, ResourceTypeInstance, s.LookupForACLSynth)
	case ResourceTypeVpeSegment:
		return s.lookupSegment(s.VpeSegments, name, t, ResourceTypeVPE, s.LookupForACLSynth)
	}
	return nil, nil // should not get here
}

func lookupSingleForACLSynth[T INWResource](m map[ID]T, subnets map[ID]*SubnetDetails, name string, t ResourceType) (*Resource, error) {
	details, ok := m[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, name, t)
	}

	res, err := lookupSingle(subnets, details.SubnetName(), ResourceTypeSubnet)
	if err != nil {
		return nil, err
	}
	res.Name = &name // save NIF's name as resource name
	return res, nil
}

func lookupContainerForACLSynth[T EndpointProvider](m map[ID]T, defs *Definitions, name string, t ResourceType) (*Resource, error) {
	containerDetails, ok := m[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, t)
	}

	res := &Resource{Name: &name, NamedAddrs: []*NamedAddrs{}, Cidrs: []*NamedAddrs{}, Type: utils.Ptr(ResourceTypeSubnet)}
	endpointMap := containerDetails.endpointMap(defs)
	for _, endpointName := range containerDetails.endpointNames() {
		subnet, err := lookupSingleForACLSynth(endpointMap, defs.Subnets, endpointName, containerDetails.endpointType())
		if err != nil {
			return nil, err
		}
		res.Cidrs = append(res.Cidrs, subnet.Cidrs...)
		res.NamedAddrs = append(res.NamedAddrs, subnet.NamedAddrs...)
	}
	return res, nil
}

func (s *Definitions) lookupCidrSegmentACL(name string) (*Resource, error) {
	segmentDetails, ok := s.CidrSegments[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, ResourceTypeCidrSegment)
	}

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
