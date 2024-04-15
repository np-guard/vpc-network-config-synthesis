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

type Action string

const (
	Allow Action = "allow"
	Deny  Action = "deny"
)

type ACLRule struct {
	Action      Action
	Direction   Direction
	Source      IP
	Destination IP
	Protocol    Protocol
	Explanation string
}

type ACL struct {
	Subnet   string
	Internal []ACLRule
	External []ACLRule
}

type ACLCollection struct {
	ACLs map[string]*ACL
}

type ACLWriter interface {
	WriteACL(*ACLCollection) error
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

func (r *ACLRule) Target() IP {
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
	if rule.isRedundant(a.Internal) {
		return
	}
	a.Internal = append(a.Internal, *rule)
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
	return &ACLCollection{ACLs: map[string]*ACL{}}
}

func MergeACLCollections(collections ...*ACLCollection) *ACLCollection {
	result := NewACLCollection()
	for _, c := range collections {
		for a := range c.ACLs {
			acl := c.LookupOrCreate(a)
			for r := range acl.Internal {
				result.LookupOrCreate(a).AppendInternal(&acl.Internal[r])
			}
			for r := range acl.External {
				result.LookupOrCreate(a).AppendExternal(&acl.External[r])
			}
		}
	}
	return result
}

func NewACL() *ACL {
	return &ACL{Internal: []ACLRule{}, External: []ACLRule{}}
}

func (c *ACLCollection) LookupOrCreate(name string) *ACL {
	acl, ok := c.ACLs[name]
	if ok {
		return acl
	}
	newACL := NewACL()
	newACL.Subnet = name
	c.ACLs[name] = newACL
	return newACL
}

func (c *ACLCollection) Write(w Writer) error {
	return w.WriteACL(c)
}

func (c *ACLCollection) SortedACLSubnets() []string {
	return utils.SortedKeys(c.ACLs)
}
