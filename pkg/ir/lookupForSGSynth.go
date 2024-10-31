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

func (s *Definitions) LookupForSGSynth(t ResourceType, name string) (*LocalRemotePair, error) {
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

func (s *Definitions) lookupNIFForSGSynth(name string) (*LocalRemotePair, error) {
	details, ok := s.NIFs[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, ResourceTypeNIF, name)
	}
	if details.LRPair != nil {
		return details.LRPair, nil
	}
	details.LRPair = &LocalRemotePair{
		Name:        &name,
		LocalCidrs:  []*NamedAddrs{{Name: &details.Instance}},
		RemoteCidrs: []*NamedAddrs{{Name: &details.Instance}},
		LocalType:   ResourceTypeNIF,
	}
	return details.LRPair, nil
}

func lookupContainerForSGSynth[T EndpointProvider](m map[string]T, name string, t ResourceType) (*LocalRemotePair, error) {
	details, ok := m[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, t, name)
	}
	if details.getLocalRemotePair() != nil {
		return details.getLocalRemotePair(), nil
	}
	lrPair := &LocalRemotePair{
		Name:        &name,
		LocalCidrs:  []*NamedAddrs{{Name: &name}},
		RemoteCidrs: []*NamedAddrs{{Name: &name}},
		LocalType:   t,
	}
	details.setLocalRemotePair(lrPair)
	return lrPair, nil
}

func (s *Definitions) lookupSubnetForSGSynth(name string) (*LocalRemotePair, error) {
	subnetDetails, ok := s.Subnets[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, ResourceTypeSubnet, name)
	}
	subnetDetails.LRPair = &LocalRemotePair{Name: &name,
		LocalCidrs:  s.containedResourcesInCidr(subnetDetails.CIDR),
		RemoteCidrs: []*NamedAddrs{{IPAddrs: subnetDetails.CIDR}},
		LocalType:   ResourceTypeSubnet,
	}
	return subnetDetails.LRPair, nil
}

func (s *Definitions) lookupCidrSegmentForSGSynth(name string) (*LocalRemotePair, error) {
	segmentDetails, ok := s.CidrSegments[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, ResourceTypeCidrSegment, name)
	}
	if segmentDetails.LRPair != nil {
		return segmentDetails.LRPair, nil
	}
	segmentDetails.LRPair = &LocalRemotePair{Name: &name,
		LocalCidrs:  s.containedResourcesInCidr(segmentDetails.Cidrs),
		RemoteCidrs: cidrToNamedAddrs(segmentDetails.Cidrs),
		LocalType:   ResourceTypeCidr,
	}
	return segmentDetails.LRPair, nil
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
		res[i] = &NamedAddrs{Name: &name}
	}
	return res
}
