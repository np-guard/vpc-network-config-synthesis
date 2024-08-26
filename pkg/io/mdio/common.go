/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package mdio implements output of ACLs and security groups in CSV format
package mdio

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/np-guard/models/pkg/netp"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const leftAlign = " :--- "

// Writer implements ir.SynthWriter
type Writer struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

func (w *Writer) writeAll(rows [][]string) error {
	for _, row := range rows {
		_, err := w.w.WriteString(strings.Join(row, " | ") + "\n")
		if err != nil {
			return err
		}
	}
	w.w.Flush()
	return nil
}

const (
	anyProtocol  = "ALL"
	nonIcmp      = "-"
	anyIcmpValue = "Any"
)

func direction(d ir.Direction) string {
	if d == ir.Inbound {
		return "Inbound"
	}
	return "Outbound"
}

func printICMPTypeCode(protocol netp.Protocol) string {
	p, ok := protocol.(netp.ICMP)
	if !ok {
		return nonIcmp
	}
	icmpType := anyIcmpValue
	icmpCode := anyIcmpValue
	if typeCode := p.ICMPTypeCode(); typeCode != nil {
		icmpType = strconv.Itoa(typeCode.Type)
		if typeCode.Code != nil {
			icmpCode = strconv.Itoa(*typeCode.Code)
		}
	}
	return fmt.Sprintf("Type: %v, Code: %v", icmpType, icmpCode)
}

func printProtocolName(protocol netp.Protocol) string {
	switch p := protocol.(type) {
	case netp.ICMP:
		return "ICMP"
	case netp.TCPUDP:
		return strings.ToUpper(string(p.ProtocolString()))
	case netp.AnyProtocol:
		return anyProtocol
	}
	return ""
}
