/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package io

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

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

func printPorts(p interval.Interval) string {
	if p.Equal(netp.AllPorts()) {
		return "any port"
	}
	return fmt.Sprintf("ports %v-%v", p.Start(), p.End())
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
