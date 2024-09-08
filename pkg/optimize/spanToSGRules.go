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
func tcpudpSGSpanToSGRules(span map[*ir.SGName]*interval.CanonicalSet, direction ir.Direction, isTCP bool) []ir.SGRule {
	result := make([]ir.SGRule, 0)
	for sgName, intervals := range span {
		for _, dstPorts := range intervals.Intervals() {
			p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(dstPorts.Start()), int(dstPorts.End()))
			result = append(result, ir.NewSGRule(direction, sgName, p, netset.GetCidrAll(), ""))
		}
	}
	return result
}

func icmpSGSpanToSGRules(span map[*ir.SGName]*netset.ICMPSet, direction ir.Direction) []ir.SGRule {
	result := make([]ir.SGRule, 0)
	for sgName, icmpSet := range span {
		for _, icmp := range icmpSet.Partitions() {
			p, _ := netp.NewICMP(icmp.TypeCode)
			result = append(result, ir.NewSGRule(direction, sgName, p, netset.GetCidrAll(), ""))
		}
	}
	return result
}

func protocolAllSGSpanToSGRules(span []*ir.SGName, direction ir.Direction) []ir.SGRule {
	result := make([]ir.SGRule, len(span))
	for i, sgName := range span {
		result[i] = ir.NewSGRule(direction, sgName, netp.AnyProtocol{}, netset.GetCidrAll(), "")
	}
	return result
}

// IPAddrs remote
func tcpudpIPSpanToSGRules(span []ds.Pair[*netset.IPBlock, *interval.CanonicalSet], _ []*netset.IPBlock,
	direction ir.Direction, isTCP bool) []ir.SGRule {
	rules := []ds.Pair[*netset.IPBlock, *interval.Interval]{} // start ip and ports
	result := make([]ir.SGRule, 0)

	for i := range span {
		if i > 0 && !touching(span[i-1].Left, span[i].Left) { // if the CIDRS are not touching
			for _, r := range rules {
				p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(r.Right.Start()), int(r.Right.End()))
				for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[i-1].Left))) {
					result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
				}
			}

			rules = []ds.Pair[*netset.IPBlock, *interval.Interval]{}
			continue
		}

		activePorts := interval.NewCanonicalSet()
		for _, r := range rules {
			if !r.Right.ToSet().IsSubset(span[i].Right) { // close old rules
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

	// close all old rules
	for _, r := range rules {
		p, _ := netp.NewTCPUDP(isTCP, netp.MinPort, netp.MaxPort, int(r.Right.Start()), int(r.Right.Start()))
		for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[len(span)-1].Left))) {
			result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
		}
	}

	return result
}

func icmpSpanToSGRules(span []ds.Pair[*netset.IPBlock, *netset.ICMPSet], _ []*netset.IPBlock, direction ir.Direction) []ir.SGRule {
	rules := []ds.Pair[*netset.IPBlock, *netp.ICMP]{}
	result := make([]ir.SGRule, 0)

	for i := range span {
		if i > 0 && !touching(span[i-1].Left, span[i].Left) { // if the CIDRS are not touching
			for _, r := range rules {
				p, _ := netp.NewICMP(r.Right.TypeCode)
				for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[i-1].Left))) {
					result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
				}
			}

			rules = []ds.Pair[*netset.IPBlock, *netp.ICMP]{}
			continue
		}

		activeICMP := netset.EmptyICMPSet()
		for _, r := range rules {
			ruleIcmpSet := netset.NewICMPSet(*r.Right)
			if !ruleIcmpSet.IsSubset(span[i].Right) { // close old rules
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

	// close all rules
	for _, r := range rules {
		p, _ := netp.NewICMP(r.Right.TypeCode)
		for _, cidr := range ToCidrs(IPBlockFromRange(r.Left, LastIPAddress(span[len(span)-1].Left))) {
			result = append(result, ir.NewSGRule(direction, cidr, p, netset.GetCidrAll(), ""))
		}
	}

	return result
}

func allSpanToSGRules(span []*netset.IPBlock, direction ir.Direction) []ir.SGRule {
	result := make([]ir.SGRule, 0)
	for _, ipAddrs := range span {
		for _, cidr := range ToCidrs(ipAddrs) {
			result = append(result, ir.NewSGRule(direction, cidr, netp.AnyProtocol{}, netset.GetCidrAll(), ""))
		}
	}
	return result
}
