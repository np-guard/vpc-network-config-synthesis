/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package jsonio handles global specification written in a JSON file, as described by spec_schema.input
package jsonio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/np-guard/models/pkg/ipblock"
	"github.com/np-guard/models/pkg/spec"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Reader implements ir.Reader
type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (*Reader) ReadSpec(filename string, configDefs *ir.ConfigDefs) (*ir.Spec, error) {
	jsonSpec, err := unmarshal(filename)
	if err != nil {
		return nil, err
	}

	err = validateSegments(jsonSpec.Segments)
	if err != nil {
		return nil, err
	}

	subnetSegments := parseSubnetSegments(jsonSpec.Segments)

	cidrSegments, err := parseCidrSegments(jsonSpec.Segments, configDefs)
	if err != nil {
		return nil, err
	}

	// replace to fully qualified name
	jsonSpec, finalSubnetSegments, err := replaceResourcesName(jsonSpec, subnetSegments, configDefs)
	if err != nil {
		return nil, err
	}

	externals, err := translateExternals(jsonSpec.Externals)
	if err != nil {
		return nil, err
	}

	defs := &ir.Definitions{
		ConfigDefs:     *configDefs,
		SubnetSegments: finalSubnetSegments,
		CidrSegments:   cidrSegments,
		Externals:      externals,
	}

	var connections []ir.Connection
	for i := range jsonSpec.RequiredConnections {
		bidiConns, err := translateConnection(defs, &jsonSpec.RequiredConnections[i], i)
		if err != nil {
			return nil, err
		}
		connections = append(connections, bidiConns...)
	}

	return &ir.Spec{
		Connections: connections,
		Defs:        *defs,
	}, nil
}

func validateSegments(jsonSegments spec.SpecSegments) error {
	for _, v := range jsonSegments {
		if v.Type != spec.SegmentTypeSubnet && v.Type != spec.SegmentTypeCidr {
			return fmt.Errorf("only subnet and cidr segments are supported, not %q", v.Type)
		}
	}
	return nil
}

func filterSegmentsByType(jsonSegments spec.SpecSegments, segmentType spec.SegmentType) map[string][]string {
	result := make(map[string][]string)
	for k, v := range jsonSegments {
		if v.Type == segmentType {
			result[k] = v.Items
		}
	}
	return result
}

func parseSubnetSegments(jsonSegments spec.SpecSegments) map[string][]ir.ID {
	subnetSegments := filterSegmentsByType(jsonSegments, spec.SegmentTypeSubnet)
	result := make(map[string][]ir.ID)
	for segmentName, subnets := range subnetSegments {
		result[segmentName] = subnets
	}
	return result
}

func parseCidrSegments(jsonSegments spec.SpecSegments, configDefs *ir.ConfigDefs) (map[ir.ID]*ir.CidrSegmentDetails, error) {
	cidrSegments := filterSegmentsByType(jsonSegments, spec.SegmentTypeCidr)
	result := make(map[ir.ID]*ir.CidrSegmentDetails)
	for segmentName, segment := range cidrSegments {
		// each cidr saves the contained subnets
		segmentMap := make(map[*ipblock.IPBlock]ir.CIDRDetails)
		for _, cidr := range segment {
			c, err := ipblock.FromCidr(cidr)
			if err != nil {
				return nil, err
			}
			vpcs := parseOverlappingVpcs(c, configDefs.VPCs)
			subnets, err := configDefs.SubnetsContainedInCidr(*c)
			if err != nil {
				return nil, err
			}
			segmentMap[c] = ir.CIDRDetails{
				OverlappingVPCs:  vpcs,
				ContainedSubnets: subnets,
			}
		}
		cidrSegmentDetails := ir.CidrSegmentDetails{
			Cidrs: segmentMap,
		}
		result[segmentName] = &cidrSegmentDetails
	}
	return result, nil
}

// read externals from conn-spec
func translateExternals(m map[string]string) (map[ir.ID]*ir.ExternalDetails, error) {
	result := make(map[ir.ID]*ir.ExternalDetails)
	for k, v := range m {
		address, err := ipblock.FromCidrOrAddress(v)
		if err != nil {
			return nil, err
		}
		result[k] = &ir.ExternalDetails{ExternalAddrs: address}
	}
	return result, nil
}

