// Package tfio implements output of ACLs in terraform format
package tfio

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/tf"
)

// Writer implements ir.Writer
type Writer struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

// Write prints an entire collection of acls as a sequence of terraform resources.
func (w *Writer) Write(c ir.Collection) error {
	output := collection(c).Print()
	_, err := w.w.WriteString(output)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
}

func portRangePair(t ir.PortRangePair, name string) tf.Block {
	var arguments []tf.Argument
	if t.DstPort.Min != ir.DefaultMinPort {
		arguments = append(arguments, tf.Argument{Name: "port_min", Value: strconv.Itoa(t.DstPort.Min)})
	}
	if t.DstPort.Max != ir.DefaultMaxPort {
		arguments = append(arguments, tf.Argument{Name: "port_max", Value: strconv.Itoa(t.DstPort.Max)})
	}
	if t.SrcPort.Min != ir.DefaultMinPort {
		arguments = append(arguments, tf.Argument{Name: "source_port_min", Value: strconv.Itoa(t.SrcPort.Min)})
	}
	if t.SrcPort.Max != ir.DefaultMaxPort {
		arguments = append(arguments, tf.Argument{Name: "source_port_max", Value: strconv.Itoa(t.SrcPort.Max)})
	}
	return tf.Block{
		Name:      name,
		Arguments: arguments,
	}
}

func protocol(t ir.Protocol) []tf.Block {
	switch p := t.(type) {
	case ir.TCPUDP:
		return []tf.Block{portRangePair(p.PortRangePair, strings.ToLower(p.Name()))}
	case ir.ICMP:
		var arguments []tf.Argument
		if p.ICMPCodeType != nil {
			arguments = append(arguments, tf.Argument{Name: "type", Value: strconv.Itoa(p.Type)})
			if p.Code != nil {
				arguments = append(arguments, tf.Argument{Name: "code", Value: strconv.Itoa(*p.Code)})
			}
		}
		return []tf.Block{{
			Name:      "icmp",
			Arguments: arguments,
		}}
	case ir.AnyProtocol:
		return []tf.Block{}
	}
	return nil
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

func action(a ir.Action) string {
	return string(a)
}

func direction(d ir.Direction) string {
	return string(d)
}

func rule(t *ir.Rule) tf.Block {
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

func singleACL(name string, t ir.ACL) tf.Block {
	blocks := make([]tf.Block, len(t.Rules))
	for i := range t.Rules {
		blocks[i] = rule(t.Rules[i])
	}
	return tf.Block{
		Name:   "resource",
		Labels: []string{quote("ibm_is_network_acl"), quote(name)},
		Arguments: []tf.Argument{
			{Name: "name", Value: quote(name + "-${var.initials}")}, //nolint:revive  // obvious false positive
			{Name: "resource_group", Value: "var.resource_group_id"},
			{Name: "vpc", Value: "var.vpc_id"},
		},
		Blocks: blocks,
	}
}

func collection(t ir.Collection) *tf.ConfigFile {
	var acls = make([]tf.Block, len(t.ACLs))
	i := 0
	for name := range t.ACLs {
		acls[i] = singleACL(name, t.ACLs[name])
		i += 1
	}
	return &tf.ConfigFile{
		Resources: acls,
	}
}
