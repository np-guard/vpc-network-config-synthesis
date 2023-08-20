// Package spec describes the input-format-agnostic specification of the required connectivity
package spec

type (
	Spec struct {
		// Subnet definitions
		Subnets map[string]string

		// Segments are a way for users to create aggregations. These can later be used in
		// src/dst fields
		SubnetSegments map[string][]string

		// Externals are a way for users to name IP addresses or ranges external to the
		// VPC. These are later used in src/dst definitions
		Externals map[string]string

		// A list of required connections
		Connections []Connection
	}

	Connection struct {
		// If true, allow both connections from src to dst and connections from dst to src
		Bidirectional bool

		// In unidirectional connection, this is the ingress endpoint
		Dst *Endpoint

		// In unidirectional connection, this is the egress endpoint
		Src *Endpoint

		// List of allowed transport-layer connections
		Protocols []Protocol
	}

	Endpoint struct {
		// Name of endpoint
		Name string

		// Type of endpoint
		Type EndpointType
	}
)

type EndpointType string

const (
	EndpointTypeExternal EndpointType = "external"
	EndpointTypeSegment  EndpointType = "segment"
	EndpointTypeSubnet   EndpointType = "subnet"
)

type Reader interface {
	ReadSpec(filename string, subnetMap map[string]string) (*Spec, error)
}
