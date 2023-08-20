// Package aclcsv implements output of ACLs in CSV format
package aclcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl"
)

// Writer implements acl.Writer
type Writer struct {
	w *csv.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: csv.NewWriter(w)}
}

// Write prints an entire collection of acls as a single CSV table.
// This is mostly useful when there is only a single acl.ACL item in the collection
func (w *Writer) Write(collection acl.Collection) error {
	if err := w.w.Write(header()); err != nil {
		return err
	}
	for _, item := range collection.Items {
		if err := w.w.WriteAll(makeTable(item)); err != nil {
			return err
		}
	}
	return nil
}

const allPorts = "All"

func makeTable(t *acl.ACL) [][]string {
	rows := make([][]string, len(t.Rules))
	for i, rule := range t.Rules {
		rows[i] = makeRow(i+1, rule)
	}
	return rows
}

func port(p acl.PortRange) string {
	switch {
	case p.Min == acl.DefaultMinPort && p.Max == acl.DefaultMaxPort:
		return allPorts
	case p.Min == p.Max:
		return fmt.Sprintf("%v", p.Max)
	default:
		return fmt.Sprintf("%v-%v", p.Min, p.Max)
	}
}

func action(a acl.Action) string {
	return string(a)
}

func direction(d acl.Direction) string {
	return string(d)
}

func printPortRange(protocol acl.Protocol, isSrcPort bool) string {
	switch p := protocol.(type) {
	case acl.ICMP:
		return "-"
	case acl.UDP:
		if isSrcPort {
			return port(p.SrcPort)
		}
		return port(p.DstPort)
	case acl.TCP:
		if isSrcPort {
			return port(p.SrcPort)
		}
		return port(p.DstPort)
	case acl.AnyProtocol:
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
	}
}

func makeRow(i int, rule *acl.Rule) []string {
	return []string{
		strconv.Itoa(i),
		direction(rule.Direction),
		action(rule.Action),
		rule.Source,
		printPortRange(rule.Protocol, true),
		rule.Destination,
		printPortRange(rule.Protocol, false),
		rule.Protocol.Name(),
	}
}
