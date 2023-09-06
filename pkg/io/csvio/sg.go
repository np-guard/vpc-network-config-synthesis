package csvio

import (
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func (w *Writer) WriteSG(collection *ir.SecurityGroupCollection) error {
	if err := w.w.Write(sgHeader()); err != nil {
		return err
	}
	for _, sgName := range collection.SortedSGNames() {
		if err := w.w.WriteAll(makeSGTable(collection.SGs[sgName], sgName)); err != nil {
			return err
		}
	}
	return nil
}

func sgHeader() []string {
	return []string{
		"SG",
		"Direction",
		"Protocol",
		"Target Type",
		"Target",
		"Value",
		"Description",
	}
}

func makeSGRow[T ir.RemoteType](rule *ir.SecurityGroupRule[T], sgName ir.SecurityGroupName) []string {
	return []string{
		string(sgName),
		direction(rule.Direction),
		printProtocolName(rule.Protocol),
		rule.Remote.String(),
		printValue(rule.Protocol, rule.Direction == ir.Inbound),
		rule.Explanation,
	}
}

func makeSGTable(t *ir.SecurityGroup, sgName ir.SecurityGroupName) [][]string {
	rules := t.Rules
	rows := make([][]string, len(rules))
	for i := range rules {
		rows[i] = makeSGRow(&rules[i], sgName)
	}
	return rows
}

func printValue(protocol ir.Protocol, isSource bool) string {
	switch p := protocol.(type) {
	case ir.ICMP:
		return printICMPTypeCode(protocol)
	case ir.TCPUDP:
		var portString string
		if isSource {
			portString = port(p.PortRangePair.SrcPort)
		} else {
			portString = port(p.PortRangePair.DstPort)
		}
		return portString
	case ir.AnyProtocol:
		return ""
	default:
		log.Panicf("Impossible protocol %v", p)
	}
	return ""
}
