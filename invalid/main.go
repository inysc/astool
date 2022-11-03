package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/inysc/astool"
	"github.com/inysc/astool/invalid/rule"
)

var (
	file_input   = flag.String("file", "", "写入到的文件名")
	tags_input   = flag.String("tags", "", "指定 build tag 的")
	types_input  = flag.String("types", "", "待生成校验方法的结构体名")
	import_input = flag.String("imports", "", "规则中使用到的包")
)

var (
	file    string
	tags    []string
	types   []string
	imports []string
)

var CONTENT = astool.NewBytes()

func init() {
	flag.Parse()

	tags = strings.Split(*tags_input, ",")
	types = strings.Split(*types_input, ",")
	imports = strings.Split(*import_input, ",")
	if *file_input == "" {
		*file_input = types[0]
	}
	file = *file_input
}

func main() {
	log.SetFlags(log.Lshortfile)

	log.Printf("current ppid<%d>", os.Getppid())

	pkg := astool.EasyStructInfos(".", types, tags...)
	fmt.Printf("pkg: %v\n", pkg)
	fmt.Printf("file: %v\n", file)
	fmt.Printf("imports: %v\n", imports)
	for _, name := range types {
		Parse(pkg.Structs[name])
	}
}

func Parse(st *astool.StructInfo) {
	for _, v := range st.Fields {
		rule.NewRule(st.Name, v)
	}
}
