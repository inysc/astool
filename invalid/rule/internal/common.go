package internal

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/inysc/astool"
)

type commonRule struct {
	Iden       string
	Rule       string
	Tags       string
	FieldName  string
	FieldType  string
	Index      int
	Message    string
	StructName string
}

func (cr *commonRule) Prio() int {
	if cr.Iden == "Default" {
		return 0
	}
	return 2
}

func (cr *commonRule) Check() string {
	bs := astool.NewBytes()
	tmplVal := map[string]any{
		"iden":        cr.Iden,
		"rule":        cr.Rule,
		"tags":        cr.Tags,
		"field_name":  cr.FieldName,
		"field_type":  cr.FieldType,
		"index":       cr.Index,
		"message":     cr.Message,
		"struct_name": cr.StructName,
	}
	typ := cr.FieldType
	if cr.Iden == "Not" && strings.HasPrefix(cr.Rule, "nil") {
		typ = strings.TrimPrefix(typ, "*")
	}
	err := GetTmpl(cr.Iden, typ).Execute(bs, tmplVal)
	if err != nil {
		log.Fatal(err)
	}
	return bs.String()
}

func (cr *commonRule) Vars() string {
	if cr.Iden == "Regexp" {
		return fmt.Sprintf("var _%s_%d = regexp.MustCompile(%s)\n", cr.StructName, cr.Index, cr.Rule)
	}
	return ""
}

func NewCommonRule(structName, fieldName, fieldType string, rule *Rule) *commonRule {
	INDEX++

	return &commonRule{
		Iden:       rule.Iden,
		Rule:       rule.Rule(),
		Tags:       rule.Tags(),
		FieldName:  fieldName,
		FieldType:  fieldType,
		Index:      INDEX,
		Message:    rule.Message(),
		StructName: structName,
	}
}

func GetTmpl(iden, typ string) *template.Template {
	if strings.HasPrefix(typ, "*[]") {
		typ = "ptrslice"
	} else if strings.HasPrefix(typ, "[]") {
		typ = "slice"
	} else if strings.HasPrefix(typ, "*") {
		typ = "ptrsingle"
	} else {
		typ = "single"
	}

	val := map[string]map[string]*template.Template{
		"Default": {
			"single":    defaultSingleTmpl,
			"ptrsingle": defaultPtrSingleTmpl,
			"slice":     defaultSingleTmpl,
			"ptrslice":  defaultPtrSingleTmpl,
		},
		"Duck": {
			"single":    duckSingleTmpl,
			"ptrsingle": duckPtrSliceTmpl,
			"slice":     duckSliceTmpl,
			"ptrslice":  duckPtrSliceTmpl,
		},
		"Enum": {
			"single":    enumSingleTmpl,
			"slice":     enumSliceTmpl,
			"ptrsingle": enumPtrSingleTmpl,
			"ptrslice":  enumPtrSliceTmpl,
		},
		"Equal": {
			"single":    equalSingleTmpl,
			"slice":     equalSliceTmpl,
			"ptrsingle": equalPtrSingleTmpl,
			"ptrslice":  equalPtrSliceTmpl,
		},
		"Length": {
			"single":    lengthSingleTmpl,
			"slice":     lengthSliceTmpl,
			"ptrsingle": lengthPtrSingleTmpl,
			"ptrslice":  lengthPtrSliceTmpl,
		},
		"LengthUtf8": {
			"single":    lengthSingleTmpl,
			"slice":     lengthSliceTmpl,
			"ptrsingle": lengthPtrSingleTmpl,
			"ptrslice":  lengthPtrSliceTmpl,
		},
		"Not": {
			"single":    notSingleTmpl,
			"slice":     notSliceTmpl,
			"ptrsingle": notPtrSingleTmpl,
			"ptrslice":  notPtrSliceTmpl,
		},
		"Range": {
			"single":    rangeSingleTmpl,
			"slice":     rangeSliceTmpl,
			"ptrsingle": rangePtrSingleTmpl,
			"ptrslice":  rangePtrSliceTmpl,
		},
		"RangeTime": {
			"single":    rangeTimeSingleTmpl,
			"slice":     rangeTimeSliceTmpl,
			"ptrsingle": rangeTimePtrSingleTmpl,
			"ptrslice":  rangeTimePtrSliceTmpl,
		},
		"Regexp": {
			"single":    regexpSingleTmpl,
			"slice":     regexpSliceTmpl,
			"ptrsingle": regexpPtrSingleTmpl,
			"ptrslice":  regexpPtrSliceTmpl,
		},
		"Unique": {
			"single":    uniqueSingleTmpl,
			"ptrsingle": uniquePtrSingleTmpl,
			"slice":     uniqueSingleTmpl,
			"ptrslice":  uniquePtrSingleTmpl,
		},
	}

	return val[iden][typ]
}
