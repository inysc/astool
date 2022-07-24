package pkg

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

const PkgMode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
	packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
	packages.NeedSyntax | packages.NeedTypesInfo

func ParsePackage(patterns []string, tags []string) *packages.Package {
	cfg := &packages.Config{
		Mode:       PkgMode,
		Tests:      false, // 不应当支持测试文件
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	return pkgs[0]
}
