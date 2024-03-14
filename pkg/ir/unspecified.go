package ir

import (
	"log"
	"strings"
)

const (
	commonWarningACL            = "The following subnets do not have required connections; "
	warningUnspecifiedACL       = commonWarningACL + "no ACLs were generated for them: "
	warningUnspecifiedSingleACL = commonWarningACL + "the generated ACL will block all traffic: "
	warningUnspecifiedSG        = "The following resources do not have required connections; no SGs were generated for them: "
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
		if s.findResourceInConnections([]string{subnet}, ResourceTypeSubnet) {
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
			for _, subnets := range cidrSegment {
				for _, s := range subnets {
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
	printUnspecifiedWarning(warning, blockedSubnets)
}

func (s *Spec) ComputeBlockedResources() {
	warning := warningUnspecifiedSG

	blockedResources := append(s.computeBlockedNIFs(), s.computeBlockedVPEs()...)

	printUnspecifiedWarning(warning, blockedResources)
}

func (s *Spec) computeBlockedVPEs() []string {
	var blockedVPEs []string
	for vpe := range s.Defs.VPEToIP {
		if !s.findResourceInConnections([]string{vpe}, ResourceTypeVPE) {
			blockedVPEs = append(blockedVPEs, vpe)
		}
	}
	return blockedVPEs
}

func (s *Spec) computeBlockedNIFs() []string {
	var blockedResources []string

	for instance, NIFs := range s.Defs.InstanceToNIFs {
		if s.findResourceInConnections([]string{instance}, ResourceTypeNIF) {
			continue
		}

		// instance is not in spec. look for its NIFs
		var blockedNIFs []string
		for _, nif := range NIFs {
			if !s.findResourceInConnections([]string{nif}, ResourceTypeNIF) {
				blockedNIFs = append(blockedNIFs, nif)
			}
		}

		// instance has only one NIF which was not found
		if len(blockedNIFs) > 0 && len(NIFs) == 1 {
			blockedResources = append(blockedResources, instance)
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
