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

func reduceSGCubes(spans *sgCubesPerProtocol) *sgCubesPerProtocol {
	spans = deleteOtherProtocolIfAllProtocolExists(spans)
	return compressThreeProtocolsToAllProtocol(spans)
}

// delete other protocols rules if all protocol rule exists
func deleteOtherProtocolIfAllProtocolExists(spans *sgCubesPerProtocol) *sgCubesPerProtocol {
	for _, sgName := range spans.all {
		delete(spans.tcp, sgName)
		delete(spans.udp, sgName)
		delete(spans.icmp, sgName)
	}
	return spans
}

// merge tcp, udp and icmp rules into all protocol rule
func compressThreeProtocolsToAllProtocol(spans *sgCubesPerProtocol) *sgCubesPerProtocol {
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
	return spans
}

// observation: It pays to switch to all protocol rule when we have rules that cover all other protocols
// on exactly the same cidr (only one protocol can exceed).
func reduceIPCubes(cubes *ipCubesPerProtocol) *ipCubesPerProtocol {
	tcpPtr := 0
	udpPtr := 0
	icmpPtr := 0

	var changed bool
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

		cubes, changed = reduce(cubes, tcpPtr, udpPtr, icmpPtr)
		if changed {
			continue
		}

		tcpIP := cubes.tcp[tcpPtr].Left
		udpIP := cubes.udp[udpPtr].Left
		icmpIP := cubes.icmp[icmpPtr].Left

		if udpIP.IsSubset(tcpIP) && icmpIP.IsSubset(tcpIP) {
			if optimize.LessIPBlock(udpIP, icmpIP) {
				udpPtr++
			} else {
				icmpPtr++
			}
			continue
		}

	}
	return cubes
}

func reduce(cubes *ipCubesPerProtocol, tcpPtr, udpPtr, icmpPtr int) (*ipCubesPerProtocol, bool) {
	tcpIP := cubes.tcp[tcpPtr].Left
	udpIP := cubes.udp[udpPtr].Left
	icmpIP := cubes.icmp[icmpPtr].Left

	if udpIP.IsSubset(tcpIP) && udpIP.Equal(icmpIP) {
		cubes.all = cubes.all.Union(udpIP)
		slices.Delete(cubes.udp, udpPtr, udpPtr+1)
		slices.Delete(cubes.icmp, icmpPtr, icmpPtr+1)
		if tcpIP.Equal(udpIP) {
			slices.Delete(cubes.tcp, tcpPtr, tcpPtr+1)
		}
		// continue
	}

	if tcpIP.IsSubset(udpIP) && tcpIP.Equal(icmpIP) {
		cubes.all = cubes.all.Union(udpIP)
		slices.Delete(cubes.udp, udpPtr, udpPtr+1)
		slices.Delete(cubes.icmp, icmpPtr, icmpPtr+1)
		if tcpIP.Equal(udpIP) {
			slices.Delete(cubes.tcp, tcpPtr, tcpPtr+1)
		}
		// continue
	}

	return cubes, false
}
