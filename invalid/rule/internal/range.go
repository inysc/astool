package internal

import (
	_ "embed"

	"github.com/inysc/astool"
)

var (
// //go:embed template/range/single.tmpl
// rangeSingleFuncStr string
// //go:embed template/range/ptr_single.tmpl
// rangePtrSingleFuncStr string
// //go:embed template/range/slice.tmpl
// rangeSliceFuncStr string
// //go:embed template/range/ptr_slice.tmpl
// rangeSlicePtrFuncStr string

// rangeSingleTmpl    *template.Template
// rangePtrSingleTmpl *template.Template
// rangeSliceTmpl     *template.Template
// rangeSlicePtrTmpl  *template.Template
)

type RangeRule struct {
	Index      int
	Tags       string
	Message    string
	StructName string
	FieldName  string
	FieldType  string
	Left       string
	LeftVal    *string
	Right      string
	RightVal   *string
}

func (rr *RangeRule) Prio() int {
	return 2
}

func (rr *RangeRule) Check() string {
	bs := astool.NewBytes()
	tmplVal := map[string]any{
		"rule":        rr.Left,
		"tags":        rr.Tags,
		"field_name":  rr.FieldName,
		"index":       rr.Index,
		"message":     rr.Message,
		"struct_name": rr.StructName,
	}
	print(tmplVal)

	return bs.String()
}
