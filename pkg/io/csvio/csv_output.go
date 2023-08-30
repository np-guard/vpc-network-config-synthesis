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
	anyProtocol  = "ALL"
	nonIcmp      = "-" // IBM cloud uses "—"
	anyIcmpValue = "Any"
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
		return "any port"
	default:
		return fmt.Sprintf("ports %v-%v", p.Min, p.Max)
	}
}

func action(a ir.Action) string {
	if a == ir.Deny {
		return "Deny"
	}
	return "Allow"
}

func direction(d ir.Direction) string {
	if d == ir.Inbound {
		return "Inbound"
	}
	return "Outbound"
}

func header() []string {
	return []string{
		"Acl",
		"Subnet",
		"Direction",
		"Rule priority",
		"Allow or deny",
		"Protocol",
		"Source",
		"Destination",
		"Value",
		"Description",
	}
}

func makeRow(priority int, rule *ir.Rule, aclName, subnet string) []string {
	return []string{
		aclName,
		subnet,
		direction(rule.Direction),
		strconv.Itoa(priority),
		action(rule.Action),
		printProtocolName(rule.Protocol),
		printIP(rule.Source, rule.Protocol, true),
		printIP(rule.Destination, rule.Protocol, false),
		printICMPTypeCode(rule.Protocol),
		rule.Explanation,
	}
}

func printIP(ip ir.IP, protocol ir.Protocol, isSource bool) string {
	ipString := ip.String()
	if ipString == "0.0.0.0/0" {
		ipString = "Any IP"
	}
	switch p := protocol.(type) {
	case ir.ICMP:
		return ipString
	case ir.TCPUDP:
		var portString string
		if isSource {
			portString = port(p.PortRangePair.SrcPort)
		} else {
			portString = port(p.PortRangePair.DstPort)
		}
		return fmt.Sprintf("%v, %v", ipString, portString)
	case ir.AnyProtocol:
		return ipString
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}

func printICMPTypeCode(protocol ir.Protocol) string {
	p, ok := protocol.(ir.ICMP)
	if !ok {
		return nonIcmp
	}
	icmpType := anyIcmpValue
	icmpCode := anyIcmpValue
	if p.ICMPCodeType != nil {
		icmpType = strconv.Itoa(p.Type)
		if p.Code != nil {
			icmpCode = strconv.Itoa(*p.Code)
		}
	}
	return fmt.Sprintf("Type: %v, Code: %v", icmpType, icmpCode)
}

func printProtocolName(protocol ir.Protocol) string {
	switch p := protocol.(type) {
	case ir.ICMP:
		return "ICMP"
	case ir.TCPUDP:
		return strings.ToUpper(string(p.Protocol))
	case ir.AnyProtocol:
		return anyProtocol
	}
	return ""
}
