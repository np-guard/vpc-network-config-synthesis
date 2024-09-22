/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package synth

import (
	"fmt"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const SGTypeNotSupported = "SG: src/dst of type %s is not supported"

type SGSynthesizer struct {
	spec   *ir.Spec
	result *ir.SGCollection
}

// NewSGSynthesizer creates and returns a new SGSynthesizer instance
func NewSGSynthesizer(s *ir.Spec, _ bool) Synthesizer {
	return &SGSynthesizer{spec: s, result: ir.NewSGCollection()}
}

func (s *SGSynthesizer) Synth() (ir.Collection, error) {
	return s.makeSG()
}

// this method translates spec to a collection of Security Groups
// 1. generate SGs for relevant endpoints for each connection
// 2. generate SGs for blocked endpoints (endpoints that do not appear in Spec)
func (s *SGSynthesizer) makeSG() (*ir.SGCollection, error) {
	for c := range s.spec.Connections {
		err := s.generateSGRulesFromConnection(&s.spec.Connections[c])
		if err != nil {
			return nil, err
		}
	}
	s.generateSGsForBlockedResources()
	return s.result, nil
}

//  1. check that both resources are supported in SG generation.
//  2. check that at least one resource is internal.
//  3. convert src and dst resources to namedAddrs slices to make it more convenient to go through the addrs
//     and add the rule to the relevant SG.
//  4. generate rules and add them to relevant SG to allow traffic for all pairs of IPAddrs of both resources.
func (s *SGSynthesizer) generateSGRulesFromConnection(conn *ir.Connection) error {
	if !resourceRelevantToSG(conn.Src.Type) {
		return fmt.Errorf(SGTypeNotSupported, string(conn.Src.Type))
	}
	if !resourceRelevantToSG(conn.Dst.Type) {
		return fmt.Errorf(SGTypeNotSupported, string(conn.Dst.Type))
	}
	internalSrc, internalDst, internalConn := internalConn(conn)
	if !internalSrc && !internalDst {
		return fmt.Errorf("SG: Both source and destination are external for connection %v", *conn)
	}

	srcEndpoints := updateEndpoints(&s.spec.Defs, conn.Src)
	dstEndpoints := updateEndpoints(&s.spec.Defs, conn.Dst)

	for _, srcEndpoint := range srcEndpoints {
		for _, dstEndpoint := range dstEndpoints {
			if srcEndpoint.Addrs.Equal(dstEndpoint.Addrs) {
				continue
			}

			for _, trackedProtocol := range conn.TrackedProtocols {
				ruleExplanation := explanation{internal: internalConn, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}.String()
				s.allowConnectionEndpoint(srcEndpoint, dstEndpoint, trackedProtocol.Protocol, ir.Outbound, internalSrc, ruleExplanation)
				s.allowConnectionEndpoint(dstEndpoint, srcEndpoint, trackedProtocol.Protocol, ir.Inbound, internalDst, ruleExplanation)
			}
		}
	}
	return nil
}

// if the endpoint in internal, a rule will be created to allow traffic.
func (s *SGSynthesizer) allowConnectionEndpoint(localEndpoint, remoteEndpoint *namedAddrs, protocol netp.Protocol,
	direction ir.Direction, internalEndpoint bool, ruleExplanation string) {
	if !internalEndpoint {
		return
	}
	localSGName := ir.SGName(localEndpoint.Name)
	localSG := s.result.LookupOrCreate(localSGName)
	localSG.Attached = []ir.ID{ir.ID(localSGName)}
	rule := &ir.SGRule{
		Remote:      sgRemote(&s.spec.Defs, remoteEndpoint),
		Direction:   direction,
		Protocol:    protocol,
		Explanation: ruleExplanation,
	}
	localSG.Add(rule)
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
