package synth

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

const examplesDir = "../../examples/"

func enterField(field, ctx string, j map[string]interface{}) (resCtx string, res map[string]interface{}) {
	resCtx = ctx + "." + field
	res, ok := j[field].(map[string]interface{})
	if !ok {
		res = nil
	}
	return
}

func enterArray[T any](i int, ctx string, j []interface{}) (resCtx string, res T) {
	resCtx = ctx + "." + fmt.Sprintf("[%v]", i)
	res, ok := j[i].(T)
	if !ok {
		res = *new(T)
	}
	return
}

func enter[T any](field, ctx string, j map[string]interface{}) (resCtx string, res T) {
	resCtx = ctx + "." + field
	res, ok := j[field].(T)
	if !ok {
		// TODO: use default argument
		res = *new(T)
	}
	return
}

func readFile(t *testing.T, filename string) map[string]interface{} {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf(`Read file %v returns %v`, filename, err)
	}

	jsonSpec := map[string]interface{}(nil)
	err = json.Unmarshal(bytes, &jsonSpec)
	if err != nil {
		t.Fatalf(`json.Unmarshal %v returns %v`, filename, err)
	}
	return jsonSpec
}

func TestUnmarshalSpecSections(t *testing.T) {
	ctx := ""
	filename := examplesDir + "generic_example.json"

	jsonSpec := readFile(t, filename)

	spec, err := UnmarshalSpec(filename)
	if err != nil {
		t.Fatalf(`Unmarshal %v returns %v`, filename, err)
	}
	ctx, jsonSectionArray := enter[[]interface{}]("sections", ctx, jsonSpec)
	if len(spec.Sections) != len(jsonSectionArray) {
		t.Fatalf(`len(%v): %v != %v`, ctx, len(spec.Externals), len(jsonSectionArray))
	}
	for i, section := range spec.Sections {
		ctx, jsonSection := enterArray[map[string]interface{}](i, ctx, jsonSectionArray)
		{
			ctx, jsonName := enter[string]("name", ctx, jsonSection)
			if section.Name != jsonName {
				t.Fatalf(`%v: %v != %v`, ctx, section.Name, jsonName)
			}
		}
		{
			ctx, jsonType := enter[string]("type", ctx, jsonSection)
			if string(section.Type) != jsonType {
				t.Fatalf(`%v: %v != %v`, ctx, section.Name, jsonType)
			}
		}
		{
			ctx, jsonSectionItemsArray := enter[[]interface{}]("items", ctx, jsonSection)
			if len(section.Items) != len(jsonSectionItemsArray) {
				t.Fatalf(`len(%v): %v != %v`, ctx, len(section.Items), len(jsonSectionItemsArray))
			}
			for j, item := range section.Items {
				ctx, jsonItem := enterArray[string](j, ctx, jsonSectionItemsArray)
				if item != jsonItem {
					t.Fatalf(`%v: %v != %v`, ctx, item, jsonItem)
				}
			}
		}
		{
			ctx, jsonFullyConnected := enter[bool]("fully-connected", ctx, jsonSection)
			if section.FullyConnected != jsonFullyConnected {
				t.Fatalf(`%v: %t != %t`, ctx, section.FullyConnected, jsonFullyConnected)
			}
		}
		// TODO: Check fully-connected-with-connection-type.
		// It is the same code as allowed-protocols, but refactoring into assertion functions is not idiomatic in Go
	}
}

func TestUnmarshalSpecExternals(t *testing.T) {
	ctx := ""
	filename := examplesDir + "generic_example.json"

	jsonSpec := readFile(t, filename)

	spec, err := UnmarshalSpec(filename)
	if err != nil {
		t.Fatalf(`Unmarshal %v returns %v`, filename, err)
	}
	ctx, jsonExtArray := enter[[]interface{}]("externals", ctx, jsonSpec)
	if len(spec.Externals) != len(jsonExtArray) {
		t.Fatalf(`len(%v): %v != %v`, ctx, len(spec.Externals), len(jsonExtArray))
	}
	for i, ext := range spec.Externals {
		ctx, jsonExt := enterArray[map[string]interface{}](i, ctx, jsonExtArray)
		{
			ctx, jsonName := enter[string]("name", ctx, jsonExt)
			if ext.Name != jsonName {
				t.Fatalf(`%v: %v != %v`, ctx, ext.Name, jsonName)
			}
		}
		{
			ctx, jsonCidr := enter[string]("cidr", ctx, jsonExt)
			if ext.Cidr != jsonCidr {
				t.Fatalf(`%v: %v != %v`, ctx, ext.Name, jsonCidr)
			}
		}
	}
}

