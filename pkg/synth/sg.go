/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package synth

import (
	"fmt"
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const SGTypeNotSupported = "SG: src/dst of type %s is not supported."

type SGSynthesizer struct {
	Spec   *ir.Spec
	Result *ir.SGCollection
}

// MakeSG translates Spec to a collection of security groups
// 1. generate SGs for relevant endpoints for each connection
// 2. generate SGs for blocked endpoints (endpoints that do not appear in Spec)
func (s *SGSynthesizer) MakeSG() *ir.SGCollection {
	for c := range s.Spec.Connections {
		s.generateSGRulesFromConnection(&s.Spec.Connections[c])
	}
	s.generateSGsForBlockedResources()
	return s.Result
}

//  1. check that both resources are supported in SG generation.
//  2. check that at least one resource is internal.
//  3. convert src and dst resources to namedAddrs slices to make it more convenient to go through the addrs
//     and add the rule to the relevant SG.
//  4. generate rules and add them to relevant SG to allow traffic for all pairs of IPAddrs of both resources.
func (s *SGSynthesizer) generateSGRulesFromConnection(conn *ir.Connection) {
	if !resourceRelevantToSG(conn.Src.Type) {
		log.Fatalf(fmt.Sprintf(SGTypeNotSupported, string(conn.Src.Type)))
	}
	if !resourceRelevantToSG(conn.Dst.Type) {
		log.Fatalf(fmt.Sprintf(SGTypeNotSupported, string(conn.Dst.Type)))
	}
	internalSrc, internalDst, _ := internalConn(conn)
	if !internalSrc && !internalDst {
		log.Fatalf("SG: Both source and destination are external for connection %v", *conn)
	}

	srcEndpoints := updateEndpoints(&s.Spec.Defs, conn.Src)
	dstEndpoints := updateEndpoints(&s.Spec.Defs, conn.Dst)

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
	reason := explanation{internal: internal, connectionOrigin: conn.Origin, protocolOrigin: trackedProtocol.Origin}.String()

	if internalSrc {
		sgSrcName := ir.SGName(srcEndpoint.Name)
		sgSrc := s.Result.LookupOrCreate(sgSrcName)
		sgSrc.Attached = []ir.ID{ir.ID(sgSrcName)}
		rule := &ir.SGRule{
			Remote:      sgRemote(&s.Spec.Defs, dstEndpoint),
			Direction:   ir.Outbound,
			Protocol:    trackedProtocol.Protocol,
			Explanation: reason,
		}
		sgSrc.Add(rule)
	}
}

// if the dst in internal, a rule will be created to allow traffic.
func (s *SGSynthesizer) allowConnectionToDst(conn *ir.Connection, trackedProtocol ir.TrackedProtocol,
	srcEndpoint, dstEndpoint *namedAddrs) {
	_, internalDst, internal := internalConn(conn)
	reason := explanation{
		internal:         internal,
		connectionOrigin: conn.Origin,
		protocolOrigin:   trackedProtocol.Origin,
	}.String()

	if internalDst {
		sgDstName := ir.SGName(dstEndpoint.Name)
		sgDst := s.Result.LookupOrCreate(sgDstName)
		sgDst.Attached = []ir.ID{ir.ID(sgDstName)}
		rule := &ir.SGRule{
			Remote:      sgRemote(&s.Spec.Defs, srcEndpoint),
			Direction:   ir.Inbound,
			Protocol:    trackedProtocol.Protocol.InverseDirection(),
			Explanation: reason,
		}
		sgDst.Add(rule)
	}
}

// generate SGs for blocked endpoints (endpoints that do not appear in Spec)
func (s *SGSynthesizer) generateSGsForBlockedResources() {
	blockedResources := s.Spec.ComputeBlockedResources()
	for _, resource := range blockedResources {
		sg := s.Result.LookupOrCreate(ir.SGName(resource)) // an empty SG allows no connections
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
