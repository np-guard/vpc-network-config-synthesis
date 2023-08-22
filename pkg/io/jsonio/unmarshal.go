// Package jsonio handles global specification written in a JSON file, as described by spec_schema.input
package jsonio

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Reader implements ir.Reader
type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (*Reader) ReadSpec(filename string, subnets map[string]string) (*ir.Spec, error) {
	jsonspec, err := unmarshal(filename)
	if err != nil {
		return nil, err
	}

	defs := ir.Definitions{
		Externals:      jsonspec.Externals,
		Subnets:        jsonspec.Subnets,
		SubnetSegments: make(map[string][]string),
	}

	if subnets != nil {
		if jsonspec.Subnets != nil && len(jsonspec.Subnets) != 0 {
			return nil, fmt.Errorf("both subnets and config_file are supplied")
		}
		defs.Subnets = subnets
	}

	for k, v := range jsonspec.Segments {
		if v.Type != TypeSubnet {
			return nil, fmt.Errorf("only subnet segments are supported, not %q", v.Type)
		}
		defs.SubnetSegments[k] = v.Items
	}

	var connections []ir.Connection
	for i := range jsonspec.RequiredConnections {
		bidiconns, err := translateConnection(&defs, &jsonspec.RequiredConnections[i])
		if err != nil {
			return nil, err
		}
		connections = append(connections, bidiconns...)
	}

	return &ir.Spec{Connections: connections}, nil
}

func translateConnection(defs *ir.Definitions, v *SpecRequiredConnectionsElem) ([]ir.Connection, error) {
	p, err := translateProtocols(v.AllowedProtocols)
	if err != nil {
		return nil, err
	}
	src, err := defs.Lookup(v.Src.Name, translateEndpointType(v.Src.Type))
	if err != nil {
		return nil, err
	}
	dst, err := defs.Lookup(v.Dst.Name, translateEndpointType(v.Dst.Type))
	if err != nil {
		return nil, err
	}
	result := ir.MakeConnection(src, dst, p, v.Bidirectional)
	return result, nil
}

func translateProtocols(protocols ProtocolList) ([]ir.Protocol, error) {
	var result = make([]ir.Protocol, len(protocols))
	for i, _p := range protocols {
		var protocol ir.Protocol
		switch p := _p.(type) {
		case *AnyProtocol:
			if i != 0 {
				return nil, fmt.Errorf("when allowing any protocol, no more protocols can be defined")
			}
			return []ir.Protocol{}, nil
		case *Icmp:
			if p.Type == nil {
				if p.Code != nil {
					return nil, fmt.Errorf("defnining ICMP code for unspecified ICMP type is not allowed")
				}
				protocol = ir.ICMP{}
			} else {
				err := ir.ValidateICMP(*p.Type, *p.Code)
				if err != nil {
					return nil, err
				}
				protocol = ir.ICMP{ICMPCodeType: &ir.ICMPCodeType{Type: *p.Type, Code: p.Code}}
			}
		case *TcpUdp:
			protocol = ir.TCPUDP{
				Protocol: ir.TransportLayerProtocolName(p.Protocol),
				PortRangePair: ir.PortRangePair{
					SrcPort: ir.PortRange{Min: p.MinSourcePort, Max: p.MaxSourcePort},
					DstPort: ir.PortRange{Min: p.MinDestinationPort, Max: p.MaxDestinationPort},
				},
			}
		}
		result[i] = protocol
	}
	return result, nil
}

func translateEndpointType(endpointType EndpointType) ir.EndpointType {
	switch endpointType {
	case EndpointTypeExternal:
		return ir.EndpointTypeExternal
	case EndpointTypeSegment:
		return ir.EndpointTypeSegment
	case EndpointTypeSubnet:
		return ir.EndpointTypeSubnet
	default:
		log.Fatalf("Unsupported endpoint type %v", endpointType)
	}
	return ir.EndpointTypeSubnet
}

// unmarshal returns a Spec struct given a file adhering to spec_schema.input
func unmarshal(filename string) (*Spec, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	jsonspec := new(Spec)
	err = json.Unmarshal(bytes, jsonspec)
	if err != nil {
		return nil, err
	}
	for i := range jsonspec.RequiredConnections {
		if jsonspec.RequiredConnections[i].AllowedProtocols == nil {
			jsonspec.RequiredConnections[i].AllowedProtocols = ProtocolList{new(AnyProtocol)}
		} else {
			err = fixProtocolList(jsonspec.RequiredConnections[i].AllowedProtocols)
			if err != nil {
				return nil, err
			}
		}
	}
	return jsonspec, err
}

func fixProtocolList(list ProtocolList) error {
	for j := range list {
		var err error
		p := list[j].(map[string]interface{})
		switch p["protocol"] {
		case "TCP":
			list[j], err = unmarshalProtocol(p, new(TcpUdp))
		case "UDP":
			list[j], err = unmarshalProtocol(p, new(TcpUdp))
		case "ICMP":
			var icmp *Icmp
			icmp, err = unmarshalProtocol(p, new(Icmp))
			if err != nil {
				return err
			}
			list[j] = icmp
		case "ANY":
			list[j], err = unmarshalProtocol(p, new(AnyProtocol))
			if err != nil {
				return err
			}
			if len(list) != 1 {
				err = errors.New("redundant protocol declaration")
			}
		default:
			return fmt.Errorf("unknown protocol name: %q, %T", p["protocol"], p["protocol"])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalProtocol[T json.Unmarshaler](p map[string]interface{}, result T) (T, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return result, err
	}
	return result, nil
}
