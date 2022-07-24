package pkg

import (
	"go/ast"
	"go/parser"
	"go/token"
)

var mode = parser.ParseComments | parser.DeclarationErrors | parser.AllErrors

func PraseFile(filename string) *ast.File {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, mode)
	if err != nil {
		panic(err)
	}

	return file
}

func GetStruct(structName string, files ...*ast.File) *ast.StructType {
	for _, file := range files {
		if file == nil || file.Scope == nil ||
			file.Scope.Objects == nil {
			break
		}
		for k, obj := range file.Scope.Objects {
			if k == structName {
				// don't worry about Panic
				return obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
			}
		}
	}
	return nil
}
