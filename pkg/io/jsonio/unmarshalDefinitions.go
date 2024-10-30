/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package jsonio handles global specification written in a JSON file, as described by spec_schema.input
package jsonio

import (
	"fmt"
	"slices"

	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/models/pkg/spec"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type segmentsTypes struct {
	subnetSegment   map[string][]string
	cidrSegment     map[string][]string
	nifSegment      map[string][]string
	instanceSegment map[string][]string
	vpeSegment      map[string][]string
}

// ReadDefinitions translates segments and externals
func (r *Reader) readDefinitions(jsonSpec *spec.Spec, configDefs *ir.ConfigDefs) (*ir.Definitions, *ir.BlockedResources, error) {
	if err := validateSegments(&jsonSpec.Segments); err != nil {
		return nil, nil, err
	}
	segments := divideSegmentsByType(&jsonSpec.Segments)
	subnetSegments := parseSegments(segments.subnetSegment)
	nifSegments := parseSegments(segments.nifSegment)
	instanceSegments := parseSegments(segments.instanceSegment)
	vpeSegments := parseSegments(segments.vpeSegment)
	cidrSegments, err := parseCidrSegments(segments.cidrSegment, configDefs)
	if err != nil {
		return nil, nil, err
	}
	externals, err := translateExternals(jsonSpec.Externals)
	if err != nil {
		return nil, nil, err
	}
	return &ir.Definitions{
		ConfigDefs:       *configDefs,
		SubnetSegments:   subnetSegments,
		CidrSegments:     cidrSegments,
		NifSegments:      nifSegments,
		InstanceSegments: instanceSegments,
		VpeSegments:      vpeSegments,
		Externals:        externals,
	}, prepareBlockedResources(configDefs), nil
}

// validateSegments validates that all segments are supported
func validateSegments(jsonSegments *spec.SpecSegments) error {
	for _, v := range *jsonSegments {
		if v.Type != spec.SegmentTypeSubnet && v.Type != spec.SegmentTypeCidr &&
			v.Type != spec.SegmentTypeInstance && v.Type != spec.SegmentTypeNif &&
			v.Type != spec.SegmentTypeVpe {
			return fmt.Errorf("only subnet, cidr, instance, nif and vpe segments are supported, not %q", v.Type)
		}
	}
	return nil
}

// translates segment to ir ds
func parseSegments(segments map[string][]string) map[ir.ID]*ir.SegmentDetails {
	result := make(map[string]*ir.SegmentDetails)
	for segmentName, elements := range segments {
		result[segmentName] = &ir.SegmentDetails{Elements: elements}
	}
	return result
}

// parseCidrSegments translates cidr segments
func parseCidrSegments(cidrSegments map[string][]string, configDefs *ir.ConfigDefs) (map[ir.ID]*ir.CidrSegmentDetails, error) {
	result := make(map[ir.ID]*ir.CidrSegmentDetails)
	for segmentName, segment := range cidrSegments {
		cidrs := netset.NewIPBlock()
		containedSubnets := make([]ir.ID, 0)

		for _, cidr := range segment {
			c, err := netset.IPBlockFromCidr(cidr)
			if err != nil {
				return nil, err
			}
			subnets, err := configDefs.SubnetsContainedInCidr(c)
			if err != nil {
				return nil, err
			}

			cidrs = cidrs.Union(c)
			containedSubnets = append(containedSubnets, subnets...)
		}
		if !internalCidr(configDefs, cidrs) {
			return nil, fmt.Errorf("only internal cidrs are supported in cidr segment resource type (segment name: %v)", segmentName)
		}
		cidrSegmentDetails := ir.CidrSegmentDetails{
			Cidrs:            cidrs,
			ContainedSubnets: slices.Compact(slices.Sorted(slices.Values(containedSubnets))),
		}
		result[segmentName] = &cidrSegmentDetails
	}
	return result, nil
}

// translateExternals reads externals from spec file
func translateExternals(m map[string]string) (map[ir.ID]*ir.ExternalDetails, error) {
	result := make(map[ir.ID]*ir.ExternalDetails)
	for k, v := range m {
		address, err := netset.IPBlockFromCidrOrAddress(v)
		if err != nil {
			return nil, err
		}
		result[k] = &ir.ExternalDetails{ExternalAddrs: address}
	}
	return result, nil
}

func divideSegmentsByType(jsonSegments *spec.SpecSegments) segmentsTypes {
	res := segmentsTypes{subnetSegment: make(map[string][]string), nifSegment: make(map[string][]string),
		instanceSegment: make(map[string][]string), vpeSegment: make(map[string][]string),
		cidrSegment: make(map[string][]string)}
	for k, v := range *jsonSegments {
		switch v.Type {
		case spec.SegmentTypeSubnet:
			res.subnetSegment[k] = v.Items
		case spec.SegmentTypeCidr:
			res.cidrSegment[k] = v.Items
		case spec.SegmentTypeInstance:
			res.instanceSegment[k] = v.Items
		case spec.SegmentTypeNif:
			res.nifSegment[k] = v.Items
		case spec.SegmentTypeVpe:
			res.vpeSegment[k] = v.Items
		}
	}
	return res
}

func internalCidr(configDefs *ir.ConfigDefs, cidr *netset.IPBlock) bool {
	res := cidr
	for _, vpcDetails := range configDefs.VPCs {
		res = res.Subtract(vpcDetails.AddressPrefixes)
	}
	return res.IsEmpty()
}

func prepareBlockedResources(configDefs *ir.ConfigDefs) *ir.BlockedResources {
	return &ir.BlockedResources{BlockedSubnets: sliceToMap(utils.SortedMapKeys(configDefs.Subnets)),
		BlockedInstances: sliceToMap(utils.SortedMapKeys(configDefs.Instances)),
		BlockedVPEs:      sliceToMap(utils.SortedMapKeys(configDefs.VPEs))}
}

func sliceToMap(slice []string) map[string]bool {
	res := make(map[string]bool, len(slice))
	for _, elem := range slice {
		res[elem] = true
	}
	return res
}
