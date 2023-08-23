// Package jsonio handles global specification written in a JSON file, as described by spec_schema.input
package jsonio

import (
	"encoding/json"
	"fmt"
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
	srcEndpointType, err := translateEndpointType(v.Src.Type)
	if err != nil {
		return nil, err
	}
	src, err := defs.Lookup(v.Src.Name, srcEndpointType)
	if err != nil {
		return nil, err
	}
	dstEndpointType, err := translateEndpointType(v.Dst.Type)
	if err != nil {
		return nil, err
	}
	dst, err := defs.Lookup(v.Dst.Name, dstEndpointType)
	if err != nil {
		return nil, err
	}
	result := ir.MakeConnection(src, dst, p, v.Bidirectional)
	return result, nil
}

func translateProtocols(protocols ProtocolList) ([]ir.Protocol, error) {
	var result = make([]ir.Protocol, len(protocols))
	for i, _p := range protocols {
		switch p := _p.(type) {
		case AnyProtocol:
			if i != 0 {
				return nil, fmt.Errorf("when allowing any protocol, no more protocols can be defined")
			}
			return []ir.Protocol{}, nil
		case Icmp:
			if p.Type == nil {
				if p.Code != nil {
					return nil, fmt.Errorf("defnining ICMP code for unspecified ICMP type is not allowed")
				}
				result[i] = ir.ICMP{}
			} else {
				err := ir.ValidateICMP(*p.Type, *p.Code)
				if err != nil {
					return nil, err
				}
				result[i] = ir.ICMP{ICMPCodeType: &ir.ICMPCodeType{Type: *p.Type, Code: p.Code}}
			}
		case TcpUdp:
			result[i] = ir.TCPUDP{
				Protocol: ir.TransportLayerProtocolName(p.Protocol),
				PortRangePair: ir.PortRangePair{
					SrcPort: ir.PortRange{Min: p.MinSourcePort, Max: p.MaxSourcePort},
					DstPort: ir.PortRange{Min: p.MinDestinationPort, Max: p.MaxDestinationPort},
				},
			}
		default:
			return nil, fmt.Errorf("impossible protocol: %v", p)
		}
	}
	return result, nil
}

func translateEndpointType(endpointType EndpointType) (ir.EndpointType, error) {
	switch endpointType {
	case EndpointTypeExternal:
		return ir.EndpointTypeExternal, nil
	case EndpointTypeSegment:
		return ir.EndpointTypeSegment, nil
	case EndpointTypeSubnet:
		return ir.EndpointTypeSubnet, nil
	default:
		return ir.EndpointTypeSubnet, fmt.Errorf("unsupported endpoint type %v", endpointType)
	}
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
		conn := &jsonspec.RequiredConnections[i]
		if conn.AllowedProtocols == nil {
			conn.AllowedProtocols = ProtocolList{AnyProtocol{}}
		} else {
			for j := range conn.AllowedProtocols {
				p := conn.AllowedProtocols[j].(map[string]interface{})
				bytes, err = json.Marshal(p)
				if err != nil {
					return nil, err
				}
				switch p["protocol"] {
				case "ANY":
					var result AnyProtocol
					err = json.Unmarshal(bytes, &result)
					conn.AllowedProtocols[j] = result
				case "TCP", "UDP":
					var result TcpUdp
					err = json.Unmarshal(bytes, &result)
					conn.AllowedProtocols[j] = result
				case "ICMP":
					var result Icmp
					err = json.Unmarshal(bytes, &result)
					conn.AllowedProtocols[j] = result
				default:
					return nil, fmt.Errorf("invalid protocol type %q", p["protocol"])
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return jsonspec, err
}
