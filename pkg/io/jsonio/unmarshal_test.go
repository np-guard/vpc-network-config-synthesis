/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package jsonio

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/np-guard/models/pkg/spec"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const (
	minPort = 1
	maxPort = 65535
)

type TestItem[T any] struct {
	input    string
	expected T
}

func TestTcpUdp_UnmarshalJSON(t *testing.T) {
	table := []TestItem[*spec.TcpUdp]{
		{`{"protocol": "TCP"}`,
			&spec.TcpUdp{Protocol: "TCP",
				MinDestinationPort: minPort, MaxDestinationPort: maxPort,
				MinSourcePort: minPort, MaxSourcePort: maxPort}},
		{`{"protocol": "UDP", "min_destination_port": 433, "max_destination_port": 433}`,
			&spec.TcpUdp{Protocol: "UDP",
				MinDestinationPort: 433, MaxDestinationPort: 433,
				MinSourcePort: minPort, MaxSourcePort: maxPort}},
	}
	for _, test := range table {
		actual := new(spec.TcpUdp)
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
	table := []TestItem[*spec.Icmp]{
		{`{"protocol": "ICMP"}`,
			&spec.Icmp{Protocol: "ICMP", Code: nil, Type: nil}},
		{`{"protocol": "ICMP", "code": 0, "type": 1}`,
			&spec.Icmp{Protocol: "ICMP", Code: utils.Ptr(0), Type: utils.Ptr(1)}},
	}
	for _, test := range table {
		actual := new(spec.Icmp)
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
	table := []TestItem[*spec.AnyProtocol]{
		{`{"protocol": "ANY"}`,
			&spec.AnyProtocol{Protocol: "ANY"}},
	}
	for _, test := range table {
		actual := new(spec.AnyProtocol)
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
	table := make([]TestItem[*spec.Resource], 7)
	for i, tp := range []string{"external", "segment", "subnet", "instance", "nif", "cidr", "vpe"} {
		name := fmt.Sprintf("ep-%v", i)
		js := fmt.Sprintf(`{"name": "%v", "type": "%v"}`, name, tp)
		resource := spec.Resource{Name: name, Type: spec.ResourceType(tp)}
		table[i] = TestItem[*spec.Resource]{js, &resource}
	}
	for _, test := range table {
		actual := new(spec.Resource)
		err := json.Unmarshal([]byte(test.input), actual)
		if err != nil {
			t.Fatalf(`Unmarshal %v returns %v`, test.input, err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Fatalf(`Unmarshal %q returns %v instead of %v`, test.input, *actual, test.expected)
		}
	}
}
