/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"slices"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func reduceCubesWithSGRemote(cubes *sgCubesPerProtocol) {
	deleteOtherProtocolIfAnyProtocolExistsSGRemote(cubes)
	compressThreeProtocolsToAnyProtocolSGRemote(cubes)
}

// delete other protocols rules if any protocol rule exists
func deleteOtherProtocolIfAnyProtocolExistsSGRemote(cubes *sgCubesPerProtocol) {
	for _, sgName := range cubes.anyProtocol {
		delete(cubes.tcp, sgName)
		delete(cubes.udp, sgName)
		delete(cubes.icmp, sgName)
	}
}

// merge tcp, udp and icmp rules into any protocol rule
func compressThreeProtocolsToAnyProtocolSGRemote(cubes *sgCubesPerProtocol) {
	for sgName, tcpPorts := range cubes.tcp {
		if udpPorts, ok := cubes.udp[sgName]; ok {
			if ic, ok := cubes.icmp[sgName]; ok {
				if ic.IsAll() && optimize.IsAllPorts(tcpPorts) && optimize.IsAllPorts(udpPorts) {
					delete(cubes.tcp, sgName)
					delete(cubes.udp, sgName)
					delete(cubes.icmp, sgName)
					cubes.anyProtocol = append(cubes.anyProtocol, sgName)
				}
			}
		}
	}
}

// observation: It pays to switch to all protocol rule when we have rules that cover all other protocols
// on exactly the same cidr (only one protocol can exceed).
//
//nolint:gocyclo // multiple checks
func reduceIPCubes(cubes *ipCubesPerProtocol) {
	tcpPtr := 0
	udpPtr := 0
	icmpPtr := 0

	for tcpPtr < len(cubes.tcp) && udpPtr < len(cubes.udp) && icmpPtr < len(cubes.icmp) {
		if !optimize.IsAllPorts(cubes.tcp[tcpPtr].Right) { // not all tcp ports
			tcpPtr++
			continue
		}
		if !optimize.IsAllPorts(cubes.udp[udpPtr].Right) { // not all udp ports
			udpPtr++
			continue
		}
		if !cubes.icmp[icmpPtr].Right.IsAll() { // not all icmp types & codes
			icmpPtr++
			continue
		}

		// all three protocols include all ports and types & codes
		// attempt to convert to any protocol rule
		if compressedToAnyProtocolCubeIPCubes(cubes, tcpPtr, udpPtr, icmpPtr) { // converted to any protocol rule
			continue
		}

		// could not compress to any protocol rule -- advance one ipblock
		// case 1: one protocol ipb contains two other ipbs ==> advance the smaller one
		// case 2: advance the smaller ipb
		tcpIPblock := cubes.tcp[tcpPtr].Left
		udpIPblock := cubes.udp[udpPtr].Left
		icmpIPblock := cubes.icmp[icmpPtr].Left

		udpComparedToICMP := udpIPblock.Compare(icmpIPblock)
		tcpComparedToUDP := tcpIPblock.Compare(udpIPblock)
		tcpComparedToICMP := tcpIPblock.Compare(icmpIPblock)

		udpAndICMPSubsetOfTCP := udpIPblock.IsSubset(tcpIPblock) && icmpIPblock.IsSubset(tcpIPblock)
		tcpAndICMPSubsetOfUDP := icmpIPblock.IsSubset(udpIPblock) && tcpIPblock.IsSubset(udpIPblock)
		tcpAndUDPSubsetOfICMP := tcpIPblock.IsSubset(icmpIPblock) && udpIPblock.IsSubset(icmpIPblock)

		switch {
		// case 1
		case udpAndICMPSubsetOfTCP && udpComparedToICMP == -1:
			udpPtr++
		case udpAndICMPSubsetOfTCP && udpComparedToICMP == 1:
			icmpPtr++
		case tcpAndICMPSubsetOfUDP && tcpComparedToICMP == -1:
			tcpPtr++
		case tcpAndICMPSubsetOfUDP && tcpComparedToICMP == 1:
			icmpPtr++
		case tcpAndUDPSubsetOfICMP && tcpComparedToUDP == -1:
			tcpPtr++
		case tcpAndUDPSubsetOfICMP && tcpComparedToUDP == 1:
			udpPtr++

		// case 2
		case tcpComparedToUDP == -1 && tcpComparedToICMP == -1:
			tcpPtr++
		case tcpComparedToUDP == 1 && udpComparedToICMP == -1:
			udpPtr++
		case tcpComparedToICMP == 1 && udpComparedToICMP == 1:
			icmpPtr++
		}
	}
}

// compress three protocol rules to any protocol rule (and maybe another protocol rule)
// returns true if the compression was successful
func compressedToAnyProtocolCubeIPCubes(cubes *ipCubesPerProtocol, tcpPtr, udpPtr, icmpPtr int) bool {
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
