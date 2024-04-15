/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"log"
	"strings"
)

const (
	commonWarningACL            = "The following subnets do not have required connections; "
	warningUnspecifiedACL       = commonWarningACL + "no ACLs were generated for them: "
	warningUnspecifiedSingleACL = commonWarningACL + "the generated ACL will block all traffic: "
	warningUnspecifiedSG        = "The following endpoints do not have required connections; no SGs were generated for them: "
)

func (s *Spec) ComputeBlockedSubnets(singleACL bool) {
	var warning string
	if singleACL {
		warning = warningUnspecifiedSingleACL
	} else {
		warning = warningUnspecifiedACL
	}
	var blockedSubnets []string

	for subnet := range s.Defs.Subnets {
		if s.findResourceInConnections([]string{string(subnet)}, ResourceTypeSubnet) {
			continue
		}

		// subnet segments which include the subnet
		segments := []string{}
		for segmentName, segment := range s.Defs.SubnetSegments {
			for _, s := range segment {
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
		cidrSegments := []string{}
		for segmentName, cidrSegment := range s.Defs.CidrSegments {
			for _, cidrDetails := range cidrSegment {
				for _, s := range cidrDetails.ContainedSubnets {
					if subnet == s {
						cidrSegments = append(cidrSegments, segmentName)
						break
					}
				}
			}
		}
		if !s.findResourceInConnections(cidrSegments, ResourceTypeCidr) {
			blockedSubnets = append(blockedSubnets, string(subnet))
		}
	}
	printUnspecifiedWarning(warning, blockedSubnets)
}

func (s *Spec) ComputeBlockedResources() {
	warning := warningUnspecifiedSG

	blockedResources := append(s.computeBlockedNIFs(), s.computeBlockedVPEs()...)

	printUnspecifiedWarning(warning, blockedResources)
}

func (s *Spec) computeBlockedVPEs() []string {
	var blockedVPEs []string
	for vpe := range s.Defs.VPEEndpoints {
		if !s.findResourceInConnections([]string{string(vpe)}, ResourceTypeVPE) {
			blockedVPEs = append(blockedVPEs, string(vpe))
		}
	}
	return blockedVPEs
}

func (s *Spec) computeBlockedNIFs() []string {
	var blockedResources []string
	for instanceName, instanceDetails := range s.Defs.Instances {
		if s.findResourceInConnections([]string{string(instanceName)}, ResourceTypeNIF) {
			continue
		}

		// instance is not in spec. look for its NIFs
		var blockedNIFs []string
		for _, nif := range instanceDetails.Nifs {
			if !s.findResourceInConnections([]string{string(nif)}, ResourceTypeNIF) {
				blockedNIFs = append(blockedNIFs, string(nif))
			}
		}

		// instance has only one NIF which was not found
		if len(blockedNIFs) > 0 && len(instanceDetails.Nifs) == 1 {
			blockedResources = append(blockedResources, string(instanceName))
		} else {
			blockedResources = append(blockedResources, blockedNIFs...)
		}
	}
	return blockedResources
}

func (s *Spec) findResourceInConnections(resources []string, epType ResourceType) bool {
	// The slice of strings represents all resources that include the resource we are looking for
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

func printUnspecifiedWarning(warning string, blockedResources []string) {
	if len(blockedResources) > 0 {
		log.Println(warning, strings.Join(blockedResources, ", "))
	}
}
