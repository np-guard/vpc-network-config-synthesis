// Package acl describes Network ACLs
package acl

type Action string

const (
	Allow Action = "allow"
	Deny  Action = "deny"
)

type Direction string

const (
	Outbound Direction = "outbound"
	Inbound  Direction = "inbound"
)

type PortRange struct {
	Min int
	Max int
}

type PortRangePair struct {
	SrcPort PortRange
	DstPort PortRange
}

const DefaultMinPort = 1
const DefaultMaxPort = 65535

func Swap(pair PortRangePair) PortRangePair {
	return PortRangePair{SrcPort: pair.DstPort, DstPort: pair.SrcPort}
}

type TCP struct {
	PortRangePair
}

type UDP struct {
	PortRangePair
}

type ICMP struct {
	Code *int
	Type *int
}

type AnyProtocol struct{}

type Protocol interface {
	SwapSrcDstPortRange() Protocol
	Name() string
}

func (t TCP) SwapSrcDstPortRange() Protocol { return TCP{Swap(t.PortRangePair)} }

func (t UDP) SwapSrcDstPortRange() Protocol { return UDP{Swap(t.PortRangePair)} }

func (t ICMP) SwapSrcDstPortRange() Protocol { return ICMP{Code: t.Code, Type: t.Type} }

func (t AnyProtocol) SwapSrcDstPortRange() Protocol { return AnyProtocol{} }

func (t TCP) Name() string { return "TCP" }

func (t UDP) Name() string { return "UDP" }

func (t ICMP) Name() string { return "ICMP" }

func (t AnyProtocol) Name() string { return "All" }

type Rule struct {
	Name        string
	Action      Action
	Direction   Direction
	Source      string
	Destination string
	Protocol    Protocol
}

type ACL struct {
	Name          string
	ResourceGroup string
	Vpc           string
	Rules         []*Rule
}

type Collection struct {
	Items []*ACL
}

type Writer interface {
	Write(Collection) error
}
