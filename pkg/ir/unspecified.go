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
		if s.findSubnetInConnections(subnet) {
			continue
		}

		segments := []string{} // segments which include the subnet
		for segmentName, segment := range s.Defs.SubnetSegments {
			for _, s := range segment {
				if subnet == s {
					segments = append(segments, segmentName)
					break
				}
			}
		}

		if !s.findSegmentInConnections(segments) {
			blockedSubnets = append(blockedSubnets, subnet)
		}
	}
	printUnspecifiedWarning(warning, blockedSubnets)
}

func (s *Spec) ComputeBlockedEndPoints() {
	var blockedEndPoints []string
	warning := warningUnspecifiedSG

	blockedEndPoints = append(blockedEndPoints, s.computeBlockedNIFs()...)
	blockedEndPoints = append(blockedEndPoints, s.computeBlockedVPEs()...)

	printUnspecifiedWarning(warning, blockedEndPoints)
}

func (s *Spec) computeBlockedVPEs() []string {
	var blockedVPEs []string
	for vpe := range s.Defs.VPEToIP {
		vpeFound := false
		for c := range s.Connections { // find the VPE in spec
			if s.Connections[c].Src.Type == EndpointTypeVPE && vpe == s.Connections[c].Src.Name {
				vpeFound = true
				break
			}
			if s.Connections[c].Dst.Type == EndpointTypeVPE && vpe == s.Connections[c].Dst.Name {
				vpeFound = true
				break
			}
		}
		if !vpeFound {
			blockedVPEs = append(blockedVPEs, vpe)
		}
	}
	return blockedVPEs
}

func (s *Spec) computeBlockedNIFs() []string {
	var blockedEndPoints []string

	for instance, NIFs := range s.Defs.InstanceToNIFs {
		if s.findInstanceInConnections(instance) {
			continue
		}

		// instance is not in spec. look for its NIFs
		blockedNifs := s.findNIFsInConnections(NIFs)

		if blockedNifs != nil && len(NIFs) == 1 { // instance has only one NIF which was not found
			blockedEndPoints = append(blockedEndPoints, instance)
		} else {
			blockedEndPoints = append(blockedEndPoints, blockedNifs...)
		}
	}
	return blockedEndPoints
}

func (s *Spec) findInstanceInConnections(instance string) bool {
	for c := range s.Connections {
		if s.Connections[c].Src.Type == EndpointTypeNIF && instance == s.Connections[c].Src.Name { // should be changed to EndpointTypeInstance
			return true
		}
		if s.Connections[c].Dst.Type == EndpointTypeNIF && instance == s.Connections[c].Dst.Name { // should be changed to EndpointTypeInstance
			return true
		}
	}
	return false
}

func (s *Spec) findNIFsInConnections(nifs []string) []string {
	var blockedNIFs []string
	for _, nif := range nifs {
		foundNif := false
		for c := range s.Connections {
			if s.Connections[c].Src.Type == EndpointTypeNIF && nif == s.Connections[c].Src.Name {
				foundNif = true
				break
			}
			if s.Connections[c].Dst.Type == EndpointTypeNIF && nif == s.Connections[c].Dst.Name {
				foundNif = true
				break
			}
		}
		if !foundNif {
			blockedNIFs = append(blockedNIFs, nif)
		}
	}
	return blockedNIFs
}

func (s *Spec) findSubnetInConnections(subnet string) bool {
	for c := range s.Connections {
		if s.Connections[c].Src.Type == EndpointTypeSubnet && subnet == s.Connections[c].Src.Name {
			return true
		}
		if s.Connections[c].Dst.Type == EndpointTypeSubnet && subnet == s.Connections[c].Dst.Name {
			return true
		}
	}
	return false
}

func (s *Spec) findSegmentInConnections(segments []string) bool {
	for c := range s.Connections {
		for _, i := range segments {
			if s.Connections[c].Src.Type == EndpointTypeSubnet && i == s.Connections[c].Src.Name { // should be changed to EndpointTypeSegment
				return true
			}
			if s.Connections[c].Dst.Type == EndpointTypeSubnet && i == s.Connections[c].Dst.Name { // should be changed to EndpointTypeSegment
				return true
			}
		}
	}
	return false
}

func printUnspecifiedWarning(warning string, blockedEndPoints []string) {
	if blockedEndPoints != nil {
		log.Println(warning, strings.Join(blockedEndPoints, ", "))
	}
}
