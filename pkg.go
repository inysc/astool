package astool

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

const PkgMode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
	packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
	packages.NeedSyntax | packages.NeedTypesInfo

func ParsePackage(pattern string, tags []string) *packages.Package {
	cfg := &packages.Config{
		Mode:       PkgMode,
		Tests:      false, // 不应当支持测试文件
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	f, err := os.OpenFile("1.ast", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ast.Fprint(f, pkgs[0].Fset, pkgs[0].Syntax[0], nil)

	return pkgs[0]
}
