// Package synth takes global specification written in a JSON file, as described by spec_schema.input,
// and generates NetworkACLs that collectively enable the connectivity described in the global specification.
package synth

import (
	"encoding/json"
	"fmt"
	"os"
)

// UnmarshalSpec returns a Spec struct given a file adhering to spec_schema.input
func UnmarshalSpec(filename string) (*Spec, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	spec := new(Spec)
	err = json.Unmarshal(bytes, spec)
	if err != nil {
		return nil, err
	}
	for _, conn := range spec.RequiredConnections {
		err = fixProtocolList(conn.AllowedProtocols)
		if err != nil {
			return nil, err
		}
	}
	for _, section := range spec.Sections {
		if !section.FullyConnected {
			continue
		}
		err = fixProtocolList(section.FullyConnectedWithConnectionType)
		if err != nil {
			return nil, err
		}
	}
	return spec, err
}

func fixProtocolList(list ProtocolList) error {
	for j := range list {
		var err error
		list[j], err = fixProtocol(list[j])
		if err != nil {
			return err
		}
	}
	return nil
}

func fixProtocol(unparsed interface{}) (interface{}, error) {
	p := unparsed.(map[string]interface{})
	switch p["protocol"] {
	case "TCP":
		return unmarshalProtocol(p, new(TcpUdp))
	case "UDP":
		return unmarshalProtocol(p, new(TcpUdp))
	case "ICMP":
		return unmarshalProtocol(p, new(Icmp))
	default:
		panic(fmt.Sprintf("Impossible protocol name: %q, %T", p["protocol"], p["protocol"]))
	}
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
