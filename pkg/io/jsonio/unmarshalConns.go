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

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/spec"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// transalteConnections translate requires connections from spec.Spec to []*ir.Connection
func (r *Reader) transalteConnections(conns []spec.SpecRequiredConnectionsElem, defs *ir.Definitions, isSG bool) ([]*ir.Connection, error) {
	var connections []*ir.Connection
	for i := range conns {
		conns, err := translateConnection(defs, &conns[i], i, isSG)
		if err != nil {
			return nil, err
		}
		connections = append(connections, conns...)
	}
	return connections, nil
}

func translateConnection(defs *ir.Definitions, conn *spec.SpecRequiredConnectionsElem, connIdx int, isSG bool) ([]*ir.Connection, error) {
	protocols, err1 := translateProtocols(conn.AllowedProtocols)
	srcResource, isSrcExternal, err2 := transalteConnectionResource(defs, &conn.Src, isSG)
	dstResource, isDstExternal, err3 := transalteConnectionResource(defs, &conn.Dst, isSG)
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

func transalteConnectionResource(defs *ir.Definitions, resource *spec.Resource, isSG bool) (r *ir.Resource, isExternal bool, err error) {
	resourceType, err := translateResourceType(defs, resource)
	if err != nil {
		return nil, false, err
	}
	var res *ir.Resource
	if isSG {
		res, err = defs.LookupForSGSynth(resourceType, resource.Name)
		updateBlockedResourcesSGSynth(defs, res)
	} else {
		res, err = defs.LookupForACLSynth(resourceType, resource.Name)
		updateBlockedResourcesACLSynth(defs, res)
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
			return nil, fmt.Errorf("impossible protocol: %v", p)
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
	}
	return ir.ResourceTypeSubnet, fmt.Errorf("unsupported resource type %v (%v)", resource.Type, resource.Name)
}

func updateBlockedResourcesSGSynth(defs *ir.Definitions, resource *ir.Resource) {
	for _, namedAddrs := range resource.NamedAddrs {
		if _, ok := defs.BlockedInstances[*namedAddrs.Name]; ok {
			defs.BlockedInstances[*namedAddrs.Name] = false
		}
		if _, ok := defs.BlockedVPEs[*namedAddrs.Name]; ok {
			defs.BlockedVPEs[*namedAddrs.Name] = false
		}
	}
}

func updateBlockedResourcesACLSynth(defs *ir.Definitions, resource *ir.Resource) {
	for _, namedAddrs := range resource.NamedAddrs {
		if _, ok := defs.BlockedSubnets[*namedAddrs.Name]; ok {
			defs.BlockedSubnets[*namedAddrs.Name] = false
		}
	}
}
