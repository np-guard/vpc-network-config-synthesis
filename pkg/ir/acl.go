/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
	"reflect"

	"github.com/np-guard/models/pkg/ipblock"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type Action string

const (
	Allow Action = "allow"
	Deny  Action = "deny"
)

type ACLRule struct {
	Action      Action
	Direction   Direction
	Source      *ipblock.IPBlock
	Destination *ipblock.IPBlock
	Protocol    Protocol
	Explanation string
}

type ACL struct {
	Subnet   string
	Internal []ACLRule
	External []ACLRule
}

type ACLCollection struct {
	ACLs map[ID]map[string]*ACL
}

type ACLWriter interface {
	WriteACL(aclColl *ACLCollection, vpc string) error
}

func (r *ACLRule) isRedundant(rules []ACLRule) bool {
	for i := range rules {
		if rules[i].mustSupersede(r) {
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

func (r *ACLRule) Target() *ipblock.IPBlock {
	if r.Direction == Inbound {
		return r.Destination
	}
	return r.Source
}

func (a *ACL) Rules() []ACLRule {
	rules := a.Internal
	if len(a.External) != 0 {
		rules = append(rules, makeDenyInternal()...)
		rules = append(rules, a.External...)
	}
	return rules
}

func (a *ACL) AppendInternal(rule *ACLRule) {
	if a.External == nil {
		panic("ACLs should be created with non-null Internal")
	}
	if !rule.isRedundant(a.Internal) {
		a.Internal = append(a.Internal, *rule)
	}
}

func (a *ACL) Name() string {
	return fmt.Sprintf("acl-%v", a.Subnet)
}

func (a *ACL) AppendExternal(rule *ACLRule) {
	if a.External == nil {
		panic("ACLs should be created with non-null External")
	}
	if rule.isRedundant(a.External) {
		return
	}
	a.External = append(a.External, *rule)
}

func NewACLCollection() *ACLCollection {
	return &ACLCollection{ACLs: map[ID]map[string]*ACL{}}
}

func NewACL() *ACL {
	return &ACL{Internal: []ACLRule{}, External: []ACLRule{}}
}

func (c *ACLCollection) LookupOrCreate(name string) *ACL {
	vpcName := VpcFromScopedResource(name)
	if acl, ok := c.ACLs[vpcName][name]; ok {
		return acl
	}
	newACL := NewACL()
	newACL.Subnet = name
	if c.ACLs[vpcName] == nil {
		c.ACLs[vpcName] = make(map[string]*ACL)
	}
	c.ACLs[vpcName][name] = newACL
	return newACL
}

func (c *ACLCollection) Write(w Writer, vpc string) error {
	return w.WriteACL(c, vpc)
}

func (c *ACLCollection) SortedACLSubnets(vpc string) []string {
	if vpc == "" {
		return utils.SortedAllInnerMapsKeys(c.ACLs)
	}
	return utils.SortedInnerMapKeys(c.ACLs, vpc)
}
