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
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

func (s *Definitions) LookupForSGSynth(t ResourceType, name string) (*Resource, error) {
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
		return lookupContainerForSGSynth(s.Instances, name, ResourceTypeVPE)
	case ResourceTypeSubnetSegment:
		return s.lookupSubnetSegmentForSGSynth(name)
	case ResourceTypeCidrSegment:
		return s.lookupCidrSegmentForSGSynth(name)
	case ResourceTypeNifSegment:
		return s.lookupSegmentForSGSynth(s.NifSegments, name, ResourceTypeNIF)
	case ResourceTypeInstanceSegment:
		return s.lookupSegmentForSGSynth(s.InstanceSegments, name, ResourceTypeInstance)
	case ResourceTypeVpeSegment:
		return s.lookupSegmentForSGSynth(s.VpeSegments, name, ResourceTypeVPE)
	}
	return nil, nil // should not get here
}

func (s *Definitions) lookupNIFForSGSynth(name string) (*Resource, error) {
	if _, ok := s.NIFs[name]; ok {
		return &Resource{
			Name:       &name,
			NamedAddrs: []*NamedAddrs{{Name: &s.NIFs[name].Instance}},
			Cidrs:      []*NamedAddrs{},
			Type:       utils.Ptr(ResourceTypeNIF),
		}, nil
	}
	return nil, fmt.Errorf(resourceNotFound, ResourceTypeNIF, name)
}

func lookupContainerForSGSynth[T EndpointProvider](m map[string]T, name string, t ResourceType) (*Resource, error) {
	if _, ok := m[name]; ok {
		return &Resource{
			Name:       &name,
			NamedAddrs: []*NamedAddrs{{Name: &name}},
			Cidrs:      []*NamedAddrs{},
			Type:       utils.Ptr(t),
		}, nil
	}
	return nil, fmt.Errorf(containerNotFound, t, name)
}

func (s *Definitions) lookupSubnetForSGSynth(name string) (*Resource, error) {
	if subnetDetails, ok := s.Subnets[name]; ok {
		return &Resource{Name: &name,
			NamedAddrs: s.containedResourcesInCidr(subnetDetails.CIDR),
			Cidrs:      []*NamedAddrs{{IPAddrs: subnetDetails.CIDR}},
			Type:       utils.Ptr(ResourceTypeSubnet),
		}, nil
	}
	return nil, fmt.Errorf(resourceNotFound, ResourceTypeSubnet, name)
}

func (s *Definitions) lookupSubnetSegmentForSGSynth(name string) (*Resource, error) {
	segmentDetails, ok := s.SubnetSegments[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, ResourceTypeSubnetSegment, name)
	}

	res := &Resource{Name: &name, NamedAddrs: []*NamedAddrs{}, Cidrs: []*NamedAddrs{}, Type: utils.Ptr(ResourceTypeSubnet)}
	for _, subnetName := range segmentDetails.Elements {
		subnetRes, err := s.lookupSubnetForSGSynth(subnetName)
		if err != nil {
			return nil, err
		}
		res.NamedAddrs = append(res.NamedAddrs, subnetRes.NamedAddrs...)
		res.Cidrs = append(res.Cidrs, subnetRes.Cidrs...)
	}
	return res, nil
}

func (s *Definitions) lookupCidrSegmentForSGSynth(name string) (*Resource, error) {
	if segmentDetails, ok := s.CidrSegments[name]; ok {
		return &Resource{Name: &name,
			NamedAddrs: s.containedResourcesInCidr(segmentDetails.Cidrs),
			Cidrs:      cidrToNamedAddrs(segmentDetails.Cidrs),
			Type:       utils.Ptr(ResourceTypeCidr),
		}, nil
	}
	return nil, fmt.Errorf(containerNotFound, ResourceTypeCidrSegment, name)
}

func (s *Definitions) lookupSegmentForSGSynth(segment map[string]*SegmentDetails, name string, t ResourceType) (*Resource, error) {
	if segmentDetails, ok := segment[name]; ok {
		return &Resource{Name: &name, NamedAddrs: namesToNamedAddrs(segmentDetails.Elements), Cidrs: []*NamedAddrs{}, Type: utils.Ptr(t)}, nil
	}
	return nil, fmt.Errorf(containerNotFound, name, t)
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
	slices.Compact(slices.Sorted(slices.Values(names)))
	return namesToNamedAddrs(names)
}

func cidrToNamedAddrs(cidr *netset.IPBlock) []*NamedAddrs {
	cidrs := cidr.SplitToCidrs()
	res := make([]*NamedAddrs, len(cidrs))
	for _, c := range cidrs {
		res = append(res, &NamedAddrs{IPAddrs: c})
	}
	return res
}

func namesToNamedAddrs(names []string) []*NamedAddrs {
	res := make([]*NamedAddrs, len(names))
	for _, name := range names {
		res = append(res, &NamedAddrs{Name: &name})
	}
	return res
}
