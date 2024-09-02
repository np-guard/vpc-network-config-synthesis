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
	allowed map[int]bool // type X code X allow/deny
}

func newIcmp() *icmp {
	return &icmp{}
}

func (i *icmp) add(*netp.ICMPTypeCode) {

}

func (i *icmp) allCodes(t int) bool {
	return true
}

func (i *icmp) allTypesAndCodes() bool {
	return true
}

func (i *icmp) toSGRulestoSG(sgName *ir.SGName) []ir.SGRule {
	return []ir.SGRule{}
}

func (i *icmp) toSGRulestoIPAddrs(ipAddrs *netset.IPBlock) []ir.SGRule {
	return []ir.SGRule{}
}

func (i *icmp) all() bool {
	return false
}
