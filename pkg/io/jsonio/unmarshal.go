// Package jsonio handles global specification written in a JSON file, as described by spec_schema.input
package jsonio

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/spec"
)

// Reader implements spec.Reader
type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (*Reader) ReadSpec(filename string, subnets map[string]string) (*spec.Spec, error) {
	jsonspec, err := unmarshal(filename)
	if err != nil {
		return nil, err
	}

	result := spec.Spec{
		Externals:      jsonspec.Externals,
		Subnets:        jsonspec.Subnets,
		SubnetSegments: make(map[string][]string),
	}

	if subnets != nil {
		if jsonspec.Subnets != nil && len(jsonspec.Subnets) != 0 {
			return nil, fmt.Errorf("both subnets and config_file are supplied")
		}
		result.Subnets = subnets
	}

	for k, v := range jsonspec.Segments {
		if v.Type != TypeSubnet {
			return nil, fmt.Errorf("only subnet segments are supported, not %q", v.Type)
		}
		result.SubnetSegments[k] = v.Items
	}

	for _, v := range jsonspec.RequiredConnections {
		conn, err := translateConnection(v)
		if err != nil {
			return nil, err
		}
		result.Connections = append(result.Connections, *conn)
	}

	return &result, nil
}

func translateConnection(v SpecRequiredConnectionsElem) (*spec.Connection, error) {
	p, err := translateProtocols(v.AllowedProtocols)
	if err != nil {
		return nil, err
	}
	result := &spec.Connection{
		Bidirectional: v.Bidirectional,
		Dst:           translateEndpoint(v.Dst),
		Src:           translateEndpoint(v.Src),
		Protocols:     p,
	}
	return result, nil
}

func translateProtocols(protocols ProtocolList) ([]spec.Protocol, error) {
	var result = make([]spec.Protocol, len(protocols))
	for i, _p := range protocols {
		var protocol spec.Protocol
		switch p := _p.(type) {
		case *AnyProtocol:
			if i != 0 {
				return nil, fmt.Errorf("when allowing any protocol, no more protocols can be defined")
			}
			return []spec.Protocol{}, nil
		case *Icmp:
			if p.Type == nil {
				if p.Code != nil {
					return nil, fmt.Errorf("defnining ICMP code for unspecified ICMP type is not allowed")
				}
				protocol = spec.ICMP{}
			} else {
				err := spec.ValidateICMP(*p.Type, *p.Code)
				if err != nil {
					return nil, err
				}
				protocol = spec.ICMP{ICMPCodeType: &spec.ICMPCodeType{Type: *p.Type, Code: p.Code}}
			}
		case *TcpUdp:
			protocol = spec.TCPUDP{
				Protocol: spec.TransportLayerProtocolName(p.Protocol),
				PortRangePair: spec.PortRangePair{
					SrcPort: spec.PortRange{Min: p.MinSourcePort, Max: p.MaxSourcePort},
					DstPort: spec.PortRange{Min: p.MinDestinationPort, Max: p.MaxDestinationPort},
				},
			}
		}
		result[i] = protocol
	}
	return result, nil
}

func translateEndpoint(dst *Endpoint) *spec.Endpoint {
	if dst == nil {
		return nil
	}
	return &spec.Endpoint{
		Name: dst.Name,
		Type: spec.EndpointType(dst.Type),
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
