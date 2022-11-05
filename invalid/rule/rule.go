package rule

import (
	"sort"
	"strings"

	"github.com/inysc/astool"
	"github.com/inysc/astool/invalid/rule/internal"
)

type Ruler interface {
	Prio() int     // 优先级
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
			case "Not": // 禁止值
				rules = append(rules, internal.NewNotRule(structName, field.Name, field.Type.String(), rule))
			case "Equal": // 等于值
				rules = append(rules, internal.NewEqualRule(structName, field.Name, field.Type.String(), rule))
			case "Enum": // 枚举值
				rules = append(rules, internal.NewEnumRule(structName, field.Name, field.Type.String(), rule))
			case "Range": // 区间限制
			case "RangeTime": // 时间区间限制
			case "Length": // 区间限制
			case "LengthUtf8": // 区间限制
			case "Regexp": // 正则约束
			case "Default": // 默认值
			case "Invalid": // 标注已实现 interface { Invalid() error
			case "Unique": // 唯一值
			default:
			}
		}
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].Prio() > rules[j].Prio() })

	if len(rules) != 0 {
		bs.Pf("func (i *%[1]s)_%[1]s_Invalid_%[2]s(tags []string) (err error) {", structName, field.Name)
		for _, v := range rules {
			bs.P(v.Check())
		}
		bs.P("    return nil")
		bs.P("}\n")
	}

	return bs.String()
}
