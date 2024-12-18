/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package jsonio handles global specification written in a JSON file, as described by spec_schema.input
package jsonio

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/spec"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// translateConnections translate required connections from spec.Spec to []*ir.Connection
func (r *Reader) translateConnections(conns []spec.SpecRequiredConnectionsElem, defs *ir.Definitions,
	blockedResources *ir.BlockedResources, isSG bool) ([]*ir.Connection, error) {
	var res []*ir.Connection
	for i := range conns {
		connections, err := translateConnection(defs, blockedResources, &conns[i], i, isSG)
		if err != nil {
			return nil, err
		}
		res = slices.Concat(res, connections)
	}
	return res, nil
}

func translateConnection(defs *ir.Definitions, blockedResources *ir.BlockedResources, conn *spec.SpecRequiredConnectionsElem,
	connIdx int, isSG bool) ([]*ir.Connection, error) {
	protocols, err1 := translateProtocols(conn.AllowedProtocols)
	srcResource, isSrcExternal, err2 := translateConnectionResource(defs, blockedResources, &conn.Src, isSG)
	dstResource, isDstExternal, err3 := translateConnectionResource(defs, blockedResources, &conn.Dst, isSG)
	if err := errors.Join(err1, err2, err3); err != nil {
		return nil, err
	}
	if isSrcExternal && isDstExternal {
		return nil, fmt.Errorf("both source (%s) and destination (%s) are external in required connection", conn.Src.Name, conn.Dst.Name)
	}

	origin := connectionOrigin{
		connectionIndex: connIdx,
		srcName:         resourceName(conn.Src),
		dstName:         resourceName(conn.Dst),
	}
	out := &ir.Connection{Src: srcResource, Dst: dstResource, TrackedProtocols: protocols, Origin: origin}
	if conn.Bidirectional {
		backOrigin := origin
		backOrigin.inverse = true
		in := &ir.Connection{Src: dstResource, Dst: srcResource, TrackedProtocols: protocols, Origin: &backOrigin}
		return []*ir.Connection{out, in}, nil
	}
	return []*ir.Connection{out}, nil
}

func translateConnectionResource(defs *ir.Definitions, blockedResources *ir.BlockedResources, resource *spec.Resource,
	isSG bool) (r *ir.ConnectedResource, isExternal bool, err error) {
	resourceType, err := translateResourceType(defs, resource)
	if err != nil {
		return nil, false, err
	}
	var res *ir.ConnectedResource
	if isSG {
		res, err = defs.LookupForSGSynth(resourceType, resource.Name)
		updateBlockedResourcesSGSynth(blockedResources, res)
	} else {
		res, err = defs.LookupForACLSynth(resourceType, resource.Name)
		updateBlockedResourcesACLSynth(blockedResources, res)
	}
	return res, resourceType == ir.ResourceTypeExternal, err
}

func translateProtocols(protocols spec.ProtocolList) ([]*ir.TrackedProtocol, error) {
	var result = make([]*ir.TrackedProtocol, len(protocols))
	for i, _p := range protocols {
		res := &ir.TrackedProtocol{Origin: protocolOrigin{protocolIndex: i}}
		switch p := _p.(type) {
		case spec.AnyProtocol:
			if len(protocols) != 1 {
				log.Println("when allowing any protocol, there is no need in other protocols")
			}
			res.Protocol = netp.AnyProtocol{}
		case spec.Icmp:
			icmp, err := netp.ICMPFromTypeAndCode(p.Type, p.Code)
			if err != nil {
				return nil, err
			}
			res.Protocol = icmp
		case spec.TcpUdp:
			tcpudp, err := netp.NewTCPUDP(p.Protocol == spec.TcpUdpProtocolTCP, p.MinSourcePort, p.MaxSourcePort,
				p.MinDestinationPort, p.MaxDestinationPort)
			if err != nil {
				return nil, err
			}
			res.Protocol = tcpudp
		default:
			return nil, fmt.Errorf("unsupported protocol: %v", p)
		}
		result[i] = res
	}
	return result, nil
}

func translateResourceType(defs *ir.Definitions, resource *spec.Resource) (ir.ResourceType, error) {
	switch resource.Type {
	case spec.ResourceTypeExternal:
		return ir.ResourceTypeExternal, nil
	case spec.ResourceTypeSubnet:
		return ir.ResourceTypeSubnet, nil
	case spec.ResourceTypeNif:
		return ir.ResourceTypeNIF, nil
	case spec.ResourceTypeInstance:
		return ir.ResourceTypeInstance, nil
	case spec.ResourceTypeVpe:
		return ir.ResourceTypeVPE, nil
	case spec.ResourceTypeSegment:
		if _, ok := defs.SubnetSegments[resource.Name]; ok {
			return ir.ResourceTypeSubnetSegment, nil
		}
		if _, ok := defs.CidrSegments[resource.Name]; ok {
			return ir.ResourceTypeCidrSegment, nil
		}
		if _, ok := defs.NifSegments[resource.Name]; ok {
			return ir.ResourceTypeNifSegment, nil
		}
		if _, ok := defs.InstanceSegments[resource.Name]; ok {
			return ir.ResourceTypeInstanceSegment, nil
		}
		if _, ok := defs.VpeSegments[resource.Name]; ok {
			return ir.ResourceTypeVpeSegment, nil
		}
		return ir.ResourceTypeSubnet, fmt.Errorf("could not find segment %v", resource.Name)
	}
	return ir.ResourceTypeSubnet, fmt.Errorf("unsupported resource type %v (%v)", resource.Type, resource.Name)
}

func updateBlockedResourcesSGSynth(blockedResources *ir.BlockedResources, resource *ir.ConnectedResource) {
	for _, namedAddrs := range resource.CidrsWhenLocal {
		// should check also resource type
		if _, ok := blockedResources.BlockedInstances[namedAddrs.Name]; ok {
			blockedResources.BlockedInstances[namedAddrs.Name] = false
		}
		if _, ok := blockedResources.BlockedVPEs[namedAddrs.Name]; ok {
			blockedResources.BlockedVPEs[namedAddrs.Name] = false
		}
	}
}

func updateBlockedResourcesACLSynth(blockedResources *ir.BlockedResources, resource *ir.ConnectedResource) {
	for _, namedAddrs := range resource.CidrsWhenLocal {
		if _, ok := blockedResources.BlockedSubnets[namedAddrs.Name]; ok {
			blockedResources.BlockedSubnets[namedAddrs.Name] = false
		}
	}
}
