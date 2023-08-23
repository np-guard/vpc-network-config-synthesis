// Package csvio implements output of ACLs in CSV format
package csvio

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Writer implements ir.Writer
type Writer struct {
	w *csv.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: csv.NewWriter(w)}
}

// Write prints an entire collection of acls as a single CSV table.
// This is mostly useful when there is only a single ir.ACL item in the collection
func (w *Writer) Write(collection *ir.Collection) error {
	if err := w.w.Write(header()); err != nil {
		return err
	}
	for _, item := range collection.ACLs {
		if err := w.w.WriteAll(makeTable(item)); err != nil {
			return err
		}
	}
	return nil
}

const allPorts = "All"

func makeTable(t ir.ACL) [][]string {
	rows := make([][]string, len(t.Rules))
	for i, rule := range t.Rules {
		rows[i] = makeRow(i+1, rule)
	}
	return rows
}

func port(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return allPorts
	case p.Min == p.Max:
		return fmt.Sprintf("%v", p.Max)
	default:
		return fmt.Sprintf("%v-%v", p.Min, p.Max)
	}
}

func action(a ir.Action) string {
	return string(a)
}

func direction(d ir.Direction) string {
	return string(d)
}

func printPortRange(protocol ir.Protocol, isSrcPort bool) string {
	switch p := protocol.(type) {
	case ir.ICMP:
		return "-"
	case ir.TCPUDP:
		if isSrcPort {
			return port(p.PortRangePair.SrcPort)
		}
		return port(p.PortRangePair.DstPort)
	case ir.AnyProtocol:
		return allPorts
	default:
		log.Fatalf("Impossible protocol %v", p)
	}
	return ""
}

func header() []string {
	return []string{
		"#",
		"direction",
		"action",
		"source",
		"source port",
		"destination",
		"destination port",
		"protocol",
		"description",
	}
}

func makeRow(i int, rule *ir.Rule) []string {
	return []string{
		strconv.Itoa(i),
		direction(rule.Direction),
		action(rule.Action),
		rule.Source,
		printPortRange(rule.Protocol, true),
		rule.Destination,
		printPortRange(rule.Protocol, false),
		rule.Protocol.Name(),
		rule.Name,
	}
}
