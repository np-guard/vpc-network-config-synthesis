/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type (
	Action string

	ACLRule struct {
		Action      Action
		Direction   Direction
		Source      *netset.IPBlock
		Destination *netset.IPBlock
		Protocol    netp.Protocol
		Explanation string
	}

	ACL struct {
		Name    string
		Subnets []string

		// Internal and External are used for synthesis
		Internal []*ACLRule
		External []*ACLRule

		// Inbound and Outbound are used for optimization
		Inbound  []*ACLRule
		Outbound []*ACLRule
	}

	ACLCollection struct {
		ACLs map[ID]map[string]*ACL
	}

	ACLWriter interface {
		WriteACL(aclColl *ACLCollection, vpc string, isSynth bool) error
	}
)

const (
	Allow Action = "allow"
	Deny  Action = "deny"
)

func (r *ACLRule) isRedundant(rules []*ACLRule) bool {
	for _, rule := range rules {
		if rule.mustSupersede(r) {
			return true
		}
	}
	return false
}

func (r *ACLRule) mustSupersede(other *ACLRule) bool {
	otherExplanation := other.Explanation
	other.Explanation = r.Explanation
	res := reflect.DeepEqual(r, other)
	other.Explanation = otherExplanation
	return res
}

func (r *ACLRule) Target() *netset.IPBlock {
	if r.Direction == Inbound {
		return r.Destination
	}
	return r.Source
}

func (a *ACL) Rules() []*ACLRule {
	if a.Internal == nil && a.External == nil { // optimization mode
		return slices.Concat(a.Inbound, a.Outbound)
	}
	rules := a.Internal
	if len(a.External) != 0 {
		rules = slices.Concat(rules, makeDenyInternal(), a.External)
	}
	return rules
}

func (a *ACL) AppendInternal(rule *ACLRule) {
	if !rule.isRedundant(a.Internal) {
		a.Internal = append(a.Internal, rule)
	}
}

func (a *ACL) AppendExternal(rule *ACLRule) {
	if !rule.isRedundant(a.External) {
		a.External = append(a.External, rule)
	}
}

func (a *ACL) AttachedSubnetsString() string {
	a.Subnets = slices.Compact(slices.Sorted(slices.Values(a.Subnets)))
	return strings.Join(a.Subnets, ", ")
}

func NewACLCollection() *ACLCollection {
	return &ACLCollection{ACLs: map[ID]map[string]*ACL{}}
}

func NewACL(aclName, subnetName string) *ACL {
	return &ACL{Name: aclName, Subnets: []string{subnetName}}
}

func NewACLRule(action Action, direction Direction, src, dst *netset.IPBlock, p netp.Protocol, e string) *ACLRule {
	return &ACLRule{Action: action, Direction: direction, Source: src, Destination: dst, Protocol: p, Explanation: e}
}

func aclSelector(subnetName ID, single bool) string {
	if single {
		return fmt.Sprintf("%s/singleACL", VpcFromScopedResource(subnetName))
	}
	return subnetName
}

func (c *ACLCollection) LookupOrCreate(subnetName string, singleACL bool) *ACL {
	vpcName := VpcFromScopedResource(subnetName)
	aclName := aclSelector(subnetName, singleACL)

	if acl, ok := c.ACLs[vpcName][aclName]; ok {
		if singleACL {
			acl.Subnets = append(acl.Subnets, subnetName)
		}
		return acl
	}
	if c.ACLs[vpcName] == nil {
		c.ACLs[vpcName] = make(map[string]*ACL)
	}
	c.ACLs[vpcName][aclName] = NewACL(aclName, subnetName)
	return c.ACLs[vpcName][aclName]
}

func (c *ACLCollection) VpcNames() []string {
	return utils.SortedMapKeys(c.ACLs)
}

func (c *ACLCollection) Write(w Writer, vpc string, isSynth bool) error {
	return w.WriteACL(c, vpc, isSynth)
}

func (c *ACLCollection) SortedACLNames(vpc string) []string {
	if vpc == "" {
		return utils.SortedAllInnerMapsKeys(c.ACLs)
	}
	return utils.SortedMapKeys(c.ACLs[vpc])
}
