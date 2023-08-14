// Package acl describes Network ACLs
package acl

import (
	"fmt"
	"strconv"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/tf"
)

type PortRange struct {
	Min int
	Max int
}

type PortRangePair struct {
	SrcPort PortRange
	DstPort PortRange
}

const DefaultMinPort = 1
const DefaultMaxPort = 65535

func Swap(pair PortRangePair) PortRangePair {
	return PortRangePair{SrcPort: pair.DstPort, DstPort: pair.SrcPort}
}

type TCP struct {
	PortRangePair
}

type UDP struct {
	PortRangePair
}

type ICMP struct {
	Code *int
	Type *int
}

type AnyProtocol struct{}

type Protocol interface {
	Rule(f Flow) *Rule
	Swap() Protocol
}

func (t TCP) Swap() Protocol { return TCP{Swap(t.PortRangePair)} }

func (t UDP) Swap() Protocol { return UDP{Swap(t.PortRangePair)} }

func (t ICMP) Swap() Protocol { return ICMP{Code: t.Code, Type: t.Type} }

func (t AnyProtocol) Swap() Protocol { return AnyProtocol{} }

func (t TCP) Rule(f Flow) *Rule { return &Rule{Flow: f, Protocol: t} }

func (t UDP) Rule(f Flow) *Rule { return &Rule{Flow: f, Protocol: t} }

func (t ICMP) Rule(f Flow) *Rule { return &Rule{Flow: f, Protocol: t} }

func (t AnyProtocol) Rule(f Flow) *Rule { return &Rule{Flow: f, Protocol: nil} }

type Flow struct {
	Name        string
	Allow       bool
	Source      string
	Destination string
	Outbound    bool
}

type Rule struct {
	Flow     Flow
	Protocol tf.Blockable
}

func (f *Flow) Terraform() []tf.Argument {
	return []tf.Argument{
		{Name: "name", Value: quote(f.Name)},
		{Name: "action", Value: quote(allow(f.Allow))},
		{Name: "direction", Value: quote(outbound(f.Outbound))},
		{Name: "source", Value: quote(f.Source)},
		{Name: "destination", Value: quote(f.Destination)},
	}
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

func (t *PortRangePair) Terraform(name string) tf.Block {
	var arguments []tf.Argument
	if t.DstPort.Min != DefaultMinPort {
		arguments = append(arguments, tf.Argument{Name: "port_min", Value: strconv.Itoa(t.DstPort.Min)})
	}
	if t.DstPort.Max != DefaultMaxPort {
		arguments = append(arguments, tf.Argument{Name: "port_max", Value: strconv.Itoa(t.DstPort.Max)})
	}
	if t.SrcPort.Min != DefaultMinPort {
		arguments = append(arguments, tf.Argument{Name: "source_port_min", Value: strconv.Itoa(t.SrcPort.Min)})
	}
	if t.SrcPort.Max != DefaultMaxPort {
		arguments = append(arguments, tf.Argument{Name: "source_port_max", Value: strconv.Itoa(t.SrcPort.Max)})
	}
	return tf.Block{
		Name:      name,
		Arguments: arguments,
	}
}

func (t TCP) Terraform() tf.Block {
	return t.PortRangePair.Terraform("tcp")
}

func (t UDP) Terraform() tf.Block {
	return t.PortRangePair.Terraform("udp")
}

func (t ICMP) Terraform() tf.Block {
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
		Arguments: t.Flow.Terraform(),
		Blocks:    blocks,
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
