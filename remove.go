package astool

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
)

// RemoveBytes 删除指定结构体的指定函数
// 删除指定结构体的指定函数
// src         : 源文件内容
// structs     : 指定结构体名
// varCb       : 如果返回 true, 则删除该变量
// methodCb    : 如果返回 true 且属于上述结构体, 则删除该函数
func RemoveBytes(src []byte, structs []string, varCb, methodCb func(string) bool) []byte {
	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.AllErrors|parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	var ps []Pair[token.Pos]
	for _, v := range file.Decls {
		switch gen := v.(type) {
		case *ast.GenDecl:
			for _, v := range gen.Specs {
				switch spec := v.(type) {
				case *ast.ValueSpec:
					if varCb(spec.Names[0].Name) {
						ps = append(ps, MakePair(gen.Pos()-1, gen.End()))
						break
					}
				}
			}
		case *ast.FuncDecl:
			if gen.Recv == nil || !methodCb(gen.Name.Name) {
				continue
			}
			for _, v := range gen.Recv.List {
				switch expr := v.Type.(type) {
				case *ast.StarExpr:
					switch x := expr.X.(type) {
					case *ast.Ident:
						if In(structs, x.Name) {
							ps = append(ps, MakePair(gen.Pos()-1, gen.End()))
							break
						}
					case *ast.IndexExpr:
						if In(structs, x.X.(*ast.Ident).Name) {
							ps = append(ps, MakePair(gen.Pos()-1, gen.End()))
							break
						}
					}
				}
			}
		}
	}
	// 去除 content 指定部分
	for i := len(ps) - 1; i >= 0; i-- {
		fornt := src[:ps[i].First]
		end := src[ps[i].Second:]
		src = append(fornt, end...)
	}
	return src
}

func ExistFunc(pkg *packages.Package, name string) bool {
	for _, v := range pkg.Syntax {
		for _, v := range v.Decls {
			if v, ok := v.(*ast.FuncDecl); ok {
				if v.Name.Name == name {
					return true
				}
			}
		}
	}
	return false
}

func In(vals []string, target string) bool {
	for _, v := range vals {
		if v == target {
			return true
		}
	}
	return false
}
