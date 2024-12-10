/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"slices"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func reduceSGCubes(spans *sgCubesPerProtocol) {
	deleteOtherProtocolIfAnyProtocolExists(spans)
	compressThreeProtocolsToAnyProtocol(spans)
}

// delete other protocols rules if any protocol rule exists
func deleteOtherProtocolIfAnyProtocolExists(spans *sgCubesPerProtocol) {
	for _, sgName := range spans.anyProtocol {
		delete(spans.tcp, sgName)
		delete(spans.udp, sgName)
		delete(spans.icmp, sgName)
	}
}

// merge tcp, udp and icmp rules into any protocol rule
func compressThreeProtocolsToAnyProtocol(spans *sgCubesPerProtocol) {
	for sgName, tcpPorts := range spans.tcp {
		if udpPorts, ok := spans.udp[sgName]; ok {
			if ic, ok := spans.icmp[sgName]; ok {
				if ic.IsAll() && optimize.AllPorts(tcpPorts) && optimize.AllPorts(udpPorts) {
					delete(spans.tcp, sgName)
					delete(spans.udp, sgName)
					delete(spans.icmp, sgName)
					spans.anyProtocol = append(spans.anyProtocol, sgName)
				}
			}
		}
	}
}

// observation: It pays to switch to all protocol rule when we have rules that cover all other protocols
// on exactly the same cidr (only one protocol can exceed).
//
//nolint:gocyclo // multiple if statments
func reduceIPCubes(cubes *ipCubesPerProtocol) {
	tcpPtr := 0
	udpPtr := 0
	icmpPtr := 0

	for tcpPtr < len(cubes.tcp) && udpPtr < len(cubes.udp) && icmpPtr < len(cubes.icmp) {
		if !optimize.AllPorts(cubes.tcp[tcpPtr].Right) { // not all tcp ports
			tcpPtr++
			continue
		}
		if !optimize.AllPorts(cubes.udp[udpPtr].Right) { // not all udp ports
			udpPtr++
			continue
		}
		if !cubes.icmp[icmpPtr].Right.IsAll() { // not all icmp types & codes
			icmpPtr++
			continue
		}

		// all three protocols include all ports and types & codes
		// attempt to convert to any protocol rule
		if compressedToAnyProtocolCube(cubes, tcpPtr, udpPtr, icmpPtr) { // converted to any protocol rule
			continue
		}

		// could not compress to any protocol rule -- advance one ipblock
		// case 1: one protocol ipb contains two other ipbs ==> advance the smaller one
		// case 2: advance the smaller ipb
		tcpIP := cubes.tcp[tcpPtr].Left
		udpIP := cubes.udp[udpPtr].Left
		icmpIP := cubes.icmp[icmpPtr].Left

		switch {
		// case 1
		case udpIP.IsSubset(tcpIP) && icmpIP.IsSubset(tcpIP) && optimize.LessIPBlock(udpIP, icmpIP):
			udpPtr++
		case udpIP.IsSubset(tcpIP) && icmpIP.IsSubset(tcpIP) && optimize.LessIPBlock(icmpIP, udpIP):
			icmpPtr++
		case tcpIP.IsSubset(udpIP) && icmpIP.IsSubset(udpIP) && optimize.LessIPBlock(tcpIP, icmpIP):
			tcpPtr++
		case tcpIP.IsSubset(udpIP) && icmpIP.IsSubset(udpIP) && optimize.LessIPBlock(icmpIP, tcpIP):
			icmpPtr++
		case tcpIP.IsSubset(icmpIP) && udpIP.IsSubset(icmpIP) && optimize.LessIPBlock(tcpIP, udpIP):
			tcpPtr++
		case tcpIP.IsSubset(icmpIP) && udpIP.IsSubset(icmpIP) && optimize.LessIPBlock(udpIP, tcpIP):
			udpPtr++

		// case 2
		case optimize.LessIPBlock(tcpIP, udpIP) && optimize.LessIPBlock(tcpIP, icmpIP):
			tcpPtr++
		case optimize.LessIPBlock(udpIP, tcpIP) && optimize.LessIPBlock(udpIP, icmpIP):
			udpPtr++
		case optimize.LessIPBlock(icmpIP, tcpIP) && optimize.LessIPBlock(icmpIP, udpIP):
			icmpPtr++
		}
	}
}

// compress three protocol rules to any protocol rule (and maybe another protocol rule)
// returns true if the compression was successful
func compressedToAnyProtocolCube(cubes *ipCubesPerProtocol, tcpPtr, udpPtr, icmpPtr int) bool {
	tcpIP := cubes.tcp[tcpPtr].Left
	udpIP := cubes.udp[udpPtr].Left
	icmpIP := cubes.icmp[icmpPtr].Left

	switch {
	case udpIP.Equal(tcpIP) && udpIP.Equal(icmpIP):
		cubes.tcp = slices.Delete(cubes.tcp, tcpPtr, tcpPtr+1)
		fallthrough
	case udpIP.IsSubset(tcpIP) && udpIP.Equal(icmpIP):
		cubes.udp = slices.Delete(cubes.udp, udpPtr, udpPtr+1)
		cubes.icmp = slices.Delete(cubes.icmp, icmpPtr, icmpPtr+1)
		cubes.anyProtocol = cubes.anyProtocol.Union(udpIP)
		return true
	case tcpIP.IsSubset(udpIP) && tcpIP.Equal(icmpIP):
		cubes.tcp = slices.Delete(cubes.tcp, tcpPtr, tcpPtr+1)
		cubes.icmp = slices.Delete(cubes.icmp, icmpPtr, icmpPtr+1)
		cubes.anyProtocol = cubes.anyProtocol.Union(tcpIP)
		return true
	case tcpIP.IsSubset(icmpIP) && tcpIP.Equal(udpIP):
		cubes.tcp = slices.Delete(cubes.tcp, tcpPtr, tcpPtr+1)
		cubes.udp = slices.Delete(cubes.udp, udpPtr, udpPtr+1)
		cubes.anyProtocol = cubes.anyProtocol.Union(tcpIP)
		return true
	}
	return false
}
