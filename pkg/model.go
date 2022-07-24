package pkg

import (
	"go/ast"
	"go/token"
	"reflect"
)

type Type struct {
	Ptr  bool
	Pkg  string
	Type string
}

func (t *Type) String() (ret string) {
	if t.Ptr {
		ret += "*"
	}
	if t.Pkg != "" {
		ret += t.Pkg + "."
	}
	return ret + t.Type
}

func NewType(expr ast.Expr) (typ *Type) {
	typ = &Type{}
	switch ex := expr.(type) {
	case *ast.StarExpr:
		typ.Ptr = true
		tmp := NewType(ex)
		typ.Pkg = tmp.Pkg
		typ.Type = tmp.Type
	case *ast.SelectorExpr:
		tmp := NewType(ex.X)
		typ.Pkg = tmp.Type
		typ.Type = ex.Sel.Name
	case *ast.Ident:
		typ.Type = ex.Name
	default:
		fset := token.NewFileSet()
		ast.Print(fset, expr)
		panic("please check the stdout")
	}
	return
}

type Tag = reflect.StructTag

type StructField struct {
	Name    string
	Type    *Type
	Tag     Tag
	Comment string
}

func (sf *StructField) Get(tagname string) string {
	return sf.Tag.Get(tagname)
}

func NewStructField(field *ast.Field) *StructField {
	tag := ""
	if field.Tag != nil {
		tag = field.Tag.Value
	}
	comment := ""
	if field.Comment != nil {
		comment = field.Comment.Text()
	}
	return &StructField{
		Name:    field.Names[0].Name,
		Type:    NewType(field.Type),
		Tag:     Tag(tag),
		Comment: comment,
	}
}

type StructInfo struct {
	Pkg    string
	Name   string
	Fields []*StructField
}

func NewStructInfo(pkg, name string, st *ast.StructType) *StructInfo {
	fields := make([]*StructField, 0, len(st.Fields.List))

	for _, field := range st.Fields.List {
		fields = append(fields, NewStructField(field))
	}

	return &StructInfo{
		Pkg:    pkg,
		Name:   name,
		Fields: fields,
	}
}

func EasyStructInfo(name, path string) (si *StructInfo) {
	pkg := ParsePackage([]string{path}, []string{})
	return NewStructInfo(pkg.Name, name, GetStruct(name, pkg.Syntax...))
}
