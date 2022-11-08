package internal

import (
	"fmt"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"len_gt_0":         lenGt0,
	"remove_ptr":       remove_ptr,
	"remove_slice":     remove_slice,
	"remove_slice_ptr": remove_slice_ptr,
	"trans_string":     tranString,
	"tRo":              tRo,
	"rRo":              rRo,
}

func lenGt0(s string) bool {
	return len(s) > 0
}

func remove_ptr(typ string) string {
	return strings.TrimPrefix(typ, "*")
}

func remove_slice(typ string) string {
	return strings.TrimPrefix(typ, "[]")
}

func remove_slice_ptr(typ string) string {
	return strings.ReplaceAll(strings.ReplaceAll(typ, "*", ""), "[]", "")
}

// 取反关系运算符
func rRo(l string) string {
	switch l {
	case ">":
		return "<="
	case "<":
		return ">="
	case ">=":
		return "<"
	case "<=":
		return ">"
	}
	return l
}

// 翻译关系运算符
func tRo(l string) string {
	switch l {
	case ">":
		return "greater than"
	case "<":
		return "less than"
	case ">=":
		return "greater than or equal to"
	case "<=":
		return "less than or equal to"
	}
	return l
}

func tranString(key string) string {
	if key != "" && key[0] == '\'' && key[len(key)-1] == '\'' {
		return fmt.Sprintf(`"%s"`, key[1:len(key)-1])
	}
	return key
}
