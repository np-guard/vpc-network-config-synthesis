/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// SG remote
func tcpudpToSGSpanToSGRules(span map[*ir.SGName]*interval.CanonicalSet, direction ir.Direction, isTCP bool) []ir.SGRule {
	result := make([]ir.SGRule, 0)
	for sgName, intervals := range span {
		for _, dstPorts := range intervals.Intervals() {
			p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(dstPorts.Start()), int(dstPorts.End()))
			rule := ir.SGRule{
				Direction: direction,
				Remote:    sgName,
				Protocol:  p,
				Local:     netset.GetCidrAll(),
			}
			result = append(result, rule)
		}
	}
	return result
}

func icmpToSGSpanToSGRules(span map[*ir.SGName]*icmp, direction ir.Direction) []ir.SGRule {
	result := make([]ir.SGRule, 0)
	for sgName, icmp := range span {
		result = append(result, icmp.toSGRulestoSG(sgName)...)
	}
	return result
}

func protocolAllToSGSpanToSGRules(span []*ir.SGName, direction ir.Direction) []ir.SGRule {
	result := make([]ir.SGRule, len(span))
	for i, sgName := range span {
		result[i] = ir.SGRule{
			Direction: direction,
			Remote:    sgName,
			Protocol:  netp.AnyProtocol{},
			Local:     netset.GetCidrAll(),
		}
	}
	return result
}

// IPAddrs remote
func tcpudpToIPAddrsSpanToSGRules(span []ds.Pair[*netset.IPBlock, *interval.CanonicalSet],
	direction ir.Direction, isTCP bool) []ir.SGRule {
	result := make([]ir.SGRule, len(span))
	for i := range span {
		for _, dstPorts := range span[i].Right.Intervals() {
			p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(dstPorts.Start()), int(dstPorts.End()))
			rule := ir.SGRule{
				Direction: direction,
				Remote:    span[i].Left,
				Protocol:  p,
				Local:     netset.GetCidrAll(),
			}
			result[i] = rule
		}
	}
	return result
}
