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

type Connection interface {
	tf.Blockable
}

const defaultMinTransportPort = 1
const defaultMaxTransportPort = 65535

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
	var arguments []tf.Argument
	if t.MinPort != defaultMinTransportPort {
		arguments = append(arguments, tf.Argument{Name: "port_min", Value: strconv.Itoa(t.MinPort)})
	}
	if t.MaxPort != defaultMaxTransportPort {
		arguments = append(arguments, tf.Argument{Name: "port_max", Value: strconv.Itoa(t.MaxPort)})
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
	var arguments []tf.Argument
	if t.Code != nil {
		arguments = append(arguments, tf.Argument{Name: "code", Value: strconv.Itoa(*t.Code)})
	}
	if t.Type != nil {
		arguments = append(arguments, tf.Argument{Name: "type", Value: strconv.Itoa(*t.Type)})
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
		Arguments: []tf.Argument{
			{Name: "name", Value: quote(t.Name)},
			{Name: "action", Value: quote(allow(t.Allow))},
			{Name: "direction", Value: quote(outbound(t.Outbound))},
			{Name: "source", Value: quote(t.Source)},
			{Name: "destination", Value: quote(t.Destination)},
		},
		Blocks: blocks,
	}
}

func (t *ACL) Terraform() tf.Block {
	return tf.Block{
		Name:   "resource",
		Labels: []string{quote("ibm_is_network_acl"), quote(t.Name)},
		Arguments: []tf.Argument{
			{Name: "name", Value: quote(t.Name + "-${var.initials}")}, //nolint:revive  // obvious false positive
			{Name: "resource_group", Value: t.ResourceGroup},
			{Name: "vpc", Value: t.Vpc},
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
