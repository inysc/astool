package rule

import (
	"log"
	"sort"
	"strings"

	"github.com/inysc/astool"
	"github.com/inysc/astool/invalid/rule/internal"
)

type Ruler interface {
	Prio() int     // 优先级
	Vars() string  // 使用的全局变量初始化
	Check() string // 对应的方法内容
}

var INDEX = 0

func NewRule(structName string, field *astool.StructField) string {
	rules := []Ruler{}
	bs := astool.NewBytes()
	for _, comment := range strings.Split(field.Comment, "\n") {

		comment = strings.TrimSpace(comment)
		if idx := strings.Index(comment, "@iv"); idx != -1 {
			// 去掉 @iv、空格
			comment = comment[idx+3:]
			comment = strings.TrimSpace(comment)
			rule := internal.ParseComment(comment)

			switch rule.Iden {
			case "Not", "Equal", "Enum", "Regexp", "Default", "Duck", "Unique":
				rules = append(rules, internal.NewCommonRule(structName, field.Name, field.Type.String(), rule))
			case "Range", "RangeTime", "Length", "LengthUtf8":
				rules = append(rules, internal.NewRangeRule(structName, field.Name, field.Type.String(), rule))
			default:
				log.Fatalf("unknown rule: <%s>", rule.Iden)
			}
		}
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].Prio() < rules[j].Prio() })

	if len(rules) != 0 {
		for _, rule := range rules {
			bs.WriteString(rule.Vars())
		}
		bs.Pf("func (i *%[1]s)_%[1]s_Invalid_%[2]s(tags []string) (err error) {", structName, field.Name)
		for _, rule := range rules {
			bs.P(rule.Check())
		}
		bs.P("    return nil")
		bs.P("}\n")
	}

	return bs.String()
}
