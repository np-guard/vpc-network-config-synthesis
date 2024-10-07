/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import (
	"fmt"
	"reflect"

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
		Subnet   string
		Internal []*ACLRule
		External []*ACLRule
	}

	ACLCollection struct {
		ACLs map[ID]map[string]*ACL
	}

	ACLWriter interface {
		WriteACL(aclColl *ACLCollection, vpc string) error
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
		a.Internal = append(a.Internal, rule)
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
	a.External = append(a.External, rule)
}

func NewACLCollection() *ACLCollection {
	return &ACLCollection{ACLs: map[ID]map[string]*ACL{}}
}

func NewACL(subnet string) *ACL {
	return &ACL{Subnet: subnet, Internal: []*ACLRule{}, External: []*ACLRule{}}
}

func (c *ACLCollection) LookupOrCreate(subnet string) *ACL {
	vpcName := VpcFromScopedResource(subnet)
	if acl, ok := c.ACLs[vpcName][subnet]; ok {
		return acl
	}
	newACL := NewACL(subnet)
	if c.ACLs[vpcName] == nil {
		c.ACLs[vpcName] = make(map[string]*ACL)
	}
	c.ACLs[vpcName][subnet] = newACL
	return newACL
}

func (c *ACLCollection) VpcNames() []string {
	return utils.SortedMapKeys(c.ACLs)
}

func (c *ACLCollection) Write(w Writer, vpc string) error {
	return w.WriteACL(c, vpc)
}

func (c *ACLCollection) SortedACLSubnets(vpc string) []string {
	if vpc == "" {
		return utils.SortedAllInnerMapsKeys(c.ACLs)
	}
	return utils.SortedMapKeys(c.ACLs[vpc])
}
