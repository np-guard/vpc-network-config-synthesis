package synth

import "fmt"

type explanation struct {
	isResponse       bool
	internal         bool
	connectionOrigin fmt.Stringer
	protocolOrigin   fmt.Stringer
}

func (e explanation) String() string {
	locality := "External"
	if e.internal {
		locality = "Internal"
	}
	result := fmt.Sprintf("%v. %v; %v", locality, e.connectionOrigin, e.protocolOrigin)
	if e.isResponse {
		result = "response to " + result
	}
	return result
}
