/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer

import (
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

// cubes (SGName X portSet) to SG rules
func tcpudpSGCubesToRules(cubes map[ir.SGName]*netset.PortSet, direction ir.Direction, isTCP bool, l *netset.IPBlock) []*ir.SGRule {
	result := make([]*ir.SGRule, 0)
	for sgName, portSet := range cubes {
		for _, dstPorts := range portSet.Intervals() {
			p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(dstPorts.Start()), int(dstPorts.End()))
			result = append(result, ir.NewSGRule(direction, sgName, p, l, ""))
		}
	}
	return result
}

// cubes (SGName X icmpset) to SG rules
func icmpSGCubesToRules(cubes map[ir.SGName]*netset.ICMPSet, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	result := make([]*ir.SGRule, 0)
	for sgName, icmpSet := range cubes {
		for _, icmp := range optimize.IcmpsetPartitions(icmpSet) {
			result = append(result, ir.NewSGRule(direction, sgName, icmp, l, ""))
		}
	}
	return result
}

// cubes (slice of SGs) to SG rules
func anyPotocolCubesToRules(span []ir.SGName, direction ir.Direction, l *netset.IPBlock) []*ir.SGRule {
	result := make([]*ir.SGRule, len(span))
	for i, sgName := range span {
		result[i] = ir.NewSGRule(direction, sgName, netp.AnyProtocol{}, l, "")
	}
	return result
}
