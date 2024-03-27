package ir

type IP struct{ string }

func (ip IP) String() string {
	return ip.string
}

func IPFromString(s string) IP {
	return IP{s}
}

func IPFromCidr(c CIDR) IP {
	return IPFromString(c.String())
}

func (i IP) getIP() IP {
	return i
}

const AnyIP = "0.0.0.0"

type CIDR struct{ string }

func (s CIDR) String() string {
	return s.string
}

func CidrFromString(s string) CIDR {
	return CIDR{s}
}

func CidrFromIP(ip IP) CIDR {
	return CidrFromString(ip.String())
}

const AnyCIDR = "0.0.0.0/0"
