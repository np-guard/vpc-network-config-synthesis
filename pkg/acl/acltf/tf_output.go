// Package acltf implements output of ACLs in terraform format
package acltf

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/tf"
)

// Writer implements acl.Writer
type Writer struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

// Write prints an entire collection of acls as a sequence of terraform resources.
func (w *Writer) Write(c acl.Collection) error {
	_, err := w.w.WriteString(collection(c).Print())
	return err
}

func portRangePair(t acl.PortRangePair, name string) tf.Block {
	var arguments []tf.Argument
	if t.DstPort.Min != acl.DefaultMinPort {
		arguments = append(arguments, tf.Argument{Name: "port_min", Value: strconv.Itoa(t.DstPort.Min)})
	}
	if t.DstPort.Max != acl.DefaultMaxPort {
		arguments = append(arguments, tf.Argument{Name: "port_max", Value: strconv.Itoa(t.DstPort.Max)})
	}
	if t.SrcPort.Min != acl.DefaultMinPort {
		arguments = append(arguments, tf.Argument{Name: "source_port_min", Value: strconv.Itoa(t.SrcPort.Min)})
	}
	if t.SrcPort.Max != acl.DefaultMaxPort {
		arguments = append(arguments, tf.Argument{Name: "source_port_max", Value: strconv.Itoa(t.SrcPort.Max)})
	}
	return tf.Block{
		Name:      name,
		Arguments: arguments,
	}
}

func protocol(t acl.Protocol) []tf.Block {
	switch p := t.(type) {
	case acl.TCP:
		return []tf.Block{portRangePair(p.PortRangePair, "tcp")}
	case acl.UDP:
		return []tf.Block{portRangePair(p.PortRangePair, "udp")}
	case acl.ICMP:
		var arguments []tf.Argument
		if p.Code != nil {
			arguments = append(arguments, tf.Argument{Name: "code", Value: strconv.Itoa(*p.Code)})
		}
		if p.Type != nil {
			arguments = append(arguments, tf.Argument{Name: "type", Value: strconv.Itoa(*p.Type)})
		}
		return []tf.Block{{
			Name:      "icmp",
			Arguments: arguments,
		}}
	case acl.AnyProtocol:
		return []tf.Block{}
	}
	return nil
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

func action(a acl.Action) string {
	switch a {
	case acl.Allow:
		return "allow"
	case acl.Deny:
		return "deny"
	}
	log.Fatalf("Impossible action %q", a)
	return ""
}

func direction(d acl.Direction) string {
	switch d {
	case acl.Outbound:
		return "outbound"
	case acl.Inbound:
		return "inbound"
	}
	log.Fatalf("Impossible direction %q", d)
	return ""
}

func rule(t *acl.Rule) tf.Block {
	arguments := []tf.Argument{
		{Name: "name", Value: quote(t.Name)},
		{Name: "action", Value: quote(action(t.Action))},
		{Name: "direction", Value: quote(direction(t.Direction))},
		{Name: "source", Value: quote(t.Source)},
		{Name: "destination", Value: quote(t.Destination)},
	}
	return tf.Block{Name: "rules",
		Arguments: arguments,
		Blocks:    protocol(t.Protocol),
	}
}

func singleACL(t *acl.ACL) tf.Block {
	blocks := make([]tf.Block, len(t.Rules))
	for i := range t.Rules {
		blocks[i] = rule(t.Rules[i])
	}
	return tf.Block{
		Name:   "resource",
		Labels: []string{quote("ibm_is_network_acl"), quote(t.Name)},
		Arguments: []tf.Argument{
			{Name: "name", Value: quote(t.Name + "-${var.initials}")}, //nolint:revive  // obvious false positive
			{Name: "resource_group", Value: t.ResourceGroup},
			{Name: "vpc", Value: t.Vpc},
		},
		Blocks: blocks,
	}
}

func collection(t acl.Collection) *tf.ConfigFile {
	acls := make([]tf.Block, len(t.Items))
	for i := range t.Items {
		acls[i] = singleACL(t.Items[i])
	}
	return &tf.ConfigFile{
		Resources: acls,
	}
}
