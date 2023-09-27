package mdio

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func (w *Writer) WriteSG(collection *ir.SGCollection) error {
	if err := w.writeAll(sgHeader()); err != nil {
		return err
	}
	for _, sgName := range collection.SortedSGNames() {
		if err := w.writeAll(makeSGTable(collection.SGs[sgName], sgName)); err != nil {
			return err
		}
	}
	return nil
}

func sgHeader() [][]string {
	return [][]string{{
		"",
		"SG",
		"Direction",
		"Protocol",
		"Remote type",
		"Remote",
		"Value",
		"Description",
		"",
	}, {
		"",
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		leftAlign,
		"",
	}}
}

func makeSGRow(rule *ir.SGRule, sgName ir.SGName) []string {
	return []string{
		"",
		string(sgName),
		direction(rule.Direction),
		printProtocolName(rule.Protocol),
		sgRemote(rule.Remote),
		printValue(rule.Protocol, rule.Direction == ir.Inbound),
		rule.Explanation,
		"",
	}
}

func makeSGTable(t *ir.SG, sgName ir.SGName) [][]string {
	rules := t.Rules
	rows := make([][]string, len(rules))
	for i := range rules {
		rows[i] = makeSGRow(&rules[i], sgName)
	}
	return rows
}

func sGPort(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return "any port"
	default:
		return fmt.Sprintf("Ports %v-%v", p.Min, p.Max)
	}
}

func sgRemote(r ir.RemoteType) string {
	switch tr := r.(type) {
	case ir.IP:
		s := tr.String()
		if s == ir.AnyIP {
			return "Any IP"
		}
	case ir.CIDR:
		s := tr.String()
		if s == ir.AnyCIDR {
			return "Any CIDR"
		}
		return s
	case ir.SGName:
		return tr.String()
	default:
		log.Panicf("Impossible remote %v (%T)", r, r)
	}
	return ""
}

func printValue(protocol ir.Protocol, isSource bool) string {
	switch p := protocol.(type) {
	case ir.ICMP:
		return printICMPTypeCode(protocol)
	case ir.TCPUDP:
		var r ir.PortRange
		if isSource {
			r = p.PortRangePair.SrcPort
		} else {
			r = p.PortRangePair.DstPort
		}
		return sGPort(r)
	case ir.AnyProtocol:
		return ""
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}
