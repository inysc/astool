package internal

import (
	_ "embed"
	"log"
	"strings"
	"text/template"

	"github.com/inysc/astool"
)

var (
	//go:embed template/equal/single.tmpl
	equalSingleFuncStr string
	//go:embed template/equal/ptr_single.tmpl
	equalPtrSingleFuncStr string
	//go:embed template/equal/slice.tmpl
	equalSliceFuncStr string
	//go:embed template/equal/ptr_slice.tmpl
	equalSlicePtrFuncStr string

	equalSingleTmpl    *template.Template
	equalPtrSingleTmpl *template.Template
	equalSliceTmpl     *template.Template
	equalSlicePtrTmpl  *template.Template
)

func init() {
	var err error
	{
		equalSingleTmpl, err = template.
			New("invalid_equal").
			Funcs(funcMap).
			Parse(equalSingleFuncStr)
		if err != nil {
			panic(err)
		}

		equalPtrSingleTmpl, err = template.
			New("invalid_equal_ptr").
			Funcs(funcMap).
			Parse(equalPtrSingleFuncStr)
		if err != nil {
			panic(err)
		}

		equalSliceTmpl, err = template.
			New("invalid_equal_slice").
			Funcs(funcMap).
			Parse(equalSliceFuncStr)
		if err != nil {
			panic(err)
		}

		equalSlicePtrTmpl, err = template.
			New("invalid_equal_slice_ptr").
			Funcs(funcMap).
			Parse(equalSlicePtrFuncStr)
		if err != nil {
			panic(err)
		}
	}
}

type equalRule struct {
	Rule       string
	Tags       string
	FieldName  string
	FieldType  string
	Index      int
	Message    string
	StructName string
}

func (nr *equalRule) Prio() int {
	return 2
}

func (nr *equalRule) Check() string {
	bs := astool.NewBytes()
	tmplVal := map[string]any{
		"rule":        nr.Rule,
		"tags":        nr.Tags,
		"field_name":  nr.FieldName,
		"index":       nr.Index,
		"message":     nr.Message,
		"struct_name": nr.StructName,
	}
	var err error
	if strings.HasPrefix(nr.FieldType, "*[]") {
		err = equalSlicePtrTmpl.Execute(bs, tmplVal)
	} else if strings.HasPrefix(nr.FieldType, "[]") {
		err = equalSliceTmpl.Execute(bs, tmplVal)
	} else if strings.HasPrefix(nr.FieldType, "*") {
		err = equalPtrSingleTmpl.Execute(bs, tmplVal)
	} else {
		err = equalSingleTmpl.Execute(bs, tmplVal)
	}
	if err != nil {
		log.Fatal(err)
	}
	return bs.String()
}

func NewEqualRule(structName, fieldName, fieldType string, rule *Rule) *equalRule {
	INDEX++

	return &equalRule{
		Rule:       rule.Rule(),
		Tags:       rule.Tags(),
		FieldName:  fieldName,
		FieldType:  fieldType,
		Index:      INDEX,
		Message:    rule.Message(),
		StructName: structName,
	}
}
