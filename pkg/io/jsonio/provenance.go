package jsonio

import "fmt"

type Origin struct {
	connectionIndex int
	srcName         string
	dstName         string
	inverse         bool
}

func endpointName(endpoint Endpoint) string {
	return fmt.Sprintf("(%v %v)", endpoint.Type, endpoint.Name)
}

func (o Origin) String() string {
	res := fmt.Sprintf("Connection #%v: %v->%v", o.connectionIndex, o.srcName, o.dstName)
	if o.inverse {
		return "inverse of " + res
	}
	return res
}
