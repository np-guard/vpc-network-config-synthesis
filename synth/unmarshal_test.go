package synth

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

const examplesDir = "../examples/"

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
		ctx := ctx + ".required-connections"
		jsonConnArray := jsonSpec["required-connections"].([]interface{})
		if len(spec.RequiredConnections) != len(jsonConnArray) {
			t.Fatalf(`len(required-connections): %v != %v`, len(spec.RequiredConnections), len(jsonConnArray))
		}
		for i, conn := range spec.RequiredConnections {
			ctx := ctx + fmt.Sprintf("[%v]", i)
			jsonConn := jsonConnArray[i].(map[string]interface{})
			{
				ctx := ctx + ".src"
				jsonConnSrc := jsonConn["src"].(map[string]interface{})
				{
					ctx := ctx + ".type"
					jsonConnSrcType := jsonConnSrc["type"].(string)
					if string(conn.Src.Type) != jsonConnSrcType {
						t.Fatalf(`%v: %v != %v`, ctx, conn.Src.Type, jsonConnSrcType)
					}
				}
				{
					ctx := ctx + ".name"
					jsonConnSrcName := jsonConnSrc["name"].(string)
					if conn.Src.Name != jsonConnSrcName {
						t.Fatalf(`%v: %v != %v`, ctx, conn.Src.Name, jsonConnSrcName)
					}
				}
			}
			{
				ctx := ctx + ".dst"
				jsonConnDst := jsonConn["dst"].(map[string]interface{})
				{
					ctx := ctx + ".type"
					jsonConnDstType := jsonConnDst["type"].(string)
					if string(conn.Dst.Type) != jsonConnDstType {
						t.Fatalf(`%v: %v != %v`, ctx, conn.Dst.Type, jsonConnDstType)
					}
				}
				{
					ctx := ctx + ".name"
					jsonConnDstName := jsonConnDst["name"].(string)
					if conn.Dst.Name != jsonConnDstName {
						t.Fatalf(`%v: %v != %v`, ctx, conn.Dst.Name, jsonConnDstName)
					}
				}
			}

			{
				ctx := ctx + ".allowed-protocols"
				jsonConnDstAllowedProtocols := jsonConn["allowed-protocols"].([]interface{})
				if len(conn.AllowedProtocols) != len(jsonConnDstAllowedProtocols) {
					t.Fatalf(`len(%v): %v != %v`, ctx, len(conn.AllowedProtocols), len(jsonConnDstAllowedProtocols))
				}
				for j, protocol := range conn.AllowedProtocols {
					ctx := ctx + fmt.Sprintf("[%v]", j)
					jsonProtocol := jsonConnDstAllowedProtocols[j].(map[string]interface{})
					jsonProtocolName := jsonProtocol["protocol"]
					switch p := protocol.(type) {
					case TcpUdp:
						{
							ctx := ctx + ".protocol"
							if string(p.Protocol) != jsonProtocolName {
								t.Fatalf(`%v: %v != %v`, ctx, p.Protocol, jsonProtocolName)
							}
						}
						{
							ctx := ctx + ".min_port"
							jsonPortMinPort := jsonProtocol["min_port"].(int)
							if p.MinPort != jsonPortMinPort {
								t.Fatalf(`%v: %v != %v`, ctx, p.MinPort, jsonPortMinPort)
							}
						}
						{
							ctx := ctx + ".max_port"
							jsonPortMaxPort := jsonProtocol["max_port"].(int)
							if p.MaxPort != jsonPortMaxPort {
								t.Fatalf(`%v: %v != %v`, ctx, p.MaxPort, jsonPortMaxPort)
							}
						}
						{
							ctx := ctx + ".bidirectional"
							jsonBidirectional := jsonProtocol["bidirectional"].(bool)
							if p.Bidirectional != jsonBidirectional {
								t.Fatalf(`%v: %t != %t`, ctx, p.Bidirectional, jsonBidirectional)
							}
						}
					case Icmp:
						{
							ctx := ctx + ".protocol"
							if p.Protocol.(string) != jsonProtocolName {
								t.Fatalf(`%v: %v != %v`, ctx, p.Protocol, jsonProtocolName)
							}
						}
						{
							ctx := ctx + ".type"
							jsonPortType := jsonProtocol["type"].(*int)
							if p.Type != jsonPortType {
								t.Fatalf(`%v: %v != %v`, ctx, p.Type, jsonPortType)
							}
						}
						{
							ctx := ctx + ".code"
							jsonPortCode := jsonProtocol["code"].(*int)
							if p.Code != jsonPortCode {
								t.Fatalf(`%v: %v != %v`, ctx, p.Code, jsonPortCode)
							}
						}
						{
							ctx := ctx + ".bidirectional"
							jsonBidirectional := jsonProtocol["bidirectional"].(bool)
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
