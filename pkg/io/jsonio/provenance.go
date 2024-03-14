package jsonio

import "fmt"

type connectionOrigin struct {
	connectionIndex int
	srcName         string
	dstName         string
	inverse         bool
}

func resourceName(resource Resource) string {
	return fmt.Sprintf("(%v %v)", resource.Type, resource.Name)
}

func (o connectionOrigin) String() string {
	res := fmt.Sprintf("required-connections[%v]: %v->%v", o.connectionIndex, o.srcName, o.dstName)
	if o.inverse {
		return "inverse of " + res
	}
	return res
}

type protocolOrigin struct {
	protocolIndex int
}

func (p protocolOrigin) String() string {
	res := fmt.Sprintf("allowed-protocols[%v]", p.protocolIndex)
	return res
}
