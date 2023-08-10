// Package acl describes Network ACLs
package acl

import (
	"fmt"
	"strconv"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/tf"
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

type Rule struct {
	Name        string
	Allow       bool
	Source      string
	Destination string
	Outbound    bool
	Protocol    tf.Blockable
}

type ACL struct {
	Name          string
	ResourceGroup string
	Vpc           string
	Rules         []*Rule
}

type Collection struct {
	Items []*ACL
}

type RuleMaker interface {
	newRule(name string, allow bool, source string, destination string, outbound bool) *Rule
}

func (t *TCP) newRule(name string, allow bool, source, destination string, outbound bool) *Rule {
	return &Rule{Name: name, Allow: allow, Source: source, Destination: destination, Outbound: outbound, Protocol: t}
}

func (t *UDP) newRule(name string, allow bool, source, destination string, outbound bool) *Rule {
	return &Rule{Name: name, Allow: allow, Source: source, Destination: destination, Outbound: outbound, Protocol: t}
}

func (t *ICMP) newRule(name string, allow bool, source, destination string, outbound bool) *Rule {
	return &Rule{Name: name, Allow: allow, Source: source, Destination: destination, Outbound: outbound, Protocol: t}
}

func NewRule(t RuleMaker, name string, allow bool, source, destination string, outbound bool) *Rule {
	if t == nil {
		return &Rule{Name: name, Allow: allow, Source: source, Destination: destination, Outbound: outbound}
	}
	return t.newRule(name, allow, source, destination, outbound)
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

func (t *PortRange) Terraform(name string) tf.Block {
	arguments := map[string]string{}
	if t.MinPort != 0 {
		arguments["port_min"] = strconv.Itoa(t.MinPort)
	}
	if t.MaxPort != maxTransportPort {
		arguments["port_max"] = strconv.Itoa(t.MaxPort)
	}
	return tf.Block{
		Name:      name,
		Arguments: arguments,
	}
}

func (t *TCP) Terraform() tf.Block {
	return t.PortRange.Terraform("tcp")
}

func (t *UDP) Terraform() tf.Block {
	return t.PortRange.Terraform("udp")
}

func (t *ICMP) Terraform() tf.Block {
	arguments := map[string]string{}
	if t.Type != nil {
		arguments["type"] = strconv.Itoa(*t.Type)
	}
	if t.Code != nil {
		arguments["code"] = strconv.Itoa(*t.Code)
	}
	return tf.Block{
		Name:      "icmp",
		Arguments: arguments,
	}
}

func (t *Rule) Terraform() tf.Block {
	var blocks []tf.Block
	if t.Protocol != nil {
		blocks = []tf.Block{
			t.Protocol.Terraform(),
		}
	}
	return tf.Block{Name: "rules",
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

func (t *ACL) Terraform() tf.Block {
	return tf.Block{
		Name:   "resource",
		Labels: []string{quote("ibm_is_network_acl"), quote(t.Name)},
		Arguments: map[string]string{
			"name":           quote(t.Name + "-${var.initials}"), //nolint:revive  // obvious false positive
			"resource_group": t.ResourceGroup,
			"vpc":            t.Vpc,
		},
		Blocks: tf.Blocks(t.Rules),
	}
}

func (t *Collection) Terraform() tf.ConfigFile {
	return tf.ConfigFile{
		Resources: tf.Blocks(t.Items),
	}
}

func (t *Collection) Print() string {
	x := t.Terraform()
	return x.Print()
}
