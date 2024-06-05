/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
	"reflect"

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
	// *ipblock.IPBlock | SGName
}

type SGRule struct {
	Direction   Direction
	Remote      RemoteType
	Protocol    Protocol
	Explanation string
}

type SG struct {
	Rules    []SGRule
	Attached []SGName
}

type SGCollection struct {
	SGs map[ID]map[SGName]*SG
}

type SGWriter interface {
	WriteSG(*SGCollection, string) error
}

func (r *SGRule) isRedundant(rules []SGRule) bool {
	for i := range rules {
		if rules[i].mustSupersede(r) {
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
	return &SG{Rules: []SGRule{}, Attached: []SGName{}}
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
	if rule.isRedundant(a.Rules) {
		return
	}
	a.Rules = append(a.Rules, *rule)
}

func MergeSGCollections(collections ...*SGCollection) *SGCollection {
	result := NewSGCollection()
	for _, c := range collections {
		for _, vpc := range c.SGs {
			for sgName := range vpc {
				sg := c.LookupOrCreate(sgName)
				for r := range sg.Rules {
					result.LookupOrCreate(sgName).Add(&sg.Rules[r])
				}
			}
		}
	}
	return result
}

func (c *SGCollection) Write(w Writer, vpc string) error {
	return w.WriteSG(c, vpc)
}

func (c *SGCollection) SortedSGNames(vpc ID) []SGName {
	if vpc == "" {
		return utils.SortedKeys(c.SGs)
	}
	return utils.SortedValuesInKey(c.SGs, vpc)
}
