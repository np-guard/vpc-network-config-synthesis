// Package jsonio handles global specification written in a JSON file, as described by spec_schema.input
package jsonio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/np-guard/models/pkg/ipblocks"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Reader implements ir.Reader
type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (*Reader) ReadSpec(filename string, configDefs *ir.ConfigDefs) (*ir.Spec, error) {
	jsonspec, err := unmarshal(filename)
	if err != nil {
		return nil, err
	}

	err = validateSegments(jsonspec.Segments)
	if err != nil {
		return nil, err
	}

	cidrSegments := translateSegments(jsonspec.Segments, TypeCidr)
	finalCidrSegments, err := computeContainedSubentsInCidrSegment(configDefs, cidrSegments)
	if err != nil {
		return nil, err
	}

	if configDefs == nil {
		configDefs = &ir.ConfigDefs{
			Subnets:        translateIPMap(jsonspec.Subnets),
			NIFToIP:        translateIPMap(jsonspec.Nifs),
			InstanceToNIFs: jsonspec.Instances,
		}
	}
	defs := &ir.Definitions{
		ConfigDefs:     *configDefs,
		SubnetSegments: translateSegments(jsonspec.Segments, TypeSubnet),
		CidrSegments:   finalCidrSegments,
		Externals:      translateIPMap(jsonspec.Externals),
	}

	var connections []ir.Connection
	for i := range jsonspec.RequiredConnections {
		bidiconns, err := translateConnection(defs, &jsonspec.RequiredConnections[i], i)
		if err != nil {
			return nil, err
		}
		connections = append(connections, bidiconns...)
	}

	return &ir.Spec{
		Connections: connections,
		Defs:        *defs,
	}, nil
}

func validateSegments(jsonssegments SpecSegments) error {
	for _, v := range jsonssegments {
		if v.Type != TypeSubnet && v.Type != TypeCidr {
			return fmt.Errorf("only subnet and cidr segments are supported, not %q", v.Type)
		}
	}
	return nil
}

func translateSegments(jsonSegments SpecSegments, segnmentType Type) map[string][]string {
	result := make(map[string][]string)
	for k, v := range jsonSegments {
		if v.Type == segnmentType {
			result[k] = v.Items
		}
	}
	return result
}

func computeContainedSubentsInCidrSegment(configDefs *ir.ConfigDefs, m map[string][]string) (map[string]map[string][]string, error) {
	finalMap := make(map[string]map[string][]string)
	for segmentName, segment := range m {
		// each cidr saves the contained subnets
		segmentMap := make(map[string][]string)
		for _, cidr := range segment {
			c, err := ipblocks.NewIPBlockFromCidrOrAddress(cidr)
			if err != nil {
				return nil, err
			}
			if !cidrContainedInVpc(*c, configDefs.AddressPrefixes) {
				return nil, fmt.Errorf("%s is not contained in the vpc", cidr)
			}
			subnets, err := configDefs.SubnetsContainedInCidr(*c)
			if err != nil {
				return nil, err
			}
			segmentMap[cidr] = subnets
		}
		finalMap[segmentName] = segmentMap
	}
	return finalMap, nil
}

func translateIPMap(m map[string]string) map[string]ir.IP {
	res := make(map[string]ir.IP)
	for k, v := range m {
		res[k] = ir.IPFromString(v)
	}
	return res
}

func translateConnection(defs *ir.Definitions, v *SpecRequiredConnectionsElem, connectionIndex int) ([]ir.Connection, error) {
	p, err := translateProtocols(v.AllowedProtocols)
	if err != nil {
		return nil, err
	}
	srcEndpointType, err := translateEndpointType(v.Src.Type)
	if err != nil {
		return nil, err
	}
	src, err := defs.Lookup(srcEndpointType, v.Src.Name)
	if err != nil {
		return nil, err
	}
	dstEndpointType, err := translateEndpointType(v.Dst.Type)
	if err != nil {
		return nil, err
	}
	dst, err := defs.Lookup(dstEndpointType, v.Dst.Name)
	if err != nil {
		return nil, err
	}

	origin := connectionOrigin{
		connectionIndex: connectionIndex,
		srcName:         endpointName(v.Src),
		dstName:         endpointName(v.Dst),
	}
	out := ir.Connection{Src: src, Dst: dst, TrackedProtocols: p, Origin: origin}
	if v.Bidirectional {
		backOrigin := origin
		backOrigin.inverse = true
		in := ir.Connection{Src: dst, Dst: src, TrackedProtocols: p, Origin: &backOrigin}
		return []ir.Connection{out, in}, nil
	}
	return []ir.Connection{out}, nil
}

func translateProtocols(protocols ProtocolList) ([]ir.TrackedProtocol, error) {
	var result = make([]ir.TrackedProtocol, len(protocols))
	for i, _p := range protocols {
		result[i].Origin = protocolOrigin{protocolIndex: i}
		switch p := _p.(type) {
		case AnyProtocol:
			if len(protocols) != 1 {
				return nil, fmt.Errorf("when allowing any protocol, no more protocols can be defined")
			}
			result[i].Protocol = ir.AnyProtocol{}
		case Icmp:
			if p.Type == nil {
				if p.Code != nil {
					return nil, fmt.Errorf("defnining ICMP code for unspecified ICMP type is not allowed")
				}
				result[i].Protocol = ir.TrackedProtocol{Protocol: ir.ICMP{}}
			} else {
				err := ir.ValidateICMP(*p.Type, *p.Code)
				if err != nil {
					return nil, err
				}
				result[i].Protocol = ir.ICMP{ICMPCodeType: &ir.ICMPCodeType{Type: *p.Type, Code: p.Code}}
			}
		case TcpUdp:
			result[i].Protocol = ir.TCPUDP{
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
	case EndpointTypeNif:
		return ir.EndpointTypeNIF, nil
	case EndpointTypeInstance:
		return ir.EndpointTypeInstance, nil
	case EndpointTypeVpe:
		return ir.EndpointTypeVPE, nil
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

func cidrContainedInVpc(cidr ipblocks.IPBlock, addressPrefixes []ipblocks.IPBlock) bool {
	for i := range addressPrefixes {
		if cidr.ContainedIn(&addressPrefixes[i]) {
			return true
		}
	}
	return false
}
