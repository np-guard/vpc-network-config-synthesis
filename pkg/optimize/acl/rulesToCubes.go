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
		tcpAllow:  ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.PortSet](),
		tcpDeny:   ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.PortSet](),
		udpAllow:  ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.PortSet](),
		udpDeny:   ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.PortSet](),
		icmpAllow: ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet](),
		icmpDeny:  ds.NewLeftTripleSet[*netset.IPBlock, *netset.IPBlock, *netset.ICMPSet](),
	}

	for _, rule := range rules {
		switch p := rule.Protocol.(type) {
		case netp.TCPUDP:
			if p.ProtocolString() == "TCP" {
				tcpRuleToCubes(res, rule)
			} else {
				udpRuleToCubes(res, rule)
			}
		case netp.ICMP:
			icmpRuleToCubes(res, rule)
		case netp.AnyProtocol:
			anyProtocolRuleToCubes(res, rule)
		}
	}
	return res
}

func tcpRuleToCubes(cubes *aclCubesPerProtocol, rule *ir.ACLRule) {
	ruleCube := ds.CartesianLeftTriple(rule.Source, rule.Destination, rule.Protocol.(netp.TCPUDP).DstPorts().ToSet())
	if rule.Action == ir.Allow {
		r := ruleCube.Subtract(cubes.tcpDeny)
		cubes.tcpAllow = cubes.tcpAllow.Union(r)
	} else {
		r := ruleCube.Subtract(cubes.tcpAllow)
		cubes.tcpDeny = cubes.tcpDeny.Union(r)
	}
}

func udpRuleToCubes(cubes *aclCubesPerProtocol, rule *ir.ACLRule) {
	ruleCube := ds.CartesianLeftTriple(rule.Source, rule.Destination, rule.Protocol.(netp.TCPUDP).DstPorts().ToSet())
	if rule.Action == ir.Allow {
		r := ruleCube.Subtract(cubes.udpDeny)
		cubes.udpAllow = cubes.udpAllow.Union(r)
	} else {
		r := ruleCube.Subtract(cubes.udpAllow)
		cubes.udpDeny = cubes.udpDeny.Union(r)
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
	tcp, _ := netp.NewTCPUDP(true, netp.MinPort, netp.MaxPort, netp.MinPort, netp.MaxPort)
	tcpRuleToCubes(cubes, ir.NewACLRule(rule.Action, rule.Direction, rule.Source, rule.Destination, tcp, rule.Explanation))

	udp, _ := netp.NewTCPUDP(false, netp.MinPort, netp.MaxPort, netp.MinPort, netp.MaxPort)
	udpRuleToCubes(cubes, ir.NewACLRule(rule.Action, rule.Direction, rule.Source, rule.Destination, udp, rule.Explanation))

	icmp, _ := netp.NewICMPWithoutRFCValidation(nil) // all types and codes
	icmpRuleToCubes(cubes, ir.NewACLRule(rule.Action, rule.Direction, rule.Source, rule.Destination, icmp, rule.Explanation))
}
