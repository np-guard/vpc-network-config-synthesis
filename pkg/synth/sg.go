/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package synth

import (
	"github.com/np-guard/models/pkg/netp"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const SGTypeNotSupported = "SG: src/dst of type %s is not supported"

type SGSynthesizer struct {
	spec   *ir.Spec
	result *ir.SGCollection
}

// NewSGSynthesizer creates and returns a new SGSynthesizer instances
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
		s.generateSGRulesFromConnection(s.spec.Connections[c])
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
	internalSrc, internalDst, internalConn := internalConn(conn)

	for _, srcEndpoint := range conn.Src.NamedAddrs {
		for _, dstCidr := range conn.Dst.Cidrs {
			if srcEndpoint.IPAddrs.Equal(dstCidr.IPAddrs) {
				continue
			}

			for _, trackedProtocol := range conn.TrackedProtocols {
				ruleExplanation := explanation{internal: internalConn, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}.String()
				s.allowConnectionEndpoint(srcEndpoint, dstCidr, trackedProtocol.Protocol, ir.Outbound, internalSrc, ruleExplanation)
			}
		}
	}

	for _, dstEndpoint := range conn.Dst.NamedAddrs {
		for _, srcCidr := range conn.Src.Cidrs {
			if dstEndpoint.IPAddrs.Equal(srcCidr.IPAddrs) {
				continue
			}

			for _, trackedProtocol := range conn.TrackedProtocols {
				ruleExplanation := explanation{internal: internalConn, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}.String()
				s.allowConnectionEndpoint(dstEndpoint, srcCidr, trackedProtocol.Protocol, ir.Inbound, internalDst, ruleExplanation)
			}
		}
	}
}

// if the endpoint in internal, a rule will be created to allow traffic.
func (s *SGSynthesizer) allowConnectionEndpoint(localEndpoint, remoteEndpoint *ir.NamedAddrs, p netp.Protocol, direction ir.Direction,
	internalEndpoint bool, ruleExplanation string) {
	if !internalEndpoint {
		return
	}
	localSGName := ir.SGName(*localEndpoint.Name)
	localSG := s.result.LookupOrCreate(localSGName)
	localSG.Attached = []ir.ID{ir.ID(localSGName)}
	rule := &ir.SGRule{
		Remote:      sgRemote(s.spec.Defs, remoteEndpoint),
		Direction:   direction,
		Protocol:    p,
		Explanation: ruleExplanation,
	}
	localSG.Add(rule)
}

// what to do if its a cidr segment?
func sgRemote(s *ir.Definitions, resource *ir.NamedAddrs) ir.RemoteType {
	if _, ok := s.Externals[*resource.Name]; ok {
		return resource.IPAddrs
	}
	if _, ok := s.Subnets[*resource.Name]; ok {
		return resource.IPAddrs
	}
	// what to do if its a cidr?
	return ir.SGName(*resource.Name)
}

// generate SGs for blocked endpoints (endpoints that do not appear in Spec)
func (s *SGSynthesizer) generateSGsForBlockedResources() {
	blockedResources := s.spec.ComputeBlockedResources()
	for _, resource := range blockedResources {
		sg := s.result.LookupOrCreate(ir.SGName(resource)) // an empty SG allows no connections
		sg.Attached = []ir.ID{resource}
	}
}
