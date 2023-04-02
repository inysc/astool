package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Pair struct{ First, Second token.Pos }

var (
	flag_file  = flag.String("file", "sqlstruct.go", "目的文件名")
	flag_types = flag.String("types", "", "结构体名，多个以英文逗号,分割")

	filename    string
	structNames []string

	DB_PTR_NAME   = regexp.MustCompile(`^D[Bb]Ptr[A-Z]`)                                           // 指针结构体的结构体名校验
	DB_FIELD_NAME = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)                                  // 字段名校验
	TARGET_FNS    = []string{"SQLValues", "SQLPtrNotNilValues", "SQLPtrNotPtrValues", "SQLFixPtr"} // 将为结构体生成这些函数
)

const (
	SQLValues          = "// 结构体里的每个字段的 tag 与值 \nfunc (i *%s) SQLValues() SQLPairs { return []SQLPair{%s} }\n"
	SQLPtrNotNilValues = "// 结构体里的每个不是 nil 的指针字段的 tag 与值 \nfunc (i *%s) SQLPtrNotNilValues() SQLPairs { %s }\n"
	SQLPtrNotPtrValues = "// 结构体里的每个不是指针字段的 tag 与值 \nfunc (i *%s) SQLPtrNotPtrValues() SQLPairs { return []SQLPair{%s} } \n"
	SQLFixPtr          = "// 为结构体里为 nil 的指针字段赋一个零值 \nfunc (i *%s) SQLFixPtr() { %s } \n" // 在 insert 时，执行该函数，给那些 nil 的字段赋一个零值
)

