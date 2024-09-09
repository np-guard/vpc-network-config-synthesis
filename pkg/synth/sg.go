/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package synth

import (
	"log"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const SGTypeNotSupported = "SG: src/dst of type %s is not supported."

type SGSynthesizer struct {
	spec   *ir.Spec
	result *ir.SGCollection
}

// NewSGSynthesizer creates and returns a new SGSynthesizer instance
func NewSGSynthesizer(s *ir.Spec, _ bool) Synthesizer {
	return &SGSynthesizer{spec: s, result: ir.NewSGCollection()}
}

func (s *SGSynthesizer) Synth() ir.Collection {
	return s.makeSG()
}

// this method translates spec to a collection of Security Groups
// 1. generate SGs for relevant endpoints for each connection
// 2. generate SGs for blocked endpoints (endpoints that do not appear in Spec)
func (s *SGSynthesizer) makeSG() *ir.SGCollection {
	for c := range s.spec.Connections {
		s.generateSGRulesFromConnection(&s.spec.Connections[c])
	}
	s.generateSGsForBlockedResources()
	return s.result
}

//  1. check that both resources are supported in SG generation.
//  2. check that at least one resource is internal.
//  3. convert src and dst resources to namedAddrs slices to make it more convenient to go through the addrs
//     and add the rule to the relevant SG.
//  4. generate rules and add them to relevant SG to allow traffic for all pairs of IPAddrs of both resources.
func (s *SGSynthesizer) generateSGRulesFromConnection(conn *ir.Connection) {
	if !resourceRelevantToSG(conn.Src.Type) {
		log.Fatalf(SGTypeNotSupported, string(conn.Src.Type))
	}
	if !resourceRelevantToSG(conn.Dst.Type) {
		log.Fatalf(SGTypeNotSupported, string(conn.Dst.Type))
	}
	internalSrc, internalDst, _ := internalConn(conn)
	if !internalSrc && !internalDst {
		log.Fatalf("SG: Both source and destination are external for connection %v", *conn)
	}

	srcEndpoints := updateEndpoints(&s.spec.Defs, conn.Src)
	dstEndpoints := updateEndpoints(&s.spec.Defs, conn.Dst)

	for _, srcEndpoint := range srcEndpoints {
		for _, dstEndpoint := range dstEndpoints {
			if srcEndpoint.Addrs.Equal(dstEndpoint.Addrs) {
				continue
			}

			for _, trackedProtocol := range conn.TrackedProtocols {
				s.allowConnectionFromSrc(conn, trackedProtocol, srcEndpoint, dstEndpoint)
				s.allowConnectionToDst(conn, trackedProtocol, srcEndpoint, dstEndpoint)
			}
		}
	}
}

// if the src in internal, a rule will be created to allow traffic.
func (s *SGSynthesizer) allowConnectionFromSrc(conn *ir.Connection, trackedProtocol ir.TrackedProtocol,
	srcEndpoint, dstEndpoint *namedAddrs) {
	internalSrc, _, internal := internalConn(conn)

	if !internalSrc {
		return
	}
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}.String()
	sgSrcName := ir.SGName(srcEndpoint.Name)
	sgSrc := s.result.LookupOrCreate(sgSrcName)
	sgSrc.Attached = []ir.ID{ir.ID(sgSrcName)}
	rule := &ir.SGRule{
		Remote:      sgRemote(&s.spec.Defs, dstEndpoint),
		Direction:   ir.Outbound,
		Protocol:    trackedProtocol.Protocol,
		Explanation: reason,
	}
	sgSrc.Add(rule)
}

// if the dst in internal, a rule will be created to allow traffic.
func (s *SGSynthesizer) allowConnectionToDst(conn *ir.Connection, trackedProtocol ir.TrackedProtocol,
	srcEndpoint, dstEndpoint *namedAddrs) {
	_, internalDst, internal := internalConn(conn)

	if !internalDst {
		return
	}
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}.String()
	sgDstName := ir.SGName(dstEndpoint.Name)
	sgDst := s.result.LookupOrCreate(sgDstName)
	sgDst.Attached = []ir.ID{ir.ID(sgDstName)}

	// udp protocol does not have inverse direction
	inverseP := trackedProtocol.Protocol.InverseDirection()
	if p, ok := trackedProtocol.Protocol.(netp.TCPUDP); ok && p.ProtocolString() == netp.ProtocolStringUDP {
		inverseP, _ = netp.NewTCPUDP(false, netp.MinPort, netp.MaxPort, netp.MinPort, netp.MaxPort)
	}

	rule := &ir.SGRule{
		Remote:      sgRemote(&s.spec.Defs, srcEndpoint),
		Direction:   ir.Inbound,
		Protocol:    inverseP,
		Explanation: reason,
	}
	sgDst.Add(rule)
}

// generate SGs for blocked endpoints (endpoints that do not appear in Spec)
func (s *SGSynthesizer) generateSGsForBlockedResources() {
	blockedResources := s.spec.ComputeBlockedResources()
	for _, resource := range blockedResources {
		sg := s.result.LookupOrCreate(ir.SGName(resource)) // an empty SG allows no connections
		sg.Attached = []ir.ID{resource}
	}
}

// convert src and dst resources to namedAddrs slices to make it more convenient to go through the addrs and add
// the rule to the relevant sg.
func updateEndpoints(s *ir.Definitions, resource ir.Resource) []*namedAddrs {
	if resource.Type == ir.ResourceTypeExternal {
		return []*namedAddrs{{Name: resource.Name, Addrs: resource.IPAddrs[0]}}
	}
	result := make([]*namedAddrs, len(resource.IPAddrs))
	for i, ip := range resource.IPAddrs {
		name := resource.Name
		// TODO: delete the following if statement when there will be a support in VSIs with multiple NIFs
		if nifDetails, ok := s.NIFs[resource.Name]; ok {
			name = nifDetails.Instance
		}
		result[i] = &namedAddrs{Name: name, Addrs: ip}
	}
	return result
}

func sgRemote(s *ir.Definitions, resource *namedAddrs) ir.RemoteType {
	if _, ok := s.Externals[resource.Name]; ok {
		return resource.Addrs
	}
	return ir.SGName(resource.Name)
}

func resourceRelevantToSG(e ir.ResourceType) bool {
	return e == ir.ResourceTypeNIF || e == ir.ResourceTypeExternal || e == ir.ResourceTypeVPE
}
