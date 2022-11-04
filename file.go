package astool

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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

func ReadFile(file string) ([]byte, error) {
	// 判断文件是否存在
	bs, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil
		}
		return nil, err
	}
	return bs, nil
}

func MustReadFile(file string) []byte {
	bs, err := ReadFile(file)
	if err != nil {
		panic(err)
	}
	return bs
}
