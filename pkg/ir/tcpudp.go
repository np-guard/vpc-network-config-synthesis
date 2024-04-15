/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

import "log"

type TransportLayerProtocolName string

const (
	TCP TransportLayerProtocolName = "TCP"
	UDP TransportLayerProtocolName = "UDP"
)

const DefaultMinPort = 1
const DefaultMaxPort = 65535

type PortRange struct {
	// Minimal port; default is DefaultMinPort
	Min int

	// Maximal port; default is DefaultMaxPort
	Max int
}

type PortRangePair struct {
	SrcPort PortRange
	DstPort PortRange
}

type TCPUDP struct {
	Protocol      TransportLayerProtocolName
	PortRangePair PortRangePair
}

func (t TCPUDP) InverseDirection() Protocol {
	switch t.Protocol {
	case TCP:
		return TCPUDP{
			Protocol:      TCP,
			PortRangePair: PortRangePair{SrcPort: t.PortRangePair.DstPort, DstPort: t.PortRangePair.SrcPort},
		}
	case UDP:
		return nil
	default:
		log.Panicf("Impossible protocol: %v", t.Protocol)
	}
	return nil
}
