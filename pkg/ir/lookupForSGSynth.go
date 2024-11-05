/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
	"slices"

	"github.com/np-guard/models/pkg/netset"
)

func (s *Definitions) LookupForSGSynth(t ResourceType, name string) (*ConnectedResource, error) {
	switch t {
	case ResourceTypeExternal:
		return lookupSingle(s.Externals, name, t)
	case ResourceTypeSubnet:
		return s.lookupSubnetForSGSynth(name)
	case ResourceTypeNIF:
		return s.lookupNIFForSGSynth(name)
	case ResourceTypeInstance:
		return lookupContainerForSGSynth(s.Instances, name, ResourceTypeInstance)
	case ResourceTypeVPE:
		return lookupContainerForSGSynth(s.VPEs, name, ResourceTypeVPE)
	case ResourceTypeSubnetSegment:
		return s.lookupSegment(s.SubnetSegments, name, t, ResourceTypeSubnet, s.LookupForSGSynth)
	case ResourceTypeCidrSegment:
		return s.lookupCidrSegmentForSGSynth(name)
	case ResourceTypeNifSegment:
		return s.lookupSegment(s.NifSegments, name, t, ResourceTypeNIF, s.LookupForSGSynth)
	case ResourceTypeInstanceSegment:
		return s.lookupSegment(s.InstanceSegments, name, t, ResourceTypeInstance, s.LookupForSGSynth)
	case ResourceTypeVpeSegment:
		return s.lookupSegment(s.VpeSegments, name, t, ResourceTypeVPE, s.LookupForSGSynth)
	}
	return nil, nil // should not get here
}

func (s *Definitions) lookupNIFForSGSynth(name string) (*ConnectedResource, error) {
	details, ok := s.NIFs[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, ResourceTypeNIF, name)
	}
	if details.ConnectedResource != nil {
		return details.ConnectedResource, nil
	}
	details.ConnectedResource = &ConnectedResource{
		Name:            name,
		CidrsWhenLocal:  []*NamedAddrs{{Name: details.Instance}},
		CidrsWhenRemote: []*NamedAddrs{{Name: details.Instance}},
		ResourceType:    ResourceTypeInstance,
	}
	return details.ConnectedResource, nil
}

func lookupContainerForSGSynth[T EndpointProvider](m map[string]T, name string, t ResourceType) (*ConnectedResource, error) {
	details, ok := m[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, t, name)
	}
	if details.getConnectedResource() != nil {
		return details.getConnectedResource(), nil
	}
	res := &ConnectedResource{
		Name:            name,
		CidrsWhenLocal:  []*NamedAddrs{{Name: name}},
		CidrsWhenRemote: []*NamedAddrs{{Name: name}},
		ResourceType:    t,
	}
	details.setConnectedResource(res)
	return res, nil
}

func (s *Definitions) lookupSubnetForSGSynth(name string) (*ConnectedResource, error) {
	subnetDetails, ok := s.Subnets[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, ResourceTypeSubnet, name)
	}
	if subnetDetails.ConnectedResource != nil {
		return subnetDetails.ConnectedResource, nil
	}
	subnetDetails.ConnectedResource = &ConnectedResource{Name: name,
		CidrsWhenLocal:  s.containedResourcesInCidr(subnetDetails.CIDR),
		CidrsWhenRemote: []*NamedAddrs{{IPAddrs: subnetDetails.CIDR}},
		ResourceType:    ResourceTypeSubnet,
	}
	return subnetDetails.ConnectedResource, nil
}

func (s *Definitions) lookupCidrSegmentForSGSynth(name string) (*ConnectedResource, error) {
	segmentDetails, ok := s.CidrSegments[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, ResourceTypeCidrSegment, name)
	}
	if segmentDetails.ConnectedResource != nil {
		return segmentDetails.ConnectedResource, nil
	}
	segmentDetails.ConnectedResource = &ConnectedResource{Name: name,
		CidrsWhenLocal:  s.containedResourcesInCidr(segmentDetails.Cidrs),
		CidrsWhenRemote: cidrToNamedAddrs(segmentDetails.Cidrs),
		ResourceType:    ResourceTypeCidr,
	}
	return segmentDetails.ConnectedResource, nil
}

func (s *Definitions) containedResourcesInCidr(cidr *netset.IPBlock) []*NamedAddrs {
	names := make([]string, 0)
	for _, nifDetails := range s.NIFs {
		if nifDetails.IP.IsSubset(cidr) {
			names = append(names, nifDetails.Instance)
		}
	}
	for _, reservedIPName := range s.VPEReservedIPs {
		if reservedIPName.IP.IsSubset(cidr) {
			names = append(names, reservedIPName.VPEName)
		}
	}
	return namesToNamedAddrs(slices.Compact(slices.Sorted(slices.Values(names))))
}

func cidrToNamedAddrs(cidr *netset.IPBlock) []*NamedAddrs {
	cidrs := cidr.SplitToCidrs()
	res := make([]*NamedAddrs, len(cidrs))
	for i, c := range cidrs {
		res[i] = &NamedAddrs{IPAddrs: c}
	}
	return res
}

func namesToNamedAddrs(names []string) []*NamedAddrs {
	res := make([]*NamedAddrs, len(names))
	for i, name := range names {
		res[i] = &NamedAddrs{Name: name}
	}
	return res
}
