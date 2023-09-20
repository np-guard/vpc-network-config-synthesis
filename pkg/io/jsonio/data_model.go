// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package jsonio

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type AnyProtocol struct {
	// Necessarily ANY
	Protocol AnyProtocolProtocol `json:"protocol"`
}

type AnyProtocolProtocol string

const AnyProtocolProtocolANY AnyProtocolProtocol = "ANY"

type Endpoint struct {
	// Name of endpoint
	Name string `json:"name"`

	// Type of endpoint
	Type EndpointType `json:"type"`
}

type EndpointType string

const EndpointTypeCidr EndpointType = "cidr"
const EndpointTypeExternal EndpointType = "external"
const EndpointTypeInstance EndpointType = "instance"
const EndpointTypeNif EndpointType = "nif"
const EndpointTypeSegment EndpointType = "segment"
const EndpointTypeSubnet EndpointType = "subnet"
const EndpointTypeVpe EndpointType = "vpe"

type Icmp struct {
	// ICMP code allowed. If omitted, any code is allowed
	Code *int `json:"code,omitempty"`

	// Necessarily ICMP
	Protocol IcmpProtocol `json:"protocol"`

	// ICMP type allowed. If omitted, any type is allowed
	Type *int `json:"type,omitempty"`
}

type IcmpProtocol string

const IcmpProtocolICMP IcmpProtocol = "ICMP"

type Protocol interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TcpUdp) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["protocol"]; !ok || v == nil {
		return fmt.Errorf("field protocol in TcpUdp: required")
	}
	type Plain TcpUdp
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["max_destination_port"]; !ok || v == nil {
		plain.MaxDestinationPort = 65535.0
	}
	if v, ok := raw["max_source_port"]; !ok || v == nil {
		plain.MaxSourcePort = 65535.0
	}
	if v, ok := raw["min_destination_port"]; !ok || v == nil {
		plain.MinDestinationPort = 1.0
	}
	if v, ok := raw["min_source_port"]; !ok || v == nil {
		plain.MinSourcePort = 1.0
	}
	*j = TcpUdp(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Endpoint) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in Endpoint: required")
	}
	if v, ok := raw["type"]; !ok || v == nil {
		return fmt.Errorf("field type in Endpoint: required")
	}
	type Plain Endpoint
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Endpoint(plain)
	return nil
}

var enumValues_EndpointType = []interface{}{
	"external",
	"segment",
	"subnet",
	"instance",
	"nif",
	"cidr",
	"vpe",
}
var enumValues_IcmpProtocol = []interface{}{
	"ICMP",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *IcmpProtocol) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_IcmpProtocol {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_IcmpProtocol, v)
	}
	*j = IcmpProtocol(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AnyProtocol) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["protocol"]; !ok || v == nil {
		return fmt.Errorf("field protocol in AnyProtocol: required")
	}
	type Plain AnyProtocol
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AnyProtocol(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AnyProtocolProtocol) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_AnyProtocolProtocol {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_AnyProtocolProtocol, v)
	}
	*j = AnyProtocolProtocol(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Icmp) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["protocol"]; !ok || v == nil {
		return fmt.Errorf("field protocol in Icmp: required")
	}
	type Plain Icmp
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Icmp(plain)
	return nil
}

var enumValues_AnyProtocolProtocol = []interface{}{
	"ANY",
}

type ProtocolList []interface{}

type TcpUdpProtocol string

var enumValues_TcpUdpProtocol = []interface{}{
	"TCP",
	"UDP",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TcpUdpProtocol) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_TcpUdpProtocol {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_TcpUdpProtocol, v)
	}
	*j = TcpUdpProtocol(v)
	return nil
}

const TcpUdpProtocolTCP TcpUdpProtocol = "TCP"
const TcpUdpProtocolUDP TcpUdpProtocol = "UDP"

type TcpUdp struct {
	// Maximal destination port; default is 65535
	MaxDestinationPort int `json:"max_destination_port,omitempty"`

	// Maximal source port; default is 65535. Unsupported in vpc synthesis
	MaxSourcePort int `json:"max_source_port,omitempty"`

	// Minimal destination port; default is 1
	MinDestinationPort int `json:"min_destination_port,omitempty"`

	// Minimal source port; default is 1. Unsupported in vpc synthesis
	MinSourcePort int `json:"min_source_port,omitempty"`

	// Is it TCP or UDP
	Protocol TcpUdpProtocol `json:"protocol"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *EndpointType) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_EndpointType {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_EndpointType, v)
	}
	*j = EndpointType(v)
	return nil
}

// Externals are a way for users to name IP addresses or ranges external to the
// VPC. These are later used in src/dst definitions
type SpecExternals map[string]string

// Lightweight way to define instance as a list of interfaces.
type SpecInstances map[string][]string

// Lightweight way to define network interfaces.
type SpecNifs map[string]string

type SpecRequiredConnectionsElem struct {
	// List of allowed protocols
	AllowedProtocols ProtocolList `json:"allowed-protocols,omitempty"`

	// If true, allow both connections from src to dst and connections from dst to src
	Bidirectional bool `json:"bidirectional,omitempty"`

	// In unidirectional connection, this is the ingress endpoint
	Dst Endpoint `json:"dst"`

	// In unidirectional connection, this is the egress endpoint
	Src Endpoint `json:"src"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SpecRequiredConnectionsElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["dst"]; !ok || v == nil {
		return fmt.Errorf("field dst in SpecRequiredConnectionsElem: required")
	}
	if v, ok := raw["src"]; !ok || v == nil {
		return fmt.Errorf("field src in SpecRequiredConnectionsElem: required")
	}
	type Plain SpecRequiredConnectionsElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["bidirectional"]; !ok || v == nil {
		plain.Bidirectional = false
	}
	*j = SpecRequiredConnectionsElem(plain)
	return nil
}

type Type string

var enumValues_Type = []interface{}{
	"subnet",
	"instance",
	"nif",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Type) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_Type {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_Type, v)
	}
	*j = Type(v)
	return nil
}

const TypeSubnet Type = "subnet"
const TypeInstance Type = "instance"
const TypeNif Type = "nif"

// Segments are a way for users to create aggregations. These can later be used in
// src/dst fields
type SpecSegments map[string]struct {
	// All items are of the type specified in the type property, identified by name
	Items []string `json:"items"`

	// The type of the elements inside the segment
	Type Type `json:"type"`
}

// Lightweight way to define subnets.
type SpecSubnets map[string]string

type Spec struct {
	// Externals are a way for users to name IP addresses or ranges external to the
	// VPC. These are later used in src/dst definitions
	Externals SpecExternals `json:"externals,omitempty"`

	// Lightweight way to define instance as a list of interfaces.
	Instances SpecInstances `json:"instances,omitempty"`

	// Lightweight way to define network interfaces.
	Nifs SpecNifs `json:"nifs,omitempty"`

	// A list of required connections
	RequiredConnections []SpecRequiredConnectionsElem `json:"required-connections"`

	// Segments are a way for users to create aggregations. These can later be used in
	// src/dst fields
	Segments SpecSegments `json:"segments,omitempty"`

	// Lightweight way to define subnets.
	Subnets SpecSubnets `json:"subnets,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Spec) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["required-connections"]; !ok || v == nil {
		return fmt.Errorf("field required-connections in Spec: required")
	}
	type Plain Spec
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Spec(plain)
	return nil
}
