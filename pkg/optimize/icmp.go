/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

type icmp struct {
	allowed map[int]bool // type x code X allow/deny
}

func newIcmp() *icmp {
	return &icmp{}
}

func (i *icmp) add(*netp.ICMPTypeCode) {

}

func (i *icmp) all() bool {
	return false
}

func (i *icmp) toSGRulestoSG(sgName *ir.SGName, direction ir.Direction) []ir.SGRule {
	return []ir.SGRule{}
}

func (i *icmp) toSGRulestoIPAddrs(ipAddrs *netset.IPBlock, direction ir.Direction) []ir.SGRule {
	return []ir.SGRule{}
}

// returns true if i is a subset of other
func (i *icmp) subset(other *icmp) bool {
	return true
}
