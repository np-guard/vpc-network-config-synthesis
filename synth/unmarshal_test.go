package synth

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

const examplesDir = "../examples/"

func enterField(field, ctx string, j map[string]interface{}) (string, map[string]interface{}) {
	return ctx + "." + field, j[field].(map[string]interface{})
}

func enterArray(i int, ctx string, j []interface{}) (string, map[string]interface{}) {
	return ctx + "." + fmt.Sprintf("[%v]", i), j[i].(map[string]interface{})
}

func enter[T any](field, ctx string, j map[string]interface{}) (string, T) {
	return ctx + "." + field, j[field].(T)
}

func TestUnmarshalSpec(t *testing.T) {
	ctx := ""
	filename := examplesDir + "generic_example.json"

	bytes, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf(`Read file %v returns %v`, filename, err)
	}

	jsonSpec := map[string]interface{}(nil)
	err = json.Unmarshal(bytes, &jsonSpec)
	if err != nil {
		t.Fatalf(`json.Unmarshal %v returns %v`, filename, err)
	}

	spec, err := UnmarshalSpec(filename)
	if err != nil {
		t.Fatalf(`Unmarshal %v returns %v`, filename, err)
	}
	{
		ctx, jsonExtArray := enter[[]interface{}]("externals", ctx, jsonSpec)
		if len(spec.Externals) != len(jsonExtArray) {
			t.Fatalf(`len(%v): %v != %v`, ctx, len(spec.Externals), len(jsonExtArray))
		}
		for i, ext := range spec.Externals {
			ctx, jsonExt := enterArray(i, ctx, jsonExtArray)
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
	{
		ctx, jsonConnArray := enter[[]interface{}]("required-connections", ctx, jsonSpec)
		if len(spec.RequiredConnections) != len(jsonConnArray) {
			t.Fatalf(`len(%v): %v != %v`, ctx, len(spec.RequiredConnections), len(jsonConnArray))
		}
		for i, conn := range spec.RequiredConnections {
			ctx, jsonConn := enterArray(i, ctx, jsonConnArray)
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
					ctx, jsonProtocol := enterArray(j, ctx, jsonConnAllowedProtocols)
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
}
