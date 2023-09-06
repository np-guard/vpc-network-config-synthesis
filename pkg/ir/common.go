package ir

type Direction string

const (
	Outbound Direction = "outbound"
	Inbound  Direction = "inbound"
)

type Protocol interface {
	// InverseDirection returns the response expected for a request made using this protocol
	InverseDirection() Protocol
}

type AnyProtocol struct{}

func (t AnyProtocol) InverseDirection() Protocol { return AnyProtocol{} }

type Writer interface {
	ACLWriter
	SGWriter
}
