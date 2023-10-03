package confio

import (
	"bufio"
	"io"

	configmodel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
)

// Writer implements ir.Writer
type Writer struct {
	w     *bufio.Writer
	model *configmodel.ResourcesContainerModel
}

func NewWriter(w io.Writer, inputFilename string) (*Writer, error) {
	model, err := readModel(inputFilename)
	if err != nil {
		return nil, err
	}
	return &Writer{w: bufio.NewWriter(w), model: model}, nil
}

func (w *Writer) writeModel() error {
	s, err := w.model.ToJSONString()
	if err != nil {
		return err
	}
	_, err = w.w.WriteString(s)
	if err != nil {
		return err
	}
	err = w.w.Flush()
	return err
}
