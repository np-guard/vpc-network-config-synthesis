/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sgoptimizer_test

import (
	"log"
	"testing"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/optimize"
)

func TestOps(t *testing.T) {
	sgRules := make([]*ir.SGRule, 0)
	p1, _ := netp.NewTCPUDP(true, 1, 100, 1, 10)
	p2, _ := netp.NewTCPUDP(true, 1, 100, 1, 20)

	ipb1, _ := netset.IPBlockFromCidrOrAddress("0.0.0.0")
	ipb2, _ := netset.IPBlockFromCidrOrAddress("0.0.0.0/31")

	sgRules = append(sgRules, ir.NewSGRule(ir.Outbound, ipb1, p1, netset.GetCidrAll(), ""))
	sgRules = append(sgRules, ir.NewSGRule(ir.Outbound, ipb2, p2, netset.GetCidrAll(), ""))

	res := tcpudpRulesToIPCubes(sgRules)
	for i, pair := range res {
		log.Println("pair ", i, ": ipblock: ", pair.Left.String(), ", ports: ", pair.Right.String())
	}

	sgRules = []*ir.SGRule{}
	sgRules = append(sgRules, ir.NewSGRule(ir.Outbound, ipb2, p1, netset.GetCidrAll(), ""))
	sgRules = append(sgRules, ir.NewSGRule(ir.Outbound, ipb1, p2, netset.GetCidrAll(), ""))

	res = tcpudpRulesToIPCubes(sgRules)
	for i, pair := range res {
		log.Println("pair ", i, ": ipblock: ", pair.Left.String(), ", ports: ", pair.Right.String())
	}

	t.Log("Hi")
}

func tcpudpRulesToIPCubes(rules []*ir.SGRule) []ds.Pair[*netset.IPBlock, *netset.PortSet] {
	cubes := ds.NewProductLeft[*netset.IPBlock, *netset.PortSet]()
	for _, rule := range rules {
		ipb := rule.Remote.(*netset.IPBlock)    // already checked
		portsSet := rule.Protocol.(netp.TCPUDP) // already checked
		r := ds.CartesianPairLeft(ipb, portsSet.DstPorts().ToSet())
		cubes = cubes.Union(r).(*ds.ProductLeft[*netset.IPBlock, *netset.PortSet])
	}
	return optimize.SortPartitionsByIPAddrs(cubes.Partitions())
}
