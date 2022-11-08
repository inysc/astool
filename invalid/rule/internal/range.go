package internal

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/inysc/astool"
)

type RangeRule struct {
	Iden       string
	Index      int
	Tags       string
	Message    string
	StructName string
	FieldName  string
	FieldType  string
	Format     string
	Left       string
	LeftVal    string
	Right      string
	RightVal   string
}

func (rr *RangeRule) Prio() int {
	return 2
}

func (rr *RangeRule) Check() string {
	bs := astool.NewBytes()
	tmplVal := map[string]any{
		"iden":        rr.Iden,
		"rule":        rr.Left,
		"tags":        rr.Tags,
		"field_name":  rr.FieldName,
		"field_type":  rr.FieldType,
		"index":       rr.Index,
		"message":     rr.Message,
		"struct_name": rr.StructName,
		"left":        rr.Left,
		"left_val":    rr.LeftVal,
		"right":       rr.Right,
		"right_val":   rr.RightVal,
		"format":      rr.Format,
	}
	err := GetTmpl(rr.Iden, rr.FieldType).Execute(bs, tmplVal)
	if err != nil {
		log.Fatal(err)
	}

	return bs.String()
}

func (rr *RangeRule) Vars() string {
	if rr.Iden == "RangeTime" {
		ret := ""
		if rr.LeftVal != "" && rr.LeftVal != "@now" {
			ret += fmt.Sprintf("var _%s_left_%d = mustTimeParse(%s, %s)\n", rr.StructName, rr.Index, rr.LeftVal, rr.Format)
		}
		if rr.RightVal != "" && rr.RightVal != "@now" {
			ret += fmt.Sprintf("var _%s_right_%d = mustTimeParse(%s, %s)\n", rr.StructName, rr.Index, rr.RightVal, rr.Format)
		}
		return ret
	}

	return ""
}

func NewRangeRule(structName, fieldName, fieldType string, rule *Rule) *RangeRule {
	INDEX++

	ruleval := rule.Rule()
	// 去除前后的引号
	ruleval = ruleval[1 : len(ruleval)-1]
	left, right := ">", "<"
	leftVal, rightVal := "", ""
	if ruleval[0] == '[' {
		left = ">="
	} else if ruleval[0] == '(' {
		left = ">"
	} else {
		log.Panicf("the unsupport rule<%s>", rule)
	}

	if ruleval[len(ruleval)-1] == ']' {
		right = "<="
	} else if ruleval[len(ruleval)-1] == ')' {
		right = "<"
	} else {
		log.Panicf("the unsupport rule<%s>", ruleval)
	}
	ruleval = strings.TrimSpace(ruleval[1 : len(ruleval)-1])

	inStr := false
	inRaw := false
	inByte := false
	for i, v := range ruleval {
		switch v {
		case '"':
			if !inRaw && !inByte {
				inStr = !inStr
			}
		case '`':
			if !inStr && !inByte {
				inRaw = !inRaw
			}
		case '\'':
			if !inStr && !inRaw {
				inByte = !inByte
			}
		case ',':
			if !inStr && !inRaw && !inByte {
				leftVal = strings.TrimSpace(ruleval[:i])
				rightVal = strings.TrimSpace(ruleval[i+1:])
				goto ok
			}
		}
	}
	log.Fatalf("the unsupport rule<%+v>", rule)
ok:
	return &RangeRule{
		Iden:       rule.Iden,
		Index:      INDEX,
		Tags:       rule.Tags(),
		Message:    rule.Message(),
		StructName: structName,
		FieldName:  fieldName,
		FieldType:  fieldType,
		Format:     rule.Get("format"),
		Left:       left,
		LeftVal:    leftVal,
		Right:      right,
		RightVal:   rightVal,
	}
}
