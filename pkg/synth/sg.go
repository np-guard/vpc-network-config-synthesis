/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package synth

import (
	"github.com/np-guard/models/pkg/netp"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type SGSynthesizer struct {
	spec   *ir.Spec
	result *ir.SGCollection
}

const warningUnspecifiedSG = "The following endpoints do not have required connections; the generated SGs will block all traffic: "

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
	for _, c := range s.spec.Connections {
		s.generateSGRulesFromConnection(c, ir.Outbound)
		s.generateSGRulesFromConnection(c, ir.Inbound)
	}
	s.generateSGsForBlockedResources()
	return s.result
}

func (s *SGSynthesizer) generateSGRulesFromConnection(conn *ir.Connection, direction ir.Direction) {
	localResource, remoteResource, internalEndpoint, internalConn := connSettings(conn, direction)

	for _, localEndpoint := range localResource.CidrsWhenLocal {
		for _, remoteCidr := range remoteResource.CidrsWhenRemote {
			for _, trackedProtocol := range conn.TrackedProtocols {
				ruleExplanation := explanation{internal: internalConn, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}.String()
				s.allowConnectionEndpoint(localEndpoint, remoteCidr, remoteResource.ResourceType, trackedProtocol.Protocol, direction,
					internalEndpoint, ruleExplanation)
			}
		}
	}
}

// if the endpoint in internal, a rule will be created to allow traffic.
func (s *SGSynthesizer) allowConnectionEndpoint(localEndpoint, remoteEndpoint *ir.NamedAddrs, remoteType ir.ResourceType,
	p netp.Protocol, direction ir.Direction, internalEndpoint bool, ruleExplanation string) {
	if !internalEndpoint {
		return
	}
	localSGName := ir.SGName(*localEndpoint.Name)
	localSG := s.result.LookupOrCreate(localSGName)
	localSG.Attached = []ir.ID{ir.ID(localSGName)}
	rule := &ir.SGRule{
		Remote:      sgRemote(remoteEndpoint, remoteType),
		Direction:   direction,
		Protocol:    p,
		Explanation: ruleExplanation,
	}
	localSG.Add(rule)
}

func sgRemote(resource *ir.NamedAddrs, t ir.ResourceType) ir.RemoteType {
	if isSGRemote(t) {
		return ir.SGName(*resource.Name)
	}
	return resource.IPAddrs
}

func connSettings(conn *ir.Connection, direction ir.Direction) (local, remote *ir.ConnectedResource, internalEndpoint, internalConn bool) {
	internalSrc, internalDst, internalConn := internalConnection(conn)
	local = conn.Src
	remote = conn.Dst
	internalEndpoint = internalSrc
	if direction == ir.Inbound {
		local = conn.Dst
		remote = conn.Src
		internalEndpoint = internalDst
	}
	return
}

func isSGRemote(t ir.ResourceType) bool {
	return t == ir.ResourceTypeInstance || t == ir.ResourceTypeNIF || t == ir.ResourceTypeVPE
}

// generate SGs for blocked endpoints (endpoints that do not appear in Spec)
func (s *SGSynthesizer) generateSGsForBlockedResources() {
	blockedResources := append(utils.TrueKeyValues(s.spec.BlockedInstances), utils.TrueKeyValues(s.spec.BlockedVPEs)...)
	printUnspecifiedWarning(warningUnspecifiedSG, blockedResources)
	for _, resource := range blockedResources {
		sg := s.result.LookupOrCreate(ir.SGName(resource)) // an empty SG allows no connections
		sg.Attached = []ir.ID{resource}
	}
}
