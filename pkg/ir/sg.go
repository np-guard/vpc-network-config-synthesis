/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type (
	SGName string

	RemoteType interface {
		fmt.Stringer
		// *netset.IPBlock | SGName
	}

	SGRule struct {
		Direction   Direction
		Remote      RemoteType
		Protocol    netp.Protocol
		Local       *netset.IPBlock
		Explanation string
	}

	SG struct {
		SGName        SGName
		InboundRules  map[string][]*SGRule // the key is the locals value
		OutboundRules map[string][]*SGRule // the key is the locals value
		Targets       []ID
	}

	SGCollection struct {
		SGs map[ID]map[SGName]*SG
	}

	SGWriter interface {
		WriteSG(sgColl *SGCollection, vpc string, isSynth bool) error
	}
)

func (s SGName) String() string {
	return string(s)
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

func NewSGRule(direction Direction, remote RemoteType, p netp.Protocol, local *netset.IPBlock, e string) *SGRule {
	return &SGRule{Direction: direction, Remote: remote, Protocol: p,
		Local: local, Explanation: e}
}

func NewSG(sgName SGName) *SG {
	return &SG{SGName: sgName,
		InboundRules:  make(map[string][]*SGRule),
		OutboundRules: make(map[string][]*SGRule),
	}
}

func NewSGCollection() *SGCollection {
	return &SGCollection{SGs: map[ID]map[SGName]*SG{}}
}

func (c *SGCollection) LookupOrCreate(name SGName) *SG {
	vpcName := VpcFromScopedResource(string(name))
	if sg, ok := c.SGs[vpcName][name]; ok {
		return sg
	}
	newSG := NewSG(name)
	if c.SGs[vpcName] == nil {
		c.SGs[vpcName] = make(map[SGName]*SG)
	}
	c.SGs[vpcName][name] = newSG
	return newSG
}

func (a *SG) Add(rule *SGRule) {
	local := rule.Local.String()
	if rule.Direction == Outbound && !rule.isRedundant(a.OutboundRules[local]) {
		a.OutboundRules[local] = append(a.OutboundRules[local], rule)
	}

	if rule.Direction == Inbound && !rule.isRedundant(a.InboundRules[local]) {
		a.InboundRules[local] = append(a.InboundRules[local], rule)
	}
}

func (a *SG) AllRules() []*SGRule {
	res := make([]*SGRule, 0)
	for _, key := range utils.SortedMapKeys(a.InboundRules) {
		res = slices.Concat(res, a.InboundRules[key])
	}
	for _, key := range utils.SortedMapKeys(a.OutboundRules) {
		res = slices.Concat(res, a.OutboundRules[key])
	}
	return res
}

func (c *SGCollection) VpcNames() []string {
	return utils.SortedMapKeys(c.SGs)
}

func (c *SGCollection) Write(w Writer, vpc string, isSynth bool) error {
	return w.WriteSG(c, vpc, isSynth)
}

func (c *SGCollection) SortedSGNames(vpc ID) []SGName {
	if vpc == "" {
		return utils.SortedAllInnerMapsKeys(c.SGs)
	}
	return utils.SortedMapKeys(c.SGs[vpc])
}
