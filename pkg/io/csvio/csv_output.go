// Package csvio implements output of ACLs in CSV format
package csvio

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

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
	for _, subnet := range collection.SortedACLSubnets() {
		if err := w.w.WriteAll(makeTable(collection.ACLs[subnet], subnet)); err != nil {
			return err
		}
	}
	return nil
}

const (
	all = "All"
	na  = "-"
)

func makeTable(t *ir.ACL, subnet string) [][]string {
	rules := t.Rules()
	rows := make([][]string, len(rules))
	for i := range rules {
		rows[i] = makeRow(i+1, &rules[i], t.Name(), subnet)
	}
	return rows
}

func port(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return all
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

func header() []string {
	return []string{
		"acl",
		"subnet",
		"#",
		"direction",
		"action",
		"source",
		"source port",
		"destination",
		"destination port",
		"protocol",
		"type",
		"code",
		"description",
	}
}

func makeRow(i int, rule *ir.Rule, aclName, subnet string) []string {
	icmpType, icmpCode := printICMPTypeCode(rule.Protocol)
	return []string{
		aclName,
		subnet,
		strconv.Itoa(i),
		direction(rule.Direction),
		action(rule.Action),
		rule.Source,
		printPortRange(rule.Protocol, true),
		rule.Destination,
		printPortRange(rule.Protocol, false),
		printProtocolName(rule.Protocol),
		icmpType,
		icmpCode,
		rule.Explanation,
	}
}

func printICMPTypeCode(protocol ir.Protocol) (icmpType, icmpCode string) {
	p, ok := protocol.(ir.ICMP)
	if !ok {
		return na, na
	}
	icmpType = all
	icmpCode = all
	if p.ICMPCodeType != nil {
		icmpType = strconv.Itoa(p.Type)
		if p.Code != nil {
			icmpCode = strconv.Itoa(*p.Code)
		}
	}
	return
}

func printPortRange(protocol ir.Protocol, isSrcPort bool) string {
	switch p := protocol.(type) {
	case ir.ICMP:
		return na
	case ir.TCPUDP:
		if isSrcPort {
			return port(p.PortRangePair.SrcPort)
		}
		return port(p.PortRangePair.DstPort)
	case ir.AnyProtocol:
		return all
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}

func printProtocolName(protocol ir.Protocol) string {
	switch p := protocol.(type) {
	case ir.ICMP:
		return "ICMP"
	case ir.TCPUDP:
		return strings.ToUpper(string(p.Protocol))
	case ir.AnyProtocol:
		return all
	}
	return ""
}
