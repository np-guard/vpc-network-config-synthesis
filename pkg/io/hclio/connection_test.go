package hcl

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const examplesDir = "../../../examples/"

func Test1(t *testing.T) {
	var ctx hcl.EvalContext
	var f File
	err := hclsimple.DecodeFile(examplesDir+"golden_eye.hcl", &ctx, &f)
	if err != nil {
		t.Fatal(err)
	}
	hf := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&f, hf.Body())
	fmt.Printf("%s", hf.Bytes())
	t.Log(f)
}
