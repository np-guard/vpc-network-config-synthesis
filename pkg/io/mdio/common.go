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

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const leftAlign = " :--- "

// Writer implements ir.Writer
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
