// Package aclcsv implements output of ACLs in CSV format
package aclcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"

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
		rows[i] = makeRow(rule)
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

func printPortRange(protocol acl.Protocol, d acl.Direction) string {
	switch p := protocol.(type) {
	case acl.ICMP:
		return "-"
	case acl.UDP:
		if d == acl.Outbound {
			return port(p.SrcPort)
		}
		return port(p.DstPort)
	case acl.TCP:
		if d == acl.Outbound {
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

func makeRow(rule *acl.Rule) []string {
	return []string{
		direction(rule.Direction),
		action(rule.Action),
		rule.Source,
		rule.Destination,
		rule.Protocol.Name(),
		printPortRange(rule.Protocol, rule.Direction),
	}
}