func TestUnmarshalSpecRequiredConnections(t *testing.T) {
	ctx := ""
	filename := examplesDir + "generic_example.json"

	jsonSpec := readFile(t, filename)

	spec, err := UnmarshalSpec(filename)
	if err != nil {
		t.Fatalf(`Unmarshal %v returns %v`, filename, err)
	}
	ctx, jsonConnArray := enter[[]interface{}]("required-connections", ctx, jsonSpec)
	if len(spec.RequiredConnections) != len(jsonConnArray) {
		t.Fatalf(`len(%v): %v != %v`, ctx, len(spec.RequiredConnections), len(jsonConnArray))
	}
	for i, conn := range spec.RequiredConnections {
		ctx, jsonConn := enterArray[map[string]interface{}](i, ctx, jsonConnArray)
		{
			ctx, jsonConnEndpoint := enterField("src", ctx, jsonConn)
			endpoint := conn.Src
			{
				ctx, jsonConnEndpointType := enter[string]("type", ctx, jsonConnEndpoint)
				if string(endpoint.Type) != jsonConnEndpointType {
					t.Fatalf(`%v: %v != %v`, ctx, endpoint.Type, jsonConnEndpointType)
				}
			}
			{
				ctx, jsonConnEndpointName := enter[string]("name", ctx, jsonConnEndpoint)
				if endpoint.Name != jsonConnEndpointName {
					t.Fatalf(`%v: %v != %v`, ctx, endpoint.Name, jsonConnEndpointName)
				}
			}
		}
		{
			ctx, jsonConnEndpoint := enterField("dst", ctx, jsonConn)
			endpoint := conn.Dst
			{
				ctx, jsonConnEndpointType := enter[string]("type", ctx, jsonConnEndpoint)
				if string(endpoint.Type) != jsonConnEndpointType {
					t.Fatalf(`%v: %v != %v`, ctx, endpoint.Type, jsonConnEndpointType)
				}
			}
			{
				ctx, jsonConnEndpointName := enter[string]("name", ctx, jsonConnEndpoint)
				if endpoint.Name != jsonConnEndpointName {
					t.Fatalf(`%v: %v != %v`, ctx, endpoint.Name, jsonConnEndpointName)
				}
			}
		}

		{
			ctx, jsonConnAllowedProtocols := enter[[]interface{}]("allowed-protocols", ctx, jsonConn)
			if len(conn.AllowedProtocols) != len(jsonConnAllowedProtocols) {
				t.Fatalf(`len(%v): %v != %v`, ctx, len(conn.AllowedProtocols), len(jsonConnAllowedProtocols))
			}
			for j, protocol := range conn.AllowedProtocols {
				ctx, jsonProtocol := enterArray[map[string]interface{}](j, ctx, jsonConnAllowedProtocols)
				switch p := protocol.(type) {
				case TcpUdp:
					{
						ctx, jsonProtocolName := enter[string]("protocol", ctx, jsonProtocol)
						if string(p.Protocol) != jsonProtocolName {
							t.Fatalf(`%v: %v != %v`, ctx, p.Protocol, jsonProtocolName)
						}
					}
					{
						ctx, jsonProtocolPort := enter[int]("min_port", ctx, jsonProtocol)
						if p.MinPort != jsonProtocolPort {
							t.Fatalf(`%v: %v != %v`, ctx, p.MinPort, jsonProtocolPort)
						}
					}
					{
						ctx, jsonProtocolPort := enter[int]("max_port", ctx, jsonProtocol)
						if p.MaxPort != jsonProtocolPort {
							t.Fatalf(`%v: %v != %v`, ctx, p.MaxPort, jsonProtocolPort)
						}
					}
					{
						ctx, jsonBidirectional := enter[bool]("bidirectional", ctx, jsonProtocol)
						if p.Bidirectional != jsonBidirectional {
							t.Fatalf(`%v: %t != %t`, ctx, p.Bidirectional, jsonBidirectional)
						}
					}
				case Icmp:
					{
						ctx, jsonProtocolName := enter[string]("protocol", ctx, jsonProtocol)
						if p.Protocol.(string) != jsonProtocolName {
							t.Fatalf(`%v: %v != %v`, ctx, p.Protocol, jsonProtocolName)
						}
					}
					{
						ctx, jsonProtocolType := enter[*int]("type", ctx, jsonProtocol)
						if p.Type != jsonProtocolType {
							t.Fatalf(`%v: %v != %v`, ctx, p.Type, jsonProtocolType)
						}
					}
					{
						ctx, jsonProtocolCode := enter[*int]("code", ctx, jsonProtocol)
						if p.Code != jsonProtocolCode {
							t.Fatalf(`%v: %v != %v`, ctx, p.Code, jsonProtocolCode)
						}
					}
					{
						ctx, jsonBidirectional := enter[bool]("bidirectional", ctx, jsonProtocol)
						if p.Bidirectional != jsonBidirectional {
							t.Fatalf(`%v: %t != %t`, ctx, p.Bidirectional, jsonBidirectional)
						}
					}
				}
			}
		}
	}
}
