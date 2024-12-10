/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package jsonio handles global specification written in a JSON file, as described by spec_schema.input
package jsonio

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/np-guard/models/pkg/spec"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Reader implements ir.Reader
type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadSpec(filename string, configDefs *ir.ConfigDefs, isSG bool) (*ir.Spec, error) {
	jsonSpec, err := unmarshal(filename)
	if err != nil {
		return nil, err
	}
	defs, blocked, err := r.readDefinitions(jsonSpec, configDefs)
	if err != nil {
		return nil, err
	}

	// replace to fully qualified name
	jsonSpec, defs, err = replaceResourcesName(jsonSpec, defs)
	if err != nil {
		return nil, err
	}

	connections, err := r.translateConnections(jsonSpec.RequiredConnections, defs, blocked, isSG)
	if err != nil {
		return nil, err
	}

	return &ir.Spec{
		Connections:      connections,
		Defs:             defs,
		BlockedResources: blocked,
	}, nil
}

// replace all resources names to fully qualified name
func replaceResourcesName(jsonSpec *spec.Spec, defs *ir.Definitions) (*spec.Spec, *ir.Definitions, error) {
	config := defs.ConfigDefs

	// calculate distinct and ambiguous names for every endpoint type
	distinctSubnets, ambiguousSubnets := detectDistinctAndAmbiguousNames(config.Subnets)
	distinctNifs, ambiguousNifs := detectDistinctAndAmbiguousNames(config.NIFs)
	distinctInstances, ambiguousInstances := detectDistinctAndAmbiguousNames(config.Instances)
	distinctVpes, ambiguousVpes := detectDistinctAndAmbiguousNames(config.VPEs)

	// translate segments to fully qualified names
	nifSegments, err1 := replaceSegmentNames(defs.NifSegments, distinctNifs, ambiguousNifs, spec.ResourceType(spec.SegmentTypeNif))
	vpeSegments, err2 := replaceSegmentNames(defs.VpeSegments, distinctVpes, ambiguousVpes, spec.ResourceType(spec.SegmentTypeVpe))
	subnetSegments, err3 := replaceSegmentNames(defs.SubnetSegments, distinctSubnets, ambiguousSubnets,
		spec.ResourceType(spec.SegmentTypeSubnet))
	instanceSegments, err4 := replaceSegmentNames(defs.InstanceSegments, distinctInstances, ambiguousInstances,
		spec.ResourceType(spec.SegmentTypeInstance))
	if err := errors.Join(err1, err2, err3, err4); err != nil {
		return nil, nil, err
	}
	defs.SubnetSegments = subnetSegments
	defs.NifSegments = nifSegments
	defs.InstanceSegments = instanceSegments
	defs.VpeSegments = vpeSegments

	// translate connections resources to fully qualified names
	var err error
	for i := range jsonSpec.RequiredConnections {
		conn := &jsonSpec.RequiredConnections[i]
		fullyQualifiedSrc := conn.Src.Name
		switch conn.Src.Type {
		case spec.ResourceTypeSubnet:
			fullyQualifiedSrc, err = replaceResourceName(distinctSubnets, ambiguousSubnets, conn.Src.Name, spec.ResourceTypeSubnet)
		case spec.ResourceTypeNif:
			fullyQualifiedSrc, err = replaceResourceName(distinctNifs, ambiguousNifs, conn.Src.Name, spec.ResourceTypeNif)
		case spec.ResourceTypeInstance:
			fullyQualifiedSrc, err = replaceResourceName(distinctInstances, ambiguousInstances, conn.Src.Name, spec.ResourceTypeInstance)
		case spec.ResourceTypeVpe:
			fullyQualifiedSrc, err = replaceResourceName(distinctVpes, ambiguousVpes, conn.Src.Name, spec.ResourceTypeVpe)
		}
		if err != nil {
			return nil, nil, err
		}
		conn.Src.Name = fullyQualifiedSrc

		fullyQualifiedDst := conn.Dst.Name
		switch conn.Dst.Type {
		case spec.ResourceTypeSubnet:
			fullyQualifiedDst, err = replaceResourceName(distinctSubnets, ambiguousSubnets, conn.Dst.Name, spec.ResourceTypeSubnet)
		case spec.ResourceTypeNif:
			fullyQualifiedDst, err = replaceResourceName(distinctNifs, ambiguousNifs, conn.Dst.Name, spec.ResourceTypeNif)
		case spec.ResourceTypeInstance:
			fullyQualifiedDst, err = replaceResourceName(distinctInstances, ambiguousInstances, conn.Dst.Name, spec.ResourceTypeInstance)
		case spec.ResourceTypeVpe:
			fullyQualifiedDst, err = replaceResourceName(distinctVpes, ambiguousVpes, conn.Dst.Name, spec.ResourceTypeVpe)
		}
		if err != nil {
			return nil, nil, err
		}
		conn.Dst.Name = fullyQualifiedDst
	}
	return jsonSpec, defs, nil
}

func replaceSegmentNames(segments map[ir.ID]*ir.SegmentDetails, distinctNames map[string]ir.ID, ambiguousNames map[string]struct{},
	resourceType spec.ResourceType) (map[ir.ID]*ir.SegmentDetails, error) {
	for i, segmentDetails := range segments {
		for j, el := range segmentDetails.Elements {
			fullName, err := replaceResourceName(distinctNames, ambiguousNames, el, resourceType)
			if err != nil {
				return nil, err
			}
			segmentDetails.Elements[j] = fullName
		}
		segments[i] = segmentDetails
	}
	return segments, nil
}

func replaceResourceName(distinctNames map[string]ir.ID, ambiguousNames map[string]struct{}, resourceName string,
	resourceType spec.ResourceType) (string, error) {
	if val, ok := distinctNames[resourceName]; ok {
		return val, nil
	}
	if _, ok := ambiguousNames[resourceName]; ok {
		return "", fmt.Errorf("ambiguous resource name: %s", resourceName)
	}
	return "", fmt.Errorf("unknown resource name %s (resource type: %q)", resourceName, resourceType)
}

// detectDistinctAndAmbiguousNames returns two maps: one from a name to a fully qualified name,
// and the second is a set of ambiguous names
func detectDistinctAndAmbiguousNames[T any](m map[ir.ID]T) (distinctNames map[string]ir.ID, ambiguousNames map[string]struct{}) {
	ambiguousNames = make(map[string]struct{})
	distinctNames = make(map[string]ir.ID)

	for fullElementName := range m {
		distinctNames[fullElementName] = fullElementName
		for i, c := range fullElementName {
			if c == '/' && i+1 < len(fullElementName) {
				currName := fullElementName[i+1:]
				if _, ok := ambiguousNames[currName]; ok {
					continue
				}
				if _, ok := distinctNames[currName]; !ok {
					distinctNames[currName] = fullElementName
				} else {
					ambiguousNames[currName] = struct{}{}
					delete(distinctNames, currName)
				}
			}
		}
	}
	return distinctNames, ambiguousNames
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
