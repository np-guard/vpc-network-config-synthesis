package synth

import (
	"os"
)

func UnmarshalSpec(filename string) (*Spec, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	spec := new(Spec)
	err = spec.UnmarshalJSON(bytes)
	if err != nil {
		return nil, err
	}
	return spec, err
}
