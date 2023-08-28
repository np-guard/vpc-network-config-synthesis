package synth

import "fmt"

type explanation struct {
	isResponse       bool
	internal         bool
	connectionOrigin fmt.Stringer
	protocolOrigin   fmt.Stringer
}

func (e explanation) response() explanation {
	e.isResponse = true
	return e
}

func (e explanation) String() string {
	locality := "External"
	if e.internal {
		locality = "Internal"
	}
	result := fmt.Sprintf("%v; %v", e.connectionOrigin, e.protocolOrigin)
	if e.isResponse {
		result = fmt.Sprintf("response to %v", result)
	}
	result = fmt.Sprintf("%v. %v", locality, result)
	return result
}
