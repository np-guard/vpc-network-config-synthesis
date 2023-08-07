// Package synth takes global specification written in a JSON file, as described by spec_schema.input,
// and generates NetworkACLs that collectively enable the connectivity described in the global specification.
package synth

import (
	"encoding/json"
	"os"
)

// UnmarshalSpec returns a Spec struct given a file adhering to spec_schema.input
func UnmarshalSpec(filename string) (*Spec, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	spec := new(Spec)
	err = json.Unmarshal(bytes, spec)
	if err != nil {
		return nil, err
	}
	return spec, err
}
