// Package tfio implements output of ACLs and security groups in terraform format
package tfio

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Writer implements ir.Writer
type Writer struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
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
		return []tf.Block{portRangePair(p.PortRangePair, strings.ToLower(string(p.Protocol)))}
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

func verifyName(name string) {
	pattern := "^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$"
	_, err := regexp.MatchString(pattern, name)
	if err != nil {
		log.Fatalf("\"name\" should match regexp %q", pattern)
	}
}
