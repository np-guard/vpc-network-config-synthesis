package ir

type IP struct{ string }

func (ip IP) String() string { return ip.string }

func IPFromString(s string) IP {
	return IP{s}
}
