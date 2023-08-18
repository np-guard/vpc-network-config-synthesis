package toml

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
)

type Endpoint struct {
	Name string
	Type string
}

type PortRangePair struct {
	Min_dst_port *int //nolint
	Max_dst_port *int //nolint
	Min_src_port *int //nolint
	Max_src_port *int //nolint
}

type Icmp struct {
	Type *int
	Code *int
}

type Connection struct {
	Bidirectional bool
	Src           Endpoint
	Dst           Endpoint
	TCP           []PortRangePair
	UDP           []PortRangePair
	Icmp          []Icmp
}

type Define struct {
	Segments map[string][]string
	External map[string]string
}

type File struct {
	Define     Define
	Connection []Connection
}

func Test1(t *testing.T) {
	f := "../../examples/example.toml"

	var config File
	_, err := toml.DecodeFile(f, &config)
	if err != nil {
		t.Fatal(err)
	}
	res, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("%s\n", string(res))
}
