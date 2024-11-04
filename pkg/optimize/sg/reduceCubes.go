/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"slices"

	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func reduceSGCubes(spans *sgCubesPerProtocol) {
	deleteOtherProtocolIfAllProtocolExists(spans)
	compressThreeProtocolsToAllProtocol(spans)
}

// delete other protocols rules if all protocol rule exists
func deleteOtherProtocolIfAllProtocolExists(spans *sgCubesPerProtocol) {
	for _, sgName := range spans.all {
		delete(spans.tcp, sgName)
		delete(spans.udp, sgName)
		delete(spans.icmp, sgName)
	}
}

// merge tcp, udp and icmp rules into all protocol rule
func compressThreeProtocolsToAllProtocol(spans *sgCubesPerProtocol) {
	for sgName, tcpPorts := range spans.tcp {
		if udpPorts, ok := spans.udp[sgName]; ok {
			if ic, ok := spans.icmp[sgName]; ok {
				if ic.IsAll() && tcpPorts.Equal(netset.AllPorts()) && udpPorts.Equal(netset.AllPorts()) {
					delete(spans.tcp, sgName)
					delete(spans.udp, sgName)
					delete(spans.icmp, sgName)
					spans.all = append(spans.all, sgName)
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
		if !cubes.tcp[tcpPtr].Right.Equal(netset.AllPorts()) {
			tcpPtr++
			continue
		}
		if !cubes.udp[udpPtr].Right.Equal(netset.AllPorts()) {
			udpPtr++
			continue
		}
		if !cubes.icmp[icmpPtr].Right.IsAll() {
			icmpPtr++
			continue
		}

		if compressedToAllCube(cubes, tcpPtr, udpPtr, icmpPtr) {
			continue
		}

		tcpIP := cubes.tcp[tcpPtr].Left
		udpIP := cubes.udp[udpPtr].Left
		icmpIP := cubes.icmp[icmpPtr].Left

		switch {
		// one protocol ipb contains two other ipbs ==> advance the smaller ipb
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

		// advance the smaller ipb
		case optimize.LessIPBlock(tcpIP, udpIP) && optimize.LessIPBlock(tcpIP, icmpIP):
			tcpPtr++
		case optimize.LessIPBlock(udpIP, tcpIP) && optimize.LessIPBlock(udpIP, icmpIP):
			udpPtr++
		case optimize.LessIPBlock(icmpIP, tcpIP) && optimize.LessIPBlock(icmpIP, udpIP):
			icmpPtr++
		}
	}
}

// compress three protocol rules to all protocol rule (and maybe another protocol rule)
func compressedToAllCube(cubes *ipCubesPerProtocol, tcpPtr, udpPtr, icmpPtr int) bool {
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
		cubes.all = cubes.all.Union(udpIP)
		return true
	case tcpIP.IsSubset(udpIP) && tcpIP.Equal(icmpIP):
		cubes.tcp = slices.Delete(cubes.tcp, tcpPtr, tcpPtr+1)
		cubes.icmp = slices.Delete(cubes.icmp, icmpPtr, icmpPtr+1)
		cubes.all = cubes.all.Union(tcpIP)
		return true
	case tcpIP.IsSubset(icmpIP) && tcpIP.Equal(udpIP):
		cubes.tcp = slices.Delete(cubes.tcp, tcpPtr, tcpPtr+1)
		cubes.udp = slices.Delete(cubes.udp, udpPtr, udpPtr+1)
		cubes.all = cubes.all.Union(tcpIP)
		return true
	}
	return false
}