func RmFns(src []byte, structNames []string) []byte {
	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.AllErrors|parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	var ps []Pair
	for _, v := range file.Decls {
		switch fn := v.(type) {
		case *ast.FuncDecl:
			if fn.Recv == nil {
				continue
			}
			if !In(TARGET_FNS, fn.Name.Name) {
				continue
			}
			for _, v := range fn.Recv.List {
				switch expr := v.Type.(type) {
				case *ast.StarExpr:
					switch x := expr.X.(type) {
					case *ast.Ident:
						if In(structNames, x.Name) {
							start := fn.Pos() - 1
							if fn.Doc != nil {
								start = fn.Doc.Pos() - 1
							}
							ps = append(ps, Pair{start, fn.End()})
							break
						}
					case *ast.IndexExpr: // 泛型
						if In(structNames, x.X.(*ast.Ident).Name) {
							start := fn.Pos() - 1
							if fn.Doc != nil {
								start = fn.Doc.Pos() - 1
							}
							ps = append(ps, Pair{start, fn.End()})
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

func In(vals []string, target string) bool {
	for _, v := range vals {
		if v == target {
			return true
		}
	}
	return false
}

func StructParase(st *types.Struct) (sqlFields, sqlPtrTags, sqlPtrNames, sqlNotFields, sqlType []string) {
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		fieldName := field.Name()
		fieldTag := reflect.StructTag(st.Tag(i))

		if !field.Exported() { // 不处理私有字段
			continue
		}

		// 解析 tag 信息、字段名信息
		tag := ""
		for _, v := range []string{"sql", "SQL", "db"} {
			get := fieldTag.Get(v)
			if get == "-" {
				continue
			} else if get != "" {
				tag = get
				break
			}
		}
		if tag == "" {
			if field.Anonymous() { // 结构体匿名嵌入，认为是将其字段展开到当前结构体
				st1, ok := field.Type().Underlying().(*types.Struct)
				if ok {
					sf, spt, spn, snf, sfp := StructParase(st1)
					sqlFields = append(sqlFields, sf...)
					sqlPtrTags = append(sqlPtrTags, spt...)
					sqlPtrNames = append(sqlPtrNames, spn...)
					sqlNotFields = append(sqlNotFields, snf...)
					sqlType = append(sqlType, sfp...)
					continue
				}
			}
			tag = fieldName
		}
		if DB_FIELD_NAME.MatchString(tag) {
			tag = fmt.Sprintf("`%s`", tag)
		}
		tag = strconv.Quote(tag)
		sqlFields = append(sqlFields, fmt.Sprintf("{K:%s, V:i.%s}", tag, fieldName))
		ptype, ok := field.Type().(*types.Pointer)
		if ok {
			sqlPtrTags = append(sqlPtrTags, tag)
			sqlPtrNames = append(sqlPtrNames, fieldName)
			sqlType = append(sqlType, ptype.Elem().String())
		} else {
			sqlNotFields = append(sqlNotFields, fmt.Sprintf("{K:%s, V:i.%s}", tag, fieldName))
		}
	}
	return
}

// => \n{elem}{sep}{elem}{sep}
func Join(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return "\n" + elems[0] + sep
	}
	n := len(sep)*len(elems) + 2
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteByte('\n')
	for _, s := range elems {
		b.WriteString(s)
		b.WriteString(sep)
	}
	return b.String()
}

const mode = packages.NeedTypes | packages.NeedTypesInfo |
	packages.NeedDeps | packages.NeedImports

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	// 处理输入参数
	filename = *flag_file
	if !strings.HasSuffix(filename, ".go") {
		filename += ".go"
	}
	structNames = strings.Split(*flag_types, ",")
	if len(structNames) == 0 {
		log.Fatal("not found struct names")
	}

	// 解析当前目录
	pkgs, err := packages.Load(&packages.Config{Mode: mode}, ".")
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatal("len(pkgs) != 1")
	}
	pkg := pkgs[0]

	// 处理目的文件
	// 如果目的文件已存在就清除相关结构体的方法
	// 不存在就写入一些基本信息
	buff := bytes.Buffer{}
	file, err := os.ReadFile(filename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	if len(file) == 0 {
		buff.WriteString("// Code generated by invalid. DO NOT EDIT.\n")
		buff.WriteString("package " + pkg.Name + "\n")
	} else {
		buff.Write(RmFns(file, structNames)) // 移除目标文件中相关方法
	}

	scope := pkg.Types.Scope()

	obj := scope.Lookup("SQLPair")
	if obj == nil {
		buff.WriteString("type SQLPair struct{ K string; V interface{} }\n")
	} /* else {
			named, ok := obj.Type().(*types.Named)
			if !ok {
					log.Fatal("命名冲突 SQLPair")
			}
			_, ok = named.Underlying().(*types.Struct)
			if !ok {
					log.Fatal("命名冲突，SQLPair 不是结构体")
			}
			// 应该继续检查字段，不想做了
	}*/

	obj = scope.Lookup("SQLPairs")
	if obj == nil {
		buff.WriteString("type SQLPairs []SQLPair\n")
		buff.WriteString(`func (sp SQLPairs)Split() (ks []string,vs []interface{}) {
							for _, v := range sp {
									ks = append(ks, v.K)
									vs = append(vs, v.V)
							}
							return
					}` + "\n")
		buff.WriteString(`func (sp SQLPairs) ToMap() map[string]interface{} {
					mp := make(map[string]interface{}, len(sp))
					for _, v := range sp {
							mp[v.K] = v.V
					}
					return mp
			}` + "\n")
	} /* else {
			named, ok := obj.Type().(*types.Named)
			if !ok {
					log.Fatal("命名冲突 SQLPairs")
			}
			_, ok = named.Underlying().(*types.Slice)
			if !ok {
					log.Fatal("命名冲突，SQLPairs 不是切片")
			}
			// 应该继续检查 elem ，不想做了
	}*/

	// 遍历结构体及其字段
	for _, structName := range structNames {
		obj := scope.Lookup(structName)
		if obj == nil {
			log.Fatal("no found " + structName)
		}
		named, ok := obj.Type().(*types.Named)
		if !ok {
			log.Fatal("can not find a named type called " + structName)
		}
		st, ok := named.Underlying().(*types.Struct)
		if !ok {
			log.Fatal(structName + " is not a struct")
		}

		sqlFields, sqlPtrTags, sqlPtrNames, sqlNotFields, sqlType := StructParase(st)
		needPtr := DB_PTR_NAME.MatchString(structName)

		// 写入到缓存
		buff.WriteString(fmt.Sprintf(SQLValues, structName, Join(sqlFields, ",\n")))
		if needPtr {
			bufn := bytes.NewBufferString("vals := []SQLPair{}\n")
			bufm := bytes.Buffer{}
			template_n := "if i.%[1]s != nil { vals = append(vals, SQLPair{K:%[2]s, V:i.%[1]s }) }\n"
			template_m := "if i.%[1]s == nil { i.%[1]s = new(%[2]s) }\n"
			for i, v := range sqlPtrNames {
				bufn.WriteString(fmt.Sprintf(template_n, v, sqlPtrTags[i]))
				bufm.WriteString(fmt.Sprintf(template_m, v, sqlType[i])) // BUG:
			}
			bufn.WriteString("return vals")

			buff.WriteString(fmt.Sprintf(SQLPtrNotPtrValues, structName, Join(sqlNotFields, ", \n")))
			buff.WriteString(fmt.Sprintf(SQLPtrNotNilValues, structName, bufn.Bytes()))
			buff.WriteString(fmt.Sprintf(SQLFixPtr, structName, bufm.Bytes()))
		}
		buff.WriteByte('\n')
	}

	// 格式化代码，并写入到文件
	buf, err := format.Source(buff.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		file = buff.Bytes()
	} else {
		file = buf
	}
	err = os.WriteFile(filename, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
