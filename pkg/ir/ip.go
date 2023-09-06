package ir

type IP struct{ string }

func (ip IP) String() string { return ip.string }

func IPFromString(s string) IP {
	return IP{s}
}

type CIDR string

func (s CIDR) String() string {
	return string(s)
}