// replace all resources names in conn-spec to fully qualified name
func replaceResourcesName(jsonSpec *spec.Spec, subnetSegments map[string][]ir.ID,
	config *ir.ConfigDefs) (*spec.Spec, map[ir.ID]*ir.SubnetSegmentDetails, error) {
	subnetsCache, ambiguousSubnets := inverseMapToFullyQualifiedName(config.Subnets)
	nifsCache, ambiguousNifs := inverseMapToFullyQualifiedName(config.NIFs)
	instancesCache, ambiguousInstances := inverseMapToFullyQualifiedName(config.Instances)
	vpesCache, ambiguousVpes := inverseMapToFullyQualifiedName(config.VPEs)

	// go over subnetSegments
	finalSubnetSegments := make(map[ir.ID]*ir.SubnetSegmentDetails)
	for segmentName, subnets := range subnetSegments {
		subnetsNames := make([]ir.ID, 0)
		VPCs := make([]ir.ID, 0)
		for _, subnet := range subnets {
			fullyQualifiedName, err := replaceResourceName(subnetsCache, ambiguousSubnets, subnet, spec.ResourceTypeSubnet)
			if err != nil {
				return nil, nil, err
			}
			subnetsNames = append(subnetsNames, fullyQualifiedName)
			VPCs = append(VPCs, config.Subnets[fullyQualifiedName].VPC)
		}
		finalSubnetSegments[segmentName] = &ir.SubnetSegmentDetails{Subnets: subnetsNames, OverlappingVPCs: ir.UniqueIDValues(VPCs)}
	}

	var err error
	// go over Spec
	for i := range jsonSpec.RequiredConnections {
		conn := &jsonSpec.RequiredConnections[i]
		fullyQualifiedSrc := conn.Src.Name
		switch conn.Src.Type {
		case spec.ResourceTypeSubnet:
			fullyQualifiedSrc, err = replaceResourceName(subnetsCache, ambiguousSubnets, conn.Src.Name, spec.ResourceTypeSubnet)
		case spec.ResourceTypeNif:
			fullyQualifiedSrc, err = replaceResourceName(nifsCache, ambiguousNifs, conn.Src.Name, spec.ResourceTypeNif)
		case spec.ResourceTypeInstance:
			fullyQualifiedSrc, err = replaceResourceName(instancesCache, ambiguousInstances, conn.Src.Name, spec.ResourceTypeInstance)
		case spec.ResourceTypeVpe:
			fullyQualifiedSrc, err = replaceResourceName(vpesCache, ambiguousVpes, conn.Src.Name, spec.ResourceTypeVpe)
		}
		if err != nil {
			return nil, nil, err
		}
		conn.Src.Name = fullyQualifiedSrc

		fullyQualifiedDst := conn.Dst.Name
		switch conn.Dst.Type {
		case spec.ResourceTypeSubnet:
			fullyQualifiedDst, err = replaceResourceName(subnetsCache, ambiguousSubnets, conn.Dst.Name, spec.ResourceTypeSubnet)
		case spec.ResourceTypeNif:
			fullyQualifiedDst, err = replaceResourceName(nifsCache, ambiguousNifs, conn.Dst.Name, spec.ResourceTypeNif)
		case spec.ResourceTypeInstance:
			fullyQualifiedDst, err = replaceResourceName(instancesCache, ambiguousInstances, conn.Dst.Name, spec.ResourceTypeInstance)
		case spec.ResourceTypeVpe:
			fullyQualifiedDst, err = replaceResourceName(vpesCache, ambiguousVpes, conn.Dst.Name, spec.ResourceTypeVpe)
		}
		if err != nil {
			return nil, nil, err
		}
		conn.Dst.Name = fullyQualifiedDst
	}
	return jsonSpec, finalSubnetSegments, nil
}

func translateConnection(defs *ir.Definitions, v *spec.SpecRequiredConnectionsElem, connectionIndex int) ([]ir.Connection, error) {
	p, err := translateProtocols(v.AllowedProtocols)
	if err != nil {
		return nil, err
	}
	srcResourceType, err := translateResourceType(v.Src.Type)
	if err != nil {
		return nil, err
	}
	src, err := defs.Lookup(srcResourceType, v.Src.Name)
	if err != nil {
		return nil, err
	}
	srcVPCs := defs.GetResourceOverlappingVPCs(srcResourceType, v.Src.Name)
	dstResourceType, err := translateResourceType(v.Dst.Type)
	if err != nil {
		return nil, err
	}
	dst, err := defs.Lookup(dstResourceType, v.Dst.Name)
	if err != nil {
		return nil, err
	}
	dstVPCs := defs.GetResourceOverlappingVPCs(dstResourceType, v.Dst.Name)
	err = defs.ValidateConnection(srcVPCs, dstVPCs)
	if err != nil {
		return nil, err
	}

	origin := connectionOrigin{
		connectionIndex: connectionIndex,
		srcName:         resourceName(v.Src),
		dstName:         resourceName(v.Dst),
	}
	out := ir.Connection{Src: src, Dst: dst, TrackedProtocols: p, Origin: origin}
	if v.Bidirectional {
		backOrigin := origin
		backOrigin.inverse = true
		in := ir.Connection{Src: dst, Dst: src, TrackedProtocols: p, Origin: &backOrigin}
		return []ir.Connection{out, in}, nil
	}
	return []ir.Connection{out}, nil
}

