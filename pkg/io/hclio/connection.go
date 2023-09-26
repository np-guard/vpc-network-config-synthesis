package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
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
	NPGuardConnection string          `hcl:"npguard_connection,label"`
	Name              string          `hcl:"name,label"`
	Src               cty.PathStep    `hcl:"src"`
	Dst               cty.PathStep    `hcl:"dst"`
	Bidirectional     bool            `hcl:"bidirectional,optional"`
	AllowTCP          []PortRangePair `hcl:"tcp,block"`
	AllowUDP          []PortRangePair `hcl:"udp,block"`
	AllowIcmp         []Icmp          `hcl:"icmp,block"`
}

type Module struct {
	Name  string         `hcl:"name,label"`
	Items hcl.Attributes `hcl:",remain"`
}

type File struct {
	SegmentsExternal []Module     `hcl:"module,block"`
	Connections      []Connection `hcl:"resource,block"`
}
