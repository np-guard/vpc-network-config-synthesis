package ir

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

type Protocol interface {
	// InverseDirection returns the response expected for a request made using this protocol
	InverseDirection() Protocol
}

type AnyProtocol struct{}

func (t AnyProtocol) InverseDirection() Protocol { return AnyProtocol{} }

type Rule struct {
	Action      Action
	Direction   Direction
	Source      IP
	Destination IP
	Protocol    Protocol
	Explanation string
}

type ACL struct {
	Subnet   string
	Internal []Rule
	External []Rule
}

type Collection struct {
	ACLs map[string]*ACL
}

type Writer interface {
	Write(*Collection) error
}
