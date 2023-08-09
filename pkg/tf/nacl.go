// Package tf deals with Terraform-specific abstract syntax
// It has two parts.
// * This file represent ACL-specific data.
// * syntax.go represent terraform syntax.
package tf

import (
	"fmt"
	"strconv"
)

type PortRange struct {
	MinPort int
	MaxPort int
}

type TCP struct {
	PortRange
}

type UDP struct {
	PortRange
}

type ICMP struct {
	Code *int
	Type *int
}

type ACLRule struct {
	Name        string
	Allow       bool
	Source      string
	Destination string
	Outbound    bool
	Protocol    blockable
}

type ACL struct {
	Name          string
	ResourceGroup string
	Vpc           string
	Rules         []*ACLRule
}

type ACLCollection struct {
	Items []*ACL
}

type ACLRuleMaker interface {
	newACLRule(name string, allow bool, source string, destination string, outbound bool) *ACLRule
}

func (t *TCP) newACLRule(name string, allow bool, source, destination string, outbound bool) *ACLRule {
	return &ACLRule{Name: name, Allow: allow, Source: source, Destination: destination, Outbound: outbound, Protocol: t}
}

func (t *UDP) newACLRule(name string, allow bool, source, destination string, outbound bool) *ACLRule {
	return &ACLRule{Name: name, Allow: allow, Source: source, Destination: destination, Outbound: outbound, Protocol: t}
}

func (t *ICMP) newACLRule(name string, allow bool, source, destination string, outbound bool) *ACLRule {
	return &ACLRule{Name: name, Allow: allow, Source: source, Destination: destination, Outbound: outbound, Protocol: t}
}

func NewACLRule(t ACLRuleMaker, name string, allow bool, source, destination string, outbound bool) *ACLRule {
	if t == nil {
		return &ACLRule{Name: name, Allow: allow, Source: source, Destination: destination, Outbound: outbound}
	}
	return t.newACLRule(name, allow, source, destination, outbound)
}

const maxTransportPort = 65535

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

func allow(b bool) string {
	if b {
		return "allow"
	}
	return "deny"
}

func outbound(b bool) string {
	if b {
		return "outbound"
	}
	return "inbound"
}

func (t *PortRange) terraform(name string) block {
	arguments := map[string]string{}
	if t.MinPort != 0 {
		arguments["port_min"] = strconv.Itoa(t.MinPort)
	}
	if t.MaxPort != maxTransportPort {
		arguments["port_max"] = strconv.Itoa(t.MaxPort)
	}
	return block{
		Name:      name,
		Arguments: arguments,
	}
}

func (t *TCP) terraform() block {
	return t.PortRange.terraform("tcp")
}

func (t *UDP) terraform() block {
	return t.PortRange.terraform("udp")
}

func (t *ICMP) terraform() block {
	arguments := map[string]string{}
	if t.Type != nil {
		arguments["type"] = strconv.Itoa(*t.Type)
	}
	if t.Code != nil {
		arguments["code"] = strconv.Itoa(*t.Code)
	}
	return block{
		Name:      "icmp",
		Arguments: arguments,
	}
}

func (t *ACLRule) terraform() block {
	var blocks []block
	if t.Protocol != nil {
		blocks = []block{
			t.Protocol.terraform(),
		}
	}
	return block{Name: "rules",
		Arguments: map[string]string{
			"name":        quote(t.Name),
			"action":      quote(allow(t.Allow)),
			"source":      quote(t.Source),
			"destination": quote(t.Destination),
			"direction":   quote(outbound(t.Outbound)),
		},
		Blocks: blocks,
	}
}

func (t *ACL) terraform() block {
	return block{
		Name:   "resource",
		Labels: []string{quote("ibm_is_network_acl"), quote(t.Name)},
		Arguments: map[string]string{
			"name":           quote(t.Name + "-${var.initials}"), //nolint:revive  // obvious false positive
			"resource_group": t.ResourceGroup,
			"vpc":            t.Vpc,
		},
		Blocks: blocks(t.Rules),
	}
}

func (t *ACLCollection) terraform() configFile {
	return configFile{
		Resources: blocks(t.Items),
	}
}

func (t *ACLCollection) Print() string {
	x := t.terraform()
	return x.print()
}
