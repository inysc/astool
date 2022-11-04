package astool

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

func RemoveBytes(content []byte, structNames ...string) []byte {
	file, err := parser.ParseFile(token.NewFileSet(), "", content, parser.AllErrors|parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	ps := make([]Pair[token.Pos], 0, len(structNames))
	for _, v := range file.Decls {
		switch gen := v.(type) {
		case *ast.GenDecl:
			for _, v := range gen.Specs {
				switch spec := v.(type) {
				case *ast.ValueSpec:
					for _, structName := range structNames {
						if strings.HasPrefix(spec.Names[0].Name, "_"+structName+"_") {
							ps = append(ps, MakePair(spec.Pos()-1, spec.End()))
							break
						}
					}
				}
			}
		case *ast.FuncDecl:
			if gen.Recv == nil {
				continue
			}
			for _, v := range gen.Recv.List {
				switch expr := v.Type.(type) {
				case *ast.StarExpr:
					switch x := expr.X.(type) {
					case *ast.Ident:
						for _, structName := range structNames {
							if x.Name == structName {
								ps = append(ps, MakePair(x.Pos()-1, x.End()))
								break
							}
						}
					case *ast.IndexExpr:
						for _, structName := range structNames {
							if x.X.(*ast.Ident).Name == structName {
								ps = append(ps, MakePair(x.Pos()-1, x.End()))
								break
							}
						}
					}
				}
			}
		}
	}
	for i := len(ps) - 1; i >= 0; i-- {
		// 去除 content 指定部分
		fornt := content[:ps[i].First]
		end := content[ps[i].Second:]
		content = append(fornt, end...)
	}
	return content
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
