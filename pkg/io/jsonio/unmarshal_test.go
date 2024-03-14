package jsonio

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const examplesDir = "../../../examples/"
const (
	minPort = 1
	maxPort = 65535
)

type TestItem[T any] struct {
	input    string
	expected T
}

func readFile(t *testing.T, filename string) map[string]interface{} {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf(`Read file %v returns %v`, filename, err)
	}

	jsonSpec := map[string]interface{}(nil)
	err = json.Unmarshal(bytes, &jsonSpec)
	if err != nil {
		t.Fatalf(`input.Unmarshal %v returns %v`, filename, err)
	}
	return jsonSpec
}

func TestTcpUdp_UnmarshalJSON(t *testing.T) {
	table := []TestItem[*TcpUdp]{
		{`{"protocol": "TCP"}`,
			&TcpUdp{Protocol: "TCP",
				MinDestinationPort: minPort, MaxDestinationPort: maxPort,
				MinSourcePort: minPort, MaxSourcePort: maxPort}},
		{`{"protocol": "UDP", "min_destination_port": 433, "max_destination_port": 433}`,
			&TcpUdp{Protocol: "UDP",
				MinDestinationPort: 433, MaxDestinationPort: 433,
				MinSourcePort: minPort, MaxSourcePort: maxPort}},
	}
	for _, test := range table {
		actual := new(TcpUdp)
		err := json.Unmarshal([]byte(test.input), actual)
		if err != nil {
			t.Fatalf(`Unmarshal %q returns %v`, test.input, err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Fatalf(`Unmarshal %q returns %v instead of %v`, test.input, *actual, *test.expected)
		}
	}
}

func TestIcmp_UnmarshalJSON(t *testing.T) {
	table := []TestItem[*Icmp]{
		{`{"protocol": "ICMP"}`,
			&Icmp{Protocol: "ICMP", Code: nil, Type: nil}},
		{`{"protocol": "ICMP", "code": 0, "type": 1}`,
			&Icmp{Protocol: "ICMP", Code: utils.Ptr(0), Type: utils.Ptr(1)}},
	}
	for _, test := range table {
		actual := new(Icmp)
		err := json.Unmarshal([]byte(test.input), actual)
		if err != nil {
			t.Fatalf(`Unmarshal %v returns %v`, test.input, err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Fatalf(`Unmarshal %q returns %v instead of %v`, test.input, *actual, *test.expected)
		}
	}
}

func TestAnyProtocol_UnmarshalJSON(t *testing.T) {
	table := []TestItem[*AnyProtocol]{
		{`{"protocol": "ANY"}`,
			&AnyProtocol{Protocol: "ANY"}},
	}
	for _, test := range table {
		actual := new(AnyProtocol)
		err := json.Unmarshal([]byte(test.input), actual)
		if err != nil {
			t.Fatalf(`Unmarshal %v returns %v`, test.input, err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Fatalf(`Unmarshal %q returns %v instead of %v`, test.input, *actual, *test.expected)
		}
	}
}

func TestResource_UnmarshalJSON(t *testing.T) {
	table := make([]TestItem[*Resource], 7)
	for i, tp := range []string{"external", "segment", "subnet", "instance", "nif", "cidr", "vpe"} {
		name := fmt.Sprintf("ep-%v", i)
		js := fmt.Sprintf(`{"name": "%v", "type": "%v"}`, name, tp)
		resource := Resource{Name: name, Type: ResourceType(tp)}
		table[i] = TestItem[*Resource]{js, &resource}
	}
	for _, test := range table {
		actual := new(Resource)
		err := json.Unmarshal([]byte(test.input), actual)
		if err != nil {
			t.Fatalf(`Unmarshal %v returns %v`, test.input, err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Fatalf(`Unmarshal %q returns %v instead of %v`, test.input, *actual, test.expected)
		}
	}
}

// Compare unmarshalled structs/arrays for "segments" in a spec file against simple json maps
//
//goland:noinspection GoShadowedVar
//nolint:govet // ctx is intentionally shadowed, allowing stack-like navigation
func TestUnmarshalSpecSegments(t *testing.T) {
	ctx := ""
	filename := examplesDir + "generic_example.json"

	jsonSpec := readFile(t, filename)

	spec, err := unmarshal(filename)
	if err != nil {
		t.Fatalf(`Unmarshal %v returns %v`, filename, err)
	}
	ctx, jsonSegmentMap := enter[map[string]interface{}]("segments", ctx, jsonSpec)
	if len(spec.Segments) != len(jsonSegmentMap) {
		t.Fatalf(`len(%v): %v != %v`, ctx, len(spec.Segments), len(jsonSegmentMap))
	}
	for field, segment := range spec.Segments {
		ctx, jsonSegment := enter[map[string]interface{}](field, ctx, jsonSegmentMap)
		{
			ctx, jsonType := enter[string]("type", ctx, jsonSegment)
			if string(segment.Type) != jsonType {
				t.Fatalf(`%v: %v != %v`, ctx, segment.Type, jsonType)
			}
		}
		{
			ctx, jsonSegmentItemsArray := enter[[]interface{}]("items", ctx, jsonSegment)
			if len(segment.Items) != len(jsonSegmentItemsArray) {
				t.Fatalf(`len(%v): %v != %v`, ctx, len(segment.Items), len(jsonSegmentItemsArray))
			}
			for j, item := range segment.Items {
				ctx, jsonItem := enterArray[string](j, ctx, jsonSegmentItemsArray)
				if item != jsonItem {
					t.Fatalf(`%v: %v != %v`, ctx, item, jsonItem)
				}
			}
		}
	}
}

// Compare unmarshalled structs/arrays for "externals" in a spec file against simple json maps
//
//goland:noinspection GoShadowedVar
//nolint:govet // ctx is intentionally shadowed, allowing stack-like navigation
func TestUnmarshalSpecExternals(t *testing.T) {
	ctx := ""
	filename := examplesDir + "generic_example.json"

	jsonSpec := readFile(t, filename)

	spec, err := unmarshal(filename)
	if err != nil {
		t.Fatalf(`Unmarshal %v returns %v`, filename, err)
	}
	ctx, jsonExtMap := enter[map[string]interface{}]("externals", ctx, jsonSpec)
	if len(spec.Externals) != len(jsonExtMap) {
		t.Fatalf(`len(%v): %v != %v`, ctx, len(spec.Externals), len(jsonExtMap))
	}
	for field, value := range spec.Externals {
		ctx, jsonExt := enter[string](field, ctx, jsonExtMap)
		if value != jsonExt {
			t.Fatalf(`%v: %v != %v`, ctx, value, jsonExt)
		}
	}
}

// Compare unmarshalled structs/arrays for "required-connections" in a spec file against simple json maps
//
//goland:noinspection GoShadowedVar
//nolint:gocyclo,govet // ctx is intentionally shadowed, allowing stack-like navigation
func TestUnmarshalSpecRequiredConnections(t *testing.T) {
	ctx := ""
	filename := examplesDir + "generic_example.json"

	jsonSpec := readFile(t, filename)

	spec, err := unmarshal(filename)
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
			ctx, jsonConnResource := enterField("src", ctx, jsonConn)
			resource := conn.Src
			{
				ctx, jsonConnResourceType := enter[string]("type", ctx, jsonConnResource)
				if string(resource.Type) != jsonConnResourceType {
					t.Fatalf(`%v: %v != %v`, ctx, resource.Type, jsonConnResourceType)
				}
			}
			{
				ctx, jsonConnResourceName := enter[string]("name", ctx, jsonConnResource)
				if resource.Name != jsonConnResourceName {
					t.Fatalf(`%v: %v != %v`, ctx, resource.Name, jsonConnResourceName)
				}
			}
		}
		{
			ctx, jsonConnResource := enterField("dst", ctx, jsonConn)
			resource := conn.Dst
			{
				ctx, jsonConnResourceType := enter[string]("type", ctx, jsonConnResource)
				if string(resource.Type) != jsonConnResourceType {
					t.Fatalf(`%v: %v != %v`, ctx, resource.Type, jsonConnResourceType)
				}
			}
			{
				ctx, jsonConnResourceName := enter[string]("name", ctx, jsonConnResource)
				if resource.Name != jsonConnResourceName {
					t.Fatalf(`%v: %v != %v`, ctx, resource.Name, jsonConnResourceName)
				}
			}
		}
		{
			ctx, jsonBidirectional := enter[bool]("bidirectional", ctx, jsonConn)
			if conn.Bidirectional != jsonBidirectional {
				t.Fatalf(`%v: %t != %t`, ctx, conn.Bidirectional, jsonBidirectional)
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
						ctx, jsonProtocolPort := enterInt("min_destination_port", ctx, jsonProtocol)
						if p.MinDestinationPort != jsonProtocolPort {
							t.Fatalf(`%v: %v != %v`, ctx, p.MinDestinationPort, jsonProtocolPort)
						}
					}
					{
						ctx, jsonProtocolPort := enterInt("max_destination_port", ctx, jsonProtocol)
						if p.MaxDestinationPort != jsonProtocolPort {
							t.Fatalf(`%v: %v != %v`, ctx, p.MaxDestinationPort, jsonProtocolPort)
						}
					}
				case Icmp:
					{
						ctx, jsonProtocolName := enter[string]("protocol", ctx, jsonProtocol)
						if string(p.Protocol) != jsonProtocolName {
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
				case AnyProtocol:
					t.Fatalf("Unsupported")
				default:
					t.Fatalf("Bad protocol %v", p)
				}
			}
		}
	}
}
