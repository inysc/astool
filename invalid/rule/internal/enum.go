package internal

import (
	_ "embed"
	"log"
	"strings"
	"text/template"

	"github.com/inysc/astool"
)

var (
	//go:embed template/enum/single.tmpl
	enumSingleFuncStr string
	//go:embed template/enum/ptr_single.tmpl
	enumPtrSingleFuncStr string
	//go:embed template/enum/slice.tmpl
	enumSliceFuncStr string
	//go:embed template/enum/ptr_slice.tmpl
	enumSlicePtrFuncStr string

	enumSingleTmpl    *template.Template
	enumPtrSingleTmpl *template.Template
	enumSliceTmpl     *template.Template
	enumSlicePtrTmpl  *template.Template
)

func init() {
	var err error
	{
		enumSingleTmpl, err = template.
			New("invalid_enum").
			Funcs(funcMap).
			Parse(enumSingleFuncStr)
		if err != nil {
			panic(err)
		}

		enumPtrSingleTmpl, err = template.
			New("invalid_enum_ptr").
			Funcs(funcMap).
			Parse(enumPtrSingleFuncStr)
		if err != nil {
			panic(err)
		}

		enumSliceTmpl, err = template.
			New("invalid_enum_slice").
			Funcs(funcMap).
			Parse(enumSliceFuncStr)
		if err != nil {
			panic(err)
		}

		enumSlicePtrTmpl, err = template.
			New("invalid_enum_slice_ptr").
			Funcs(funcMap).
			Parse(enumSlicePtrFuncStr)
		if err != nil {
			panic(err)
		}
	}
}

type enumRule struct {
	Index      int
	Rule       string
	Tags       string
	Message    string
	FieldName  string
	FieldType  string
	StructName string
}

func (nr *enumRule) Prio() int {
	return 2
}

func (nr *enumRule) Check() string {
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
		err = enumSlicePtrTmpl.Execute(bs, tmplVal)
	} else if strings.HasPrefix(nr.FieldType, "[]") {
		err = enumSliceTmpl.Execute(bs, tmplVal)
	} else if strings.HasPrefix(nr.FieldType, "*") {
		err = enumPtrSingleTmpl.Execute(bs, tmplVal)
	} else {
		err = enumSingleTmpl.Execute(bs, tmplVal)
	}
	if err != nil {
		log.Fatal(err)
	}
	return bs.String()
}

func NewEnumRule(structName, fieldName, fieldType string, rule *Rule) *enumRule {
	INDEX++

	return &enumRule{
		Rule:       rule.Rule(),
		Tags:       rule.Tags(),
		FieldName:  fieldName,
		FieldType:  fieldType,
		Index:      INDEX,
		Message:    rule.Message(),
		StructName: structName,
	}
}
