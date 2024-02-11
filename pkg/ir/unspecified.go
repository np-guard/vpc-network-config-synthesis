package ir

import (
	"log"
)

func (s *Spec) ComputeBlockedEndPoints() {
	printWarning := false
	warningString := "The following endpoints do not have required connections; no SGs were generated for them: "

	instances := make([]string, 0, len(s.Defs.InstanceToNIFs))
	for instance := range s.Defs.InstanceToNIFs {
		instances = append(instances, instance)
	}
	warningString, printWarning = s.findEndPointsInConnections(instances, warningString, printWarning)

	vpes := make([]string, 0, len(s.Defs.VPEToIP))
	for vpe := range s.Defs.VPEToIP {
		vpes = append(vpes, vpe)
	}
	warningString, printWarning = s.findEndPointsInConnections(vpes, warningString, printWarning)

	if printWarning {
		log.Println(warningString + ".")
	}
}

func (s *Spec) ComputeBlockedSubnets(singleACL bool) {
	var printWarning = false
	var warningString string

	if singleACL {
		warningString = "The following subnets do not have required connections; the generated ACL will block all traffic: "
	} else {
		warningString = "The following subnets do not have required connections; no ACLs were generated for them: "
	}

	for subnet := range s.Defs.Subnets {
		includingSubnet := []string{subnet}
		for segmentName, segment := range s.Defs.SubnetSegments {
			for _, s := range segment {
				if subnet == s {
					includingSubnet = append(includingSubnet, segmentName)
					break
				}
			}
		}

		subnetFound := false

		for c := range s.Connections {
			if subnetFound {
				break
			}
			for _, i := range includingSubnet {
				if i == s.Connections[c].Src.Name || i == s.Connections[c].Dst.Name {
					subnetFound = true
					break
				}
			}
		}

		if !subnetFound {
			warningString, printWarning = s.updateWarningString(printWarning, warningString, subnet)
		}
	}
	if printWarning {
		log.Println(warningString + ".")
	}
}

func (s *Spec) findEndPointsInConnections(endpoints []string, warningString string, printWarning bool) (string, bool) {
	for _, endpoint := range endpoints {
		endpointFound := false
		for c := range s.Connections {
			if endpoint == s.Connections[c].Src.Name || endpoint == s.Connections[c].Dst.Name {
				endpointFound = true
				break
			}
		}
		if !endpointFound {
			warningString, printWarning = s.updateWarningString(printWarning, warningString, endpoint)
		}
	}
	return warningString, printWarning
}

func (s *Spec) updateWarningString(printWarning bool, warningString, endPoint string) (string, bool) {
	if !printWarning {
		warningString += endPoint
	} else {
		warningString = warningString + ", " + endPoint
	}
	return warningString, true
}
