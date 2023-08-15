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
	SwapSrcDstPortRange() Protocol
}

func (t TCP) SwapSrcDstPortRange() Protocol { return TCP{Swap(t.PortRangePair)} }

func (t UDP) SwapSrcDstPortRange() Protocol { return UDP{Swap(t.PortRangePair)} }

func (t ICMP) SwapSrcDstPortRange() Protocol { return ICMP{Code: t.Code, Type: t.Type} }

func (t AnyProtocol) SwapSrcDstPortRange() Protocol { return AnyProtocol{} }

type Rule struct {
	Name        string
	Deny        bool
	Source      string
	Destination string
	Outbound    bool
	Protocol    Protocol
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

func action(deny bool) string {
	if deny {
		return "deny"
	}
	return "allow"
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
	switch t.Protocol.(type) {
	case AnyProtocol:
		break
	default:
		blocks = []tf.Block{
			t.Protocol.(tf.Blockable).Terraform(),
		}
	}
	arguments := []tf.Argument{
		{Name: "name", Value: quote(t.Name)},
		{Name: "action", Value: quote(action(t.Deny))},
		{Name: "direction", Value: quote(outbound(t.Outbound))},
		{Name: "source", Value: quote(t.Source)},
		{Name: "destination", Value: quote(t.Destination)},
	}
	return tf.Block{Name: "rules",
		Arguments: arguments,
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
