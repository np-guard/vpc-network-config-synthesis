package ir

type IP struct{ string }

func (ip IP) String() string {
	return ip.string
}

func IPFromString(s string) IP {
	return IP{s}
}

const AnyIP = "0.0.0.0"

type CIDR struct{ string }

func (s CIDR) String() string {
	return s.string
}

func CidrFromString(s string) CIDR {
	return CIDR{s}
}

const AnyCIDR = "0.0.0.0/0"