func translateProtocols(protocols spec.ProtocolList) ([]ir.TrackedProtocol, error) {
	var result = make([]ir.TrackedProtocol, len(protocols))
	for i, _p := range protocols {
		result[i].Origin = protocolOrigin{protocolIndex: i}
		switch p := _p.(type) {
		case spec.AnyProtocol:
			if len(protocols) != 1 {
				return nil, fmt.Errorf("when allowing any protocol, no more protocols can be defined")
			}
			result[i].Protocol = ir.AnyProtocol{}
		case spec.Icmp:
			if p.Type == nil {
				if p.Code != nil {
					return nil, fmt.Errorf("defining ICMP code for unspecified ICMP type is not allowed")
				}
				result[i].Protocol = ir.TrackedProtocol{Protocol: ir.ICMP{}}
			} else {
				err := ir.ValidateICMP(*p.Type, *p.Code)
				if err != nil {
					return nil, err
				}
				result[i].Protocol = ir.ICMP{ICMPCodeType: &ir.ICMPCodeType{Type: *p.Type, Code: p.Code}}
			}
		case spec.TcpUdp:
			result[i].Protocol = ir.TCPUDP{
				Protocol: ir.TransportLayerProtocolName(p.Protocol),
				PortRangePair: ir.PortRangePair{
					SrcPort: ir.PortRange{Min: p.MinSourcePort, Max: p.MaxSourcePort},
					DstPort: ir.PortRange{Min: p.MinDestinationPort, Max: p.MaxDestinationPort},
				},
			}
		default:
			return nil, fmt.Errorf("impossible protocol: %v", p)
		}
	}
	return result, nil
}

func translateResourceType(resourceType spec.ResourceType) (ir.ResourceType, error) {
	switch resourceType {
	case spec.ResourceTypeExternal:
		return ir.ResourceTypeExternal, nil
	case spec.ResourceTypeSegment:
		return ir.ResourceTypeSegment, nil
	case spec.ResourceTypeSubnet:
		return ir.ResourceTypeSubnet, nil
	case spec.ResourceTypeNif:
		return ir.ResourceTypeNIF, nil
	case spec.ResourceTypeInstance:
		return ir.ResourceTypeInstance, nil
	case spec.ResourceTypeVpe:
		return ir.ResourceTypeVPE, nil
	default:
		return ir.ResourceTypeSubnet, fmt.Errorf("unsupported resource type %v", resourceType)
	}
}

// unmarshal returns a Spec struct given a file adhering to spec_schema.input
func unmarshal(filename string) (*spec.Spec, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	jsonSpec := new(spec.Spec)
	err = json.Unmarshal(bytes, jsonSpec)
	if err != nil {
		return nil, err
	}
	for i := range jsonSpec.RequiredConnections {
		conn := &jsonSpec.RequiredConnections[i]
		if conn.AllowedProtocols == nil {
			conn.AllowedProtocols = spec.ProtocolList{spec.AnyProtocol{}}
		} else {
			for j := range conn.AllowedProtocols {
				p := conn.AllowedProtocols[j].(map[string]interface{})
				bytes, err = json.Marshal(p)
				if err != nil {
					return nil, err
				}
				switch p["protocol"] {
				case "ANY":
					var result spec.AnyProtocol
					err = json.Unmarshal(bytes, &result)
					conn.AllowedProtocols[j] = result
				case "TCP", "UDP":
					var result spec.TcpUdp
					err = json.Unmarshal(bytes, &result)
					conn.AllowedProtocols[j] = result
				case "ICMP":
					var result spec.Icmp
					err = json.Unmarshal(bytes, &result)
					conn.AllowedProtocols[j] = result
				default:
					return nil, fmt.Errorf("invalid protocol type %q", p["protocol"])
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return jsonSpec, err
}

func parseOverlappingVpcs(cidr *ipblock.IPBlock, vpcs map[ir.ID]*ir.VPCDetails) []ir.ID {
	result := make([]ir.ID, 0)
	for vpcName, vpcDetails := range vpcs {
		if vpcDetails.AddressPrefixes.Overlap(cidr) {
			result = append(result, vpcName)
		}
	}
	return result
}

func replaceResourceName(cache map[string]ir.ID, ambiguous map[string]struct{}, resourceName string,
	resourceType spec.ResourceType) (string, error) {
	if len(ir.ScopingComponents(resourceName)) != 1 {
		return resourceName, nil
	}
	if val, ok := cache[resourceName]; ok {
		return val, nil
	}
	if _, ok := ambiguous[resourceName]; ok {
		return "", fmt.Errorf("ambiguous resource name: %s", resourceName)
	}
	return "", fmt.Errorf("unknown resource name %s (resource type: %q)", resourceName, resourceType)
}

func inverseMapToFullyQualifiedName[T ir.Named](m map[ir.ID]T) (cache map[string]ir.ID, ambiguous map[string]struct{}) {
	ambiguous = make(map[string]struct{})
	cache = make(map[string]ir.ID)

	for nifName, nif := range m {
		if _, ok := ambiguous[nif.Name()]; ok {
			continue
		}
		if _, ok := cache[nif.Name()]; !ok {
			cache[nif.Name()] = nifName
		} else {
			delete(cache, nif.Name())
			ambiguous[nif.Name()] = struct{}{}
		}
	}
	return cache, ambiguous
}
