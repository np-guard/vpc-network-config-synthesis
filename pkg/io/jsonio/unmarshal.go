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
	jsonSpec, err := unmarshal(filename)
	if err != nil {
		return nil, err
	}

	err = validateSegments(jsonSpec.Segments)
	if err != nil {
		return nil, err
	}

	if configDefs == nil {
		configDefs = &ir.ConfigDefs{
			Subnets:         translateIPMap(jsonSpec.Subnets),
			NIFToIP:         translateIPMap(jsonSpec.Nifs),
			InstanceToNIFs:  jsonSpec.Instances,
			AddressPrefixes: []ir.CIDR{},
		}
	}

	cidrSegments, err := parseCidrSegments(jsonSpec.Segments, configDefs)
	if err != nil {
		return nil, err
	}

	defs := &ir.Definitions{
		ConfigDefs:     *configDefs,
		SubnetSegments: translateSegments(jsonSpec.Segments, TypeSubnet),
		CidrSegments:   cidrSegments,
		Externals:      translateIPMap(jsonSpec.Externals),
	}

	var connections []ir.Connection
	for i := range jsonSpec.RequiredConnections {
		bidiConns, err := translateConnection(defs, &jsonSpec.RequiredConnections[i], i)
		if err != nil {
			return nil, err
		}
		connections = append(connections, bidiConns...)
	}

	return &ir.Spec{
		Connections: connections,
		Defs:        *defs,
	}, nil
}

func validateSegments(jsonSegments SpecSegments) error {
	for _, v := range jsonSegments {
		if v.Type != TypeSubnet && v.Type != TypeCidr {
			return fmt.Errorf("only subnet and cidr segments are supported, not %q", v.Type)
		}
	}
	return nil
}

func translateSegments(jsonSegments SpecSegments, segmentType Type) map[string][]string {
	result := make(map[string][]string)
	for k, v := range jsonSegments {
		if v.Type == segmentType {
			result[k] = v.Items
		}
	}
	return result
}

func parseCidrSegments(jsonSegments SpecSegments, configDefs *ir.ConfigDefs) (map[string]map[string][]string, error) {
	cidrSegments := translateSegments(jsonSegments, TypeCidr)
	finalMap := make(map[string]map[string][]string)
	for segmentName, segment := range cidrSegments {
		// each cidr saves the contained subnets
		segmentMap := make(map[string][]string)
		for _, cidr := range segment {
			c, err := ipblocks.NewIPBlockFromCidrOrAddress(cidr)
			if err != nil {
				return nil, err
			}
			validCidr, err := cidrContainedInVpc(*c, configDefs.AddressPrefixes)
			if err != nil {
				return nil, err
			}
			if !validCidr {
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
	srcResourceType, err := translateResourceType(v.Src.Type)
	if err != nil {
		return nil, err
	}
	src, err := defs.Lookup(srcResourceType, v.Src.Name)
	if err != nil {
		return nil, err
	}
	dstResourceType, err := translateResourceType(v.Dst.Type)
	if err != nil {
		return nil, err
	}
	dst, err := defs.Lookup(dstResourceType, v.Dst.Name)
	if err != nil {
		return nil, err
	}

	origin := connectionOrigin{
		connectionIndex: connectionIndex,
		srcName:         resourceName(v.Src),
		dstName:         resourceName(v.Dst),
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
					return nil, fmt.Errorf("defining ICMP code for unspecified ICMP type is not allowed")
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

func translateResourceType(resourceType ResourceType) (ir.ResourceType, error) {
	switch resourceType {
	case ResourceTypeExternal:
		return ir.ResourceTypeExternal, nil
	case ResourceTypeSegment:
		return ir.ResourceTypeSegment, nil
	case ResourceTypeSubnet:
		return ir.ResourceTypeSubnet, nil
	case ResourceTypeNif:
		return ir.ResourceTypeNIF, nil
	case ResourceTypeInstance:
		return ir.ResourceTypeInstance, nil
	case ResourceTypeVpe:
		return ir.ResourceTypeVPE, nil
	default:
		return ir.ResourceTypeSubnet, fmt.Errorf("unsupported resource type %v", resourceType)
	}
}

// unmarshal returns a Spec struct given a file adhering to spec_schema.input
func unmarshal(filename string) (*Spec, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	jsonSpec := new(Spec)
	err = json.Unmarshal(bytes, jsonSpec)
	if err != nil {
		return nil, err
	}
	for i := range jsonSpec.RequiredConnections {
		conn := &jsonSpec.RequiredConnections[i]
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
	return jsonSpec, err
}

func cidrContainedInVpc(cidr ipblocks.IPBlock, addressPrefixes []ir.CIDR) (bool, error) {
	for i := range addressPrefixes {
		addressPrefix, err := ipblocks.NewIPBlockFromCidrOrAddress(addressPrefixes[i].String())
		if err != nil {
			return false, err
		}
		if cidr.ContainedIn(addressPrefix) {
			return true, nil
		}
	}
	return false, nil
}
