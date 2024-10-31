/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
)

func (s *Definitions) LookupForACLSynth(t ResourceType, name string) (*LocalRemotePair, error) {
	switch t {
	case ResourceTypeExternal:
		return lookupSingle(s.Externals, name, t)
	case ResourceTypeSubnet:
		return lookupSingle(s.Subnets, name, t)
	case ResourceTypeNIF:
		return s.lookupNIFForACLSynth(name)
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

func lookupContainerForACLSynth[T EndpointProvider](m map[ID]T, defs *Definitions, name string, t ResourceType) (*LocalRemotePair, error) {
	containerDetails, ok := m[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, t)
	}
	if containerDetails.getLocalRemotePair() != nil {
		return containerDetails.getLocalRemotePair(), nil
	}

	seenSubnets := make(map[string]struct{})
	res := &LocalRemotePair{Name: &name, LocalCidrs: []*NamedAddrs{}, RemoteCidrs: []*NamedAddrs{}, LocalType: ResourceTypeSubnet}
	endpointMap := containerDetails.endpointMap(defs)
	for _, endpointName := range containerDetails.endpointNames() {
		subnetName := endpointMap[endpointName].SubnetName()
		if _, ok := seenSubnets[subnetName]; ok {
			continue
		}
		seenSubnets[subnetName] = struct{}{}

		namedAddrs := &NamedAddrs{Name: &subnetName, IPAddrs: defs.Subnets[subnetName].CIDR}
		res.RemoteCidrs = append(res.RemoteCidrs, namedAddrs)
		res.LocalCidrs = append(res.LocalCidrs, namedAddrs)
	}
	containerDetails.setLocalRemotePair(res)
	return res, nil
}

func (s *Definitions) lookupNIFForACLSynth(name string) (*LocalRemotePair, error) {
	details, ok := s.NIFs[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, name, ResourceTypeNIF)
	}

	NifSubnetName := details.Subnet
	NifSubnetCidr := s.Subnets[NifSubnetName].CIDR
	return &LocalRemotePair{
		Name:        &name,
		LocalCidrs:  []*NamedAddrs{{Name: &NifSubnetName, IPAddrs: NifSubnetCidr}},
		RemoteCidrs: []*NamedAddrs{{Name: &NifSubnetName, IPAddrs: NifSubnetCidr}},
		LocalType:   ResourceTypeSubnet,
	}, nil
}

func (s *Definitions) lookupCidrSegmentACL(name string) (*LocalRemotePair, error) {
	segmentDetails, ok := s.CidrSegments[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, ResourceTypeCidrSegment)
	}
	if segmentDetails.LRPair != nil {
		return segmentDetails.LRPair, nil
	}

	res := &LocalRemotePair{Name: &name, LocalCidrs: []*NamedAddrs{}, RemoteCidrs: []*NamedAddrs{}, LocalType: ResourceTypeSubnet}
	for _, subnetName := range segmentDetails.ContainedSubnets {
		subnet, err := lookupSingle(s.Subnets, subnetName, ResourceTypeSubnet)
		if err != nil {
			return nil, fmt.Errorf("%w while looking up %v %v for cidr segment %v", err, ResourceTypeSubnet, subnetName, name)
		}
		res.LocalCidrs = append(res.LocalCidrs, subnet.LocalCidrs...)
	}
	for _, cidr := range segmentDetails.Cidrs.SplitToCidrs() {
		res.RemoteCidrs = append(res.RemoteCidrs, &NamedAddrs{Name: &name, IPAddrs: cidr})
	}
	segmentDetails.LRPair = res
	return res, nil
}
