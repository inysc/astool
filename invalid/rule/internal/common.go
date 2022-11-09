package internal

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

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
	switch cr.Iden {
	case "Default":
		return 0
	case "Not":
		if strings.HasPrefix(cr.Rule, "nil") {
			return 1
		}
		return 2
	default:
		return 3
	}
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
	if (cr.Iden == Not || cr.Iden == Equal) && strings.HasPrefix(cr.Rule, "nil") {
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
		Rule:       rule.Rule(rule.Iden, fieldType),
		Tags:       rule.Tags(),
		FieldName:  fieldName,
		FieldType:  fieldType,
		Index:      INDEX,
		Message:    rule.Message(),
		StructName: structName,
	}
}
