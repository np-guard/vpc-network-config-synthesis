package csvio

import (
	"fmt"
	"log"
	"strconv"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func makeACLTable(t *ir.ACL, subnet string) [][]string {
	rules := t.Rules()
	rows := make([][]string, len(rules))
	for i := range rules {
		rows[i] = makeACLRow(i+1, &rules[i], t.Name(), subnet)
	}
	return rows
}

// Write prints an entire collection of acls as a single CSV table.
func (w *Writer) WriteACL(collection *ir.ACLCollection) error {
	if err := w.w.Write(aclHeader()); err != nil {
		return err
	}
	for _, subnet := range collection.SortedACLSubnets() {
		if err := w.w.WriteAll(makeACLTable(collection.ACLs[subnet], subnet)); err != nil {
			return err
		}
	}
	return nil
}

func action(a ir.Action) string {
	if a == ir.Deny {
		return "Deny"
	}
	return "Allow"
}

func aclHeader() []string {
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

func makeACLRow(priority int, rule *ir.ACLRule, aclName, subnet string) []string {
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
