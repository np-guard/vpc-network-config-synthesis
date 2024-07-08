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

// MakeSG translates Spec to a collection of security groups
func MakeSG(s *ir.Spec) *ir.SGCollection {
	collections := []*ir.SGCollection{}
	for c := range s.Connections {
		collection := generateSGCollectionFromConnection(&s.Defs, &s.Connections[c])
		collections = append(collections, collection)
	}
	collections = append(collections, generateSGCollectionForBlockedResources(s))
	return ir.MergeSGCollections(collections...)
}

func generateSGCollectionFromConnection(s *ir.Definitions, conn *ir.Connection) *ir.SGCollection {
	internalSrc := conn.Src.Type != ir.ResourceTypeExternal
	internalDst := conn.Dst.Type != ir.ResourceTypeExternal
	if !internalSrc && !internalDst {
		log.Fatalf("SG: Both source and destination are external for connection %v", *conn)
	}
	if !resourceRelevantToSG(conn.Src.Type) {
		log.Fatalf(fmt.Sprintf(SGTypeNotSupported, string(conn.Src.Type)))
	}
	if !resourceRelevantToSG(conn.Dst.Type) {
		log.Fatalf(fmt.Sprintf(SGTypeNotSupported, string(conn.Dst.Type)))
	}

	result := ir.NewSGCollection()

	srcEndpoints := updateEndpoints(s, conn.Src)
	dstEndpoints := updateEndpoints(s, conn.Dst)

	for _, srcEndpoint := range srcEndpoints {
		for _, dstEndpoint := range dstEndpoints {
			if srcEndpoint.Addrs.Equal(dstEndpoint.Addrs) {
				continue
			}

			for _, trackedProtocol := range conn.TrackedProtocols {
				reason := explanation{
					internal:         internalSrc && internalDst,
					connectionOrigin: conn.Origin,
					protocolOrigin:   trackedProtocol.Origin,
				}.String()

				if internalSrc {
					sgSrcName := ir.SGName(srcEndpoint.Name)
					sgSrc := result.LookupOrCreate(sgSrcName)
					sgSrc.Attached = []ir.ID{ir.ID(sgSrcName)}
					rule := ir.SGRule{
						Remote:      remote(s, dstEndpoint),
						Direction:   ir.Outbound,
						Protocol:    trackedProtocol.Protocol,
						Explanation: reason,
					}
					sgSrc.Rules = append(sgSrc.Rules, rule)
				}
				if internalDst {
					sgDstName := ir.SGName(dstEndpoint.Name)
					sgDst := result.LookupOrCreate(sgDstName)
					sgDst.Attached = []ir.ID{ir.ID(sgDstName)}
					rule := ir.SGRule{
						Remote:      remote(s, srcEndpoint),
						Direction:   ir.Inbound,
						Protocol:    trackedProtocol.Protocol.InverseDirection(),
						Explanation: reason,
					}
					sgDst.Rules = append(sgDst.Rules, rule)
				}
			}
		}
	}
	return result
}

func generateSGCollectionForBlockedResources(s *ir.Spec) *ir.SGCollection {
	blockedResources := s.ComputeBlockedResources()
	result := ir.NewSGCollection()
	for _, resource := range blockedResources {
		sg := result.LookupOrCreate(ir.SGName(resource)) // an empty SG allows no connections
		sg.Attached = []ir.ID{resource}
	}
	return result
}

func updateEndpoints(s *ir.Definitions, resource ir.Resource) []ir.ConnResource {
	if resource.Type == ir.ResourceTypeExternal {
		return []ir.ConnResource{{Name: resource.Name, Addrs: resource.IPAddrs[0]}}
	}
	result := make([]ir.ConnResource, len(resource.IPAddrs))
	for i, ip := range resource.IPAddrs {
		name := resource.Name
		// TODO: delete the if statement when there will be a support in VSIs with multiple NIFs
		if nifDetails, ok := s.NIFs[resource.Name]; ok {
			name = nifDetails.Instance
		}
		result[i] = ir.ConnResource{Name: name, Addrs: ip}
	}
	return result
}

func remote(s *ir.Definitions, resource ir.ConnResource) ir.RemoteType {
	if _, ok := s.Externals[resource.Name]; ok {
		return resource.Addrs
	}
	return ir.SGName(resource.Name)
}

func resourceRelevantToSG(e ir.ResourceType) bool {
	return e == ir.ResourceTypeNIF || e == ir.ResourceTypeExternal || e == ir.ResourceTypeVPE
}
