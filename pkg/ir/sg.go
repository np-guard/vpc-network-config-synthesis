/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
	"reflect"

	"github.com/np-guard/models/pkg/netp"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type SGResource string

const (
	SGResourceNIF                  SGResource = "nif"
	SGResourceBareMetalNIF         SGResource = "bnif"
	SGResourceLoadBalancer         SGResource = "loadbalancer"
	SGResourceVPE                  SGResource = "vpe"
	SGResourceVPNServer            SGResource = "vpn"
	SGResourceFileShareMountTarget SGResource = "fsmt"
)

type SGName string

func (s SGName) String() string {
	return string(s)
}

type RemoteType interface {
	fmt.Stringer
	// *netset.IPBlock | SGName
}

type SGRule struct {
	Direction   Direction
	Remote      RemoteType
	Protocol    netp.Protocol
	Explanation string
}

type SG struct {
	InboundRules  []*SGRule
	OutboundRules []*SGRule
	Attached      []ID
}

type SGCollection struct {
	SGs map[ID]map[SGName]*SG
}

type SGWriter interface {
	WriteSG(sgColl *SGCollection, vpc string) error
}

func (r *SGRule) isRedundant(rules []*SGRule) bool {
	for _, rule := range rules {
		if rule.mustSupersede(r) {
			return true
		}
	}
	return false
}

func (r *SGRule) mustSupersede(other *SGRule) bool {
	otherExplanation := other.Explanation
	other.Explanation = r.Explanation
	res := reflect.DeepEqual(r, other)
	other.Explanation = otherExplanation
	return res
}

func NewSG() *SG {
	return &SG{InboundRules: []*SGRule{}, OutboundRules: []*SGRule{}, Attached: []ID{}}
}

func NewSGCollection() *SGCollection {
	return &SGCollection{SGs: map[ID]map[SGName]*SG{}}
}

func (c *SGCollection) LookupOrCreate(name SGName) *SG {
	vpcName := VpcFromScopedResource(string(name))
	if sg, ok := c.SGs[vpcName][name]; ok {
		return sg
	}
	newSG := NewSG()
	if c.SGs[vpcName] == nil {
		c.SGs[vpcName] = make(map[SGName]*SG)
	}
	c.SGs[vpcName][name] = newSG
	return newSG
}

func (a *SG) Add(rule *SGRule) {
	if rule.Direction == Outbound && !rule.isRedundant(a.OutboundRules) {
		a.OutboundRules = append(a.OutboundRules, rule)
	}
	if rule.Direction == Inbound && !rule.isRedundant(a.InboundRules) {
		a.InboundRules = append(a.InboundRules, rule)
	}
}

func (a *SG) AllRules() []*SGRule {
	return append(a.InboundRules, a.OutboundRules...)
}

func (c *SGCollection) Write(w Writer, vpc string) error {
	return w.WriteSG(c, vpc)
}

func (c *SGCollection) SortedSGNames(vpc ID) []SGName {
	if vpc == "" {
		return utils.SortedAllInnerMapsKeys(c.SGs)
	}
	return utils.SortedMapKeys(c.SGs[vpc])
}
