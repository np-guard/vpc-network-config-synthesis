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

// IPAddrs remote
func allSpanIPToSGRules(span *netset.IPBlock, direction ir.Direction) []*ir.SGRule {
	result := make([]*ir.SGRule, 0)
	for _, cidr := range ToCidrs(span) {
		result = append(result, ir.NewSGRule(direction, cidr, netp.AnyProtocol{}, netset.GetCidrAll(), ""))
	}
	return result
}

func tcpudpIPSpanToSGRules(span []ds.Pair[*netset.IPBlock, *interval.CanonicalSet], allSpan *netset.IPBlock,
	direction ir.Direction, isTCP bool) []*ir.SGRule {
	rules := []ds.Pair[*netset.IPBlock, *interval.Interval]{} // start ip and ports
	result := make([]*ir.SGRule, 0)

	for i := range span {
		if i > 0 {
			prevIPBlock := span[i-1].Left
			currIPBlock := span[i].Left
			if !touching(prevIPBlock, currIPBlock) { // the cidrs are not touching
				hole := IPBlockFromRange(NextIP(prevIPBlock), BeforeIP(currIPBlock))
				if !hole.IsSubset(allSpan) { // there in no all rule covering the hole
					for _, r := range rules {
						p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(r.Right.Start()), int(r.Right.End()))
						for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[i-1].Left))) {
							result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
						}
					}
					rules = []ds.Pair[*netset.IPBlock, *interval.Interval]{}
				}
			}
		}

		activePorts := interval.NewCanonicalSet()
		for _, r := range rules {
			if !r.Right.ToSet().IsSubset(span[i].Right) { // create rules
				p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(r.Right.Start()), int(r.Right.End()))
				for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[i-1].Left))) {
					result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
				}
			} else {
				activePorts.AddInterval(*r.Right)
			}
		}

		// new rules
		for _, ports := range span[i].Right.Intervals() {
			if !ports.ToSet().IsSubset(activePorts) { // it is not contained in other rules
				r := ds.Pair[*netset.IPBlock, *interval.Interval]{Left: FirstIPAddress(span[i].Left), Right: &ports}
				rules = append(rules, r)
			}
		}
	}

	// create the rest of the rules
	for _, r := range rules {
		p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(r.Right.Start()), int(r.Right.Start()))
		for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[len(span)-1].Left))) {
			result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
		}
	}

	return result
}

// problem: where should I end the rule?
// func createTcpudpRules(rules []ds.Pair[*netset.IPBlock, *interval.Interval], span []ds.Pair[*netset.IPBlock, *interval.CanonicalSet],
// 	direction ir.Direction, isTCP bool) (res []ir.SGRule) {
// 	res = make([]ir.SGRule, 0)
// 	for _, r := range rules {
// 		p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(r.Right.Start()), int(r.Right.Start()))
// 		for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[len(span)-1].Left))) {
// 			res = append(res, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
// 		}
// 	}
// 	return res
// }

func icmpSpanToSGRules(span []ds.Pair[*netset.IPBlock, *netset.ICMPSet], allSpan *netset.IPBlock, direction ir.Direction) []*ir.SGRule {
	rules := []ds.Pair[*netset.IPBlock, *netp.ICMP]{}
	result := make([]*ir.SGRule, 0)

	for i := range span {
		if i > 0 {
			prevIPBlock := span[i-1].Left
			currIPBlock := span[i].Left
			if !touching(prevIPBlock, currIPBlock) { // the cidrs are not touching
				hole := IPBlockFromRange(NextIP(prevIPBlock), BeforeIP(currIPBlock))
				if !hole.IsSubset(allSpan) { // there in no all rule covering the hole
					for _, r := range rules {
						p, _ := netp.NewICMP(r.Right.TypeCode)
						for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[i-1].Left))) {
							result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
						}
					}
					rules = []ds.Pair[*netset.IPBlock, *netp.ICMP]{}
				}
			}
		}

		activeICMP := netset.EmptyICMPSet()
		for _, r := range rules {
			ruleIcmpSet := netset.NewICMPSet(*r.Right)
			if !ruleIcmpSet.IsSubset(span[i].Right) { // create rules
				p, _ := netp.NewICMP(r.Right.TypeCode)
				for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[i-1].Left))) {
					result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
				}
			} else {
				activeICMP.Union(ruleIcmpSet)
			}
		}

		// new rules
		for _, p := range span[i].Right.Partitions() {
			if !netset.NewICMPSet(p).IsSubset(activeICMP) {
				r := ds.Pair[*netset.IPBlock, *netp.ICMP]{Left: FirstIPAddress(span[i].Left), Right: &p}
				rules = append(rules, r)
			}
		}
	}

	// create the rest of the rules
	for _, r := range rules {
		p, _ := netp.NewICMP(r.Right.TypeCode)
		for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[len(span)-1].Left))) {
			result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
		}
	}

	return result
}