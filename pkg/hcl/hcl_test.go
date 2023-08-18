package hcl

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Endpoint struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type,label"`
}

type PortRangePair struct {
	MinDestinationPort *int `hcl:"min_port,optional"`
	MaxDestinationPort *int `hcl:"max_port,optional"`
	MinSourcePort      *int `hcl:"min_src_port,optional"`
	MaxSourcePort      *int `hcl:"max_src_port,optional"`
}

type Icmp struct {
	Code *int `hcl:"code,optional"`
	Type *int `hcl:"type,optional"`
}

type Connection struct {
	Bidirectional bool            `hcl:"bidirectional,optional"`
	Src           Endpoint        `hcl:"src,block"`
	Dst           Endpoint        `hcl:"dst,block"`
	AllowTCP      []PortRangePair `hcl:"tcp,block"`
	AllowUDP      []PortRangePair `hcl:"udp,block"`
	AllowIcmp     []Icmp          `hcl:"icmp,block"`
}

type Definition struct {
	Items hcl.Attributes `hcl:",remain"`
}

type Define struct {
	Segments []Definition `hcl:"segment,block"`
	External []Definition `hcl:"external,block"`
}

type File struct {
	Define      Define       `hcl:"define,block"`
	Connections []Connection `hcl:"connection,block"`
}

func Test1(t *testing.T) {
	var ctx hcl.EvalContext
	var f File
	err := hclsimple.DecodeFile("../../examples/example.hcl", &ctx, &f)
	if err != nil {
		t.Fatal(err)
	}
	hf := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&f, hf.Body())
	fmt.Printf("%s", hf.Bytes())
	t.Log(f)
}
