// Package spec takes global specification written in a JSON file, as described by spec_schema.input
package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/acl"
)

// Unmarshal returns a Spec struct given a file adhering to spec_schema.input
func Unmarshal(filename string) (*Spec, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	spec := new(Spec)
	err = json.Unmarshal(bytes, spec)
	if err != nil {
		return nil, err
	}
	for i := range spec.RequiredConnections {
		if spec.RequiredConnections[i].AllowedProtocols == nil {
			spec.RequiredConnections[i].AllowedProtocols = ProtocolList{new(AnyProtocol)}
		} else {
			err = fixProtocolList(spec.RequiredConnections[i].AllowedProtocols)
			if err != nil {
				return nil, err
			}
		}
	}
	return spec, err
}

func (s *Spec) SetSubnets(subnets map[string]string) error {
	if subnets != nil {
		if len(s.Subnets) != 0 {
			return errors.New("both subnets and config_file are supplied")
		}
		s.Subnets = subnets
	}
	return nil
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
			err = validateICMP(icmp)
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

func validateICMP(icmp *Icmp) error {
	if icmp.Type == nil {
		if icmp.Code != nil {
			return errors.New("cannot define ICMP code for unknown ICMP type")
		}
		return nil
	}
	return acl.ValidateICMP(*icmp.Type, *icmp.Code)
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
