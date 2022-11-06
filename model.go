package astool

import (
	"go/ast"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/packages"
)

type Type struct {
	Ptr  bool
	Arr  bool
	Pkg  string
	Type string
}

func (t *Type) String() (ret string) {
	if t.Ptr {
		ret += "*"
	}
	if t.Arr {
		ret += "[]"
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
		tmp := NewType(ex.X)
		typ.Pkg = tmp.Pkg
		typ.Type = tmp.Type
	case *ast.SelectorExpr:
		tmp := NewType(ex.X)
		typ.Pkg = tmp.Type
		typ.Type = ex.Sel.Name
	case *ast.Ident:
		typ.Type = ex.Name
	case *ast.ArrayType:
		tmp := NewType(ex.Elt)
		typ.Arr = true
		typ.Pkg = tmp.Pkg
		typ.Type = tmp.Type
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
		tag = tag[1 : len(tag)-1]
	}
	comment := field.Doc.Text()
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
	pkg := ParsePackage(path, []string{})
	return NewStructInfo(pkg.Name, name, GetStruct(name, pkg.Syntax...))
}

type Package struct {
	Pkg     *packages.Package
	Structs map[string]*StructInfo
}

func EasyStructInfos(path string, name []string, tags ...string) *Package {
	pkg := Package{
		Pkg:     ParsePackage(path, []string{}),
		Structs: make(map[string]*StructInfo, len(name)),
	}
	set := NewSet(name...)
	for _, v := range pkg.Pkg.Syntax {
		for _, v := range v.Decls {
			if decl, ok := v.(*ast.GenDecl); ok {
				for _, v := range decl.Specs {
					if ts, ok := v.(*ast.TypeSpec); ok {
						if set.Has(ts.Name.Name) {
							pkg.Structs[ts.Name.Name] = NewStructInfo(pkg.Pkg.Name, ts.Name.Name, ts.Type.(*ast.StructType))
							set.Delete(ts.Name.Name)
						}
					}
				}
			}
		}
	}
	return &pkg
}
