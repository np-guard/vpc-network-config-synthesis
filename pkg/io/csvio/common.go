package csvio

import (
	"encoding/csv"
	"fmt"
	"io"
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

const (
	anyProtocol  = "ALL"
	nonIcmp      = "-" // IBM cloud uses "—"
	anyIcmpValue = "Any"
)

func direction(d ir.Direction) string {
	if d == ir.Inbound {
		return "Inbound"
	}
	return "Outbound"
}

func port(p ir.PortRange) string {
	switch {
	case p.Min == ir.DefaultMinPort && p.Max == ir.DefaultMaxPort:
		return "any port"
	default:
		return fmt.Sprintf("ports %v-%v", p.Min, p.Max)
	}
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
