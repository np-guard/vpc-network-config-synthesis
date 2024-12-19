/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package acloptimizer

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func aclRulesToCubes(rules []*ir.ACLRule) *aclCubesPerProtocol {
	res := &aclCubesPerProtocol{
		tcpudpAllow: ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet](),
		tcpudpDeny:  ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.TCPUDPSet](),
		icmpAllow:   ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet](),
		icmpDeny:    ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet](),
	}

	for _, rule := range rules {
		switch rule.Protocol.(type) {
		case netp.TCPUDP:
			tcpudpRuleToCubes(res, rule)
		case netp.ICMP:
			icmpRuleToCubes(res, rule)
		case netp.AnyProtocol:
			anyProtocolRuleToCubes(res, rule)
		}
	}
	return res
}

func tcpudpRuleToCubes(cubes *aclCubesPerProtocol, rule *ir.ACLRule) {
	tcpudp := rule.Protocol.(netp.TCPUDP)
	tcpudpSrcPorts := tcpudp.SrcPorts()
	tcpudpDstPorts := tcpudp.DstPorts()
	tcpudpSet := netset.NewTCPorUDPSet(tcpudp.ProtocolString(), tcpudpSrcPorts.Start(), tcpudpSrcPorts.End(), tcpudpDstPorts.Start(),
		tcpudpDstPorts.End())
	ruleCube := ds.CartesianLeftTriple(rule.Source, rule.Destination, tcpudpSet)
	if rule.Action == ir.Allow {
		r := ruleCube.Subtract(cubes.tcpudpDeny)
		cubes.tcpudpAllow = cubes.tcpudpAllow.Union(r)
	} else {
		r := ruleCube.Subtract(cubes.tcpudpAllow)
		cubes.tcpudpDeny = cubes.tcpudpDeny.Union(r)
	}
}

func icmpRuleToCubes(cubes *aclCubesPerProtocol, rule *ir.ACLRule) {
	ruleCube := ds.CartesianLeftTriple(rule.Source, rule.Destination, optimize.IcmpToIcmpSet(rule.Protocol.(netp.ICMP)))
	if rule.Action == ir.Allow {
		r := ruleCube.Subtract(cubes.icmpDeny)
		cubes.icmpAllow = cubes.icmpAllow.Union(r)
	} else {
		r := ruleCube.Subtract(cubes.icmpAllow)
		cubes.icmpDeny = cubes.icmpDeny.Union(r)
	}
}

func anyProtocolRuleToCubes(cubes *aclCubesPerProtocol, rule *ir.ACLRule) {
	tcp, _ := netp.NewTCPUDP(true, netp.MinPort, netp.MaxPort, netp.MinPort, netp.MaxPort) // all ports
	tcpudpRuleToCubes(cubes, ir.NewACLRule(rule.Action, rule.Direction, rule.Source, rule.Destination, tcp, rule.Explanation))

	udp, _ := netp.NewTCPUDP(false, netp.MinPort, netp.MaxPort, netp.MinPort, netp.MaxPort) // all ports
	tcpudpRuleToCubes(cubes, ir.NewACLRule(rule.Action, rule.Direction, rule.Source, rule.Destination, udp, rule.Explanation))

	icmp, _ := netp.NewICMPWithoutRFCValidation(nil) // all types and codes
	icmpRuleToCubes(cubes, ir.NewACLRule(rule.Action, rule.Direction, rule.Source, rule.Destination, icmp, rule.Explanation))
}
