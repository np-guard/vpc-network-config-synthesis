/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package ir describes the input-format-agnostic specification of the required connectivity
package ir

import (
	"fmt"
	"sort"

	"github.com/np-guard/models/pkg/netset"
)

func (s *Definitions) LookupForACLSynth(t ResourceType, name string) (*ConnectedResource, error) {
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

func lookupContainerForACLSynth[T EndpointProvider](m map[ID]T, defs *Definitions, name string,
	t ResourceType) (*ConnectedResource, error) {
	containerDetails, ok := m[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, t)
	}
	if containerDetails.getConnectedResource() != nil {
		return containerDetails.getConnectedResource(), nil
	}

	seenSubnets := make(map[string]struct{})
	res := &ConnectedResource{Name: name, ResourceType: ResourceTypeSubnet}
	endpointMap := containerDetails.endpointMap(defs)
	for _, endpointName := range containerDetails.endpointNames() {
		subnetName := endpointMap[endpointName].SubnetName()
		if _, ok := seenSubnets[subnetName]; ok {
			continue
		}
		seenSubnets[subnetName] = struct{}{}

		namedAddrs := &NamedAddrs{Name: subnetName, IPAddrs: defs.Subnets[subnetName].CIDR}
		res.CidrsWhenRemote = append(res.CidrsWhenRemote, namedAddrs)
		res.CidrsWhenLocal = append(res.CidrsWhenLocal, namedAddrs)
	}
	containerDetails.setConnectedResource(res)
	return res, nil
}

func (s *Definitions) lookupNIFForACLSynth(name string) (*ConnectedResource, error) {
	details, ok := s.NIFs[name]
	if !ok {
		return nil, fmt.Errorf(resourceNotFound, name, ResourceTypeNIF)
	}
	if details.ConnectedResource != nil {
		return details.ConnectedResource, nil
	}

	NifSubnetName := details.Subnet
	NifSubnetCidr := s.Subnets[NifSubnetName].CIDR
	details.ConnectedResource = &ConnectedResource{
		Name:            name,
		CidrsWhenLocal:  []*NamedAddrs{{Name: NifSubnetName, IPAddrs: NifSubnetCidr}},
		CidrsWhenRemote: []*NamedAddrs{{Name: NifSubnetName, IPAddrs: NifSubnetCidr}},
		ResourceType:    ResourceTypeSubnet,
	}
	return details.ConnectedResource, nil
}

func (s *Definitions) lookupCidrSegmentACL(name string) (*ConnectedResource, error) {
	segmentDetails, ok := s.CidrSegments[name]
	if !ok {
		return nil, fmt.Errorf(containerNotFound, name, ResourceTypeCidrSegment)
	}
	if segmentDetails.ConnectedResource != nil {
		return segmentDetails.ConnectedResource, nil
	}
	res := &ConnectedResource{Name: name,
		CidrsWhenLocal: s.containedSubnetsInCidr(segmentDetails.Cidrs),
		ResourceType:   ResourceTypeSubnet,
	}
	for _, cidr := range segmentDetails.Cidrs.SplitToCidrs() {
		res.CidrsWhenRemote = append(res.CidrsWhenRemote, &NamedAddrs{Name: name, IPAddrs: cidr})
	}
	segmentDetails.ConnectedResource = res
	return res, nil
}

func (s *Definitions) containedSubnetsInCidr(cidr *netset.IPBlock) []*NamedAddrs {
	res := make([]*NamedAddrs, 0)
	for subnet, subnetDetails := range s.Subnets {
		if subnetDetails.CIDR.IsSubset(cidr) {
			res = append(res, &NamedAddrs{Name: subnet, IPAddrs: subnetDetails.CIDR})
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Name < res[j].Name })
	return res
}
