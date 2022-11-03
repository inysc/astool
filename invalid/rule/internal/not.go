package internal

import (
	_ "embed"
	"log"
	"strings"
	"text/template"

	"github.com/inysc/astool"
)

var (
	//go:embed template/not/single.tmpl
	notSingleFuncStr string
	//go:embed template/not/ptr_single.tmpl
	notPtrSingleFuncStr string
	//go:embed template/not/slice.tmpl
	notSliceFuncStr string
	//go:embed template/not/ptr_slice.tmpl
	notSlicePtrFuncStr string

	notSingleTmpl    *template.Template
	notPtrSingleTmpl *template.Template
	notSliceTmpl     *template.Template
	notSlicePtrTmpl  *template.Template
)

func init() {
	var err error
	{
		notSingleTmpl, err = template.
			New("invalid_not").
			Funcs(funcMap).
			Parse(notSingleFuncStr)
		if err != nil {
			panic(err)
		}

		notPtrSingleTmpl, err = template.
			New("invalid_not_ptr").
			Funcs(funcMap).
			Parse(notPtrSingleFuncStr)
		if err != nil {
			panic(err)
		}

		notSliceTmpl, err = template.
			New("invalid_not_slice").
			Funcs(funcMap).
			Parse(notSliceFuncStr)
		if err != nil {
			panic(err)
		}

		notSlicePtrTmpl, err = template.
			New("invalid_not_slice_ptr").
			Funcs(funcMap).
			Parse(notSlicePtrFuncStr)
		if err != nil {
			panic(err)
		}
	}
}

type notRule struct {
	Rule       string
	Tags       string
	FieldName  string
	FieldType  string
	Index      int
	Message    string
	StructName string
}

func (nr *notRule) Prio() int {
	return 2
}

func (nr *notRule) Check() string {
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
		err = notSlicePtrTmpl.Execute(bs, tmplVal)
	} else if strings.HasPrefix(nr.FieldType, "[]") {
		err = notSliceTmpl.Execute(bs, tmplVal)
	} else if strings.HasPrefix(nr.FieldType, "*") {
		err = notPtrSingleTmpl.Execute(bs, tmplVal)
	} else {
		err = notSingleTmpl.Execute(bs, tmplVal)
	}
	if err != nil {
		log.Fatal(err)
	}
	return bs.String()
}

func NewNotRule(structName, fieldName, fieldType string, rule *Rule) *notRule {
	INDEX++

	return &notRule{
		Rule:       rule.Rule(),
		Tags:       rule.Tags(),
		FieldName:  fieldName,
		FieldType:  fieldType,
		Index:      INDEX,
		Message:    rule.Message(),
		StructName: structName,
	}
}
