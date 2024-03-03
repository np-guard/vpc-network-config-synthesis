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
		if s.findEndPointInConnections([]string{subnet}, EndpointTypeSubnet) {
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
		if s.findEndPointInConnections(segments, EndpointTypeSubnet) {
			continue
		}

		// cidr segments which include the subnet
		cidrSements := []string{}
		for segmentName, cidrSegment := range s.Defs.CidrSegments {
			for _, subnets := range cidrSegment {
				for _, s := range subnets {
					if subnet == s {
						cidrSements = append(cidrSements, segmentName)
						break
					}
				}
			}
		}
		if !s.findEndPointInConnections(cidrSements, EndpointTypeCidr) {
			blockedSubnets = append(blockedSubnets, subnet)
		}
	}
	printUnspecifiedWarning(warning, blockedSubnets)
}

func (s *Spec) ComputeBlockedEndPoints() {
	warning := warningUnspecifiedSG

	blockedEndPoints := s.computeBlockedNIFs()
	blockedEndPoints = append(blockedEndPoints, s.computeBlockedVPEs()...)

	printUnspecifiedWarning(warning, blockedEndPoints)
}

func (s *Spec) computeBlockedVPEs() []string {
	var blockedVPEs []string
	for vpe := range s.Defs.VPEToIP {
		if !s.findEndPointInConnections([]string{vpe}, EndpointTypeVPE) {
			blockedVPEs = append(blockedVPEs, vpe)
		}
	}
	return blockedVPEs
}

func (s *Spec) computeBlockedNIFs() []string {
	var blockedEndPoints []string

	for instance, NIFs := range s.Defs.InstanceToNIFs {
		if s.findEndPointInConnections([]string{instance}, EndpointTypeNIF) {
			continue
		}

		// instance is not in spec. look for its NIFs
		var blockedNIFs []string
		for _, nif := range NIFs {
			if !s.findEndPointInConnections([]string{nif}, EndpointTypeNIF) {
				blockedNIFs = append(blockedNIFs, nif)
			}
		}

		// instance has only one NIF which was not found
		if len(blockedNIFs) > 0 && len(NIFs) == 1 {
			blockedEndPoints = append(blockedEndPoints, instance)
		} else {
			blockedEndPoints = append(blockedEndPoints, blockedNIFs...)
		}
	}
	return blockedEndPoints
}

func (s *Spec) findEndPointInConnections(endPoints []string, epType EndpointType) bool {
	// The slice of strings represents all endpoints that include the endpoint we are looking for
	for c := range s.Connections {
		for _, ep := range endPoints {
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

func printUnspecifiedWarning(warning string, blockedEndPoints []string) {
	if len(blockedEndPoints) > 0 {
		log.Println(warning, strings.Join(blockedEndPoints, ", "))
	}
}
