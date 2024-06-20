/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"log"
	"sort"
	"strings"
)

const (
	warningUnspecifiedACL = "The following subnets do not have required connections; the generated ACL will block all traffic: "
	warningUnspecifiedSG  = "The following endpoints do not have required connections; the generated SGs will block all traffic: "
)

func (s *Spec) ComputeBlockedSubnets() []ID {
	var blockedSubnets []ID

	for subnet := range s.Defs.Subnets {
		if s.findResourceInConnections([]ID{subnet}, ResourceTypeSubnet) {
			continue
		}

		// subnet segments which include the subnet
		segments := []ID{}
		for segmentName, segmentDetails := range s.Defs.SubnetSegments {
			for _, s := range segmentDetails.Subnets {
				if subnet == s {
					segments = append(segments, segmentName)
					break
				}
			}
		}
		if s.findResourceInConnections(segments, ResourceTypeSubnet) {
			continue
		}

		// cidr segments which include the subnet
		cidrSegments := []ID{}
		for segmentName, cidrSegmentDetails := range s.Defs.CidrSegments {
			for _, cidrDetails := range cidrSegmentDetails.Cidrs {
				for _, s := range cidrDetails.ContainedSubnets {
					if subnet == s {
						cidrSegments = append(cidrSegments, segmentName)
						break
					}
				}
			}
		}
		if !s.findResourceInConnections(cidrSegments, ResourceTypeCidr) {
			blockedSubnets = append(blockedSubnets, subnet)
		}
	}
	sort.Strings(blockedSubnets)
	printUnspecifiedWarning(warningUnspecifiedACL, blockedSubnets)
	return blockedSubnets
}

func (s *Spec) ComputeBlockedResources() []ID {
	blockedResources := append(s.computeBlockedNIFs(), s.computeBlockedVPEs()...)
	sort.Strings(blockedResources)
	printUnspecifiedWarning(warningUnspecifiedSG, blockedResources)
	return blockedResources
}

func (s *Spec) computeBlockedVPEs() []ID {
	var blockedVPEs []ID
	for vpe := range s.Defs.VPEs {
		if !s.findResourceInConnections([]ID{vpe}, ResourceTypeVPE) {
			blockedVPEs = append(blockedVPEs, vpe)
		}
	}
	return blockedVPEs
}

func (s *Spec) computeBlockedNIFs() []ID {
	var blockedResources []ID
	for instanceName, instanceDetails := range s.Defs.Instances {
		if s.findResourceInConnections([]ID{instanceName}, ResourceTypeNIF) {
			continue
		}

		// instance is not in spec. look for its NIFs
		var blockedNIFs []ID
		for _, nif := range instanceDetails.Nifs {
			if !s.findResourceInConnections([]ID{nif}, ResourceTypeNIF) {
				blockedNIFs = append(blockedNIFs, nif)
			}
		}

		// instance has only one NIF which was not found
		if len(blockedNIFs) > 0 && len(instanceDetails.Nifs) == 1 {
			blockedResources = append(blockedResources, instanceName)
		} else {
			blockedResources = append(blockedResources, blockedNIFs...)
		}
	}
	return blockedResources
}

func (s *Spec) findResourceInConnections(resources []ID, epType ResourceType) bool {
	// The slice of IDs represents all resources that include the resource we are looking for
	for c := range s.Connections {
		for _, ep := range resources {
			if s.Connections[c].Src.Type == epType && ep == s.Connections[c].Src.Name {
				return true
			}
			if s.Connections[c].Dst.Type == epType && ep == s.Connections[c].Dst.Name {
				return true
			}
		}
	}
	return false
}

func printUnspecifiedWarning(warning string, blockedResources []ID) {
	if len(blockedResources) > 0 {
		log.Println(warning, strings.Join(blockedResources, ", "))
	}
}
