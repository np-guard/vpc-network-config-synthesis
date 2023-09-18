package synth

import (
	"log"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// MakeSG translates Spec to a collection of security groups
func MakeSG(s *ir.Spec, opt Options) *ir.SGCollection {
	collections := []*ir.SGCollection{}
	for c := range s.Connections {
		collection := GenerateSGCollectionFromConnection(&s.Connections[c], s.Defs.RemoteFromIP)
		collections = append(collections, collection)
	}
	return ir.MergeSGCollections(collections...)
}

func GenerateSGCollectionFromConnection(conn *ir.Connection, sgSelector func(target ir.IP) ir.RemoteType) *ir.SGCollection {
	internalSrc := conn.Src.Type != ir.EndpointTypeExternal
	internalDst := conn.Dst.Type != ir.EndpointTypeExternal
	if !internalSrc && !internalDst {
		log.Fatalf("SG: Both source and destination are external for connection %v", *conn)
	}

	result := ir.NewSGCollection()

	if !sgTrigger(conn.Src.Type) && !sgTrigger(conn.Dst.Type) {
		return result
	}

	for _, src := range conn.Src.Values {
		for _, dst := range conn.Dst.Values {
			if src == dst {
				continue
			}

			for _, trackedProtocol := range conn.TrackedProtocols {
				reason := explanation{
					internal:         internalSrc && internalDst,
					connectionOrigin: conn.Origin,
					protocolOrigin:   trackedProtocol.Origin,
				}.String()

				if internalSrc {
					sgSrcName, ok := sgSelector(src).(ir.SGName)
					if !ok {
						log.Panicf("Source is not security group name: %v", src)
					}
					sgSrc := result.LookupOrCreate(sgSrcName)
					sgSrc.Attached = []ir.SGName{sgSrcName}
					sgSrc.Rules = append(sgSrc.Rules, ir.SGRule{
						Remote:      sgSelector(dst),
						Direction:   ir.Outbound,
						Protocol:    trackedProtocol.Protocol,
						Explanation: reason,
					})
				}
				if internalDst {
					sgDstName, ok := sgSelector(dst).(ir.SGName)
					if !ok {
						log.Panicf("Dst is not security group name: %v", dst)
					}
					sgDst := result.LookupOrCreate(sgDstName)
					sgDst.Attached = []ir.SGName{sgDstName}
					sgDst.Rules = append(sgDst.Rules, ir.SGRule{
						Remote:      sgSelector(src),
						Direction:   ir.Inbound,
						Protocol:    trackedProtocol.Protocol.InverseDirection(),
						Explanation: reason,
					})
				}
			}
		}
	}

	return result
}

func sgTrigger(e ir.EndpointType) bool {
	return e == ir.EndpointTypeNif || e == ir.EndpointTypeInstance
}
