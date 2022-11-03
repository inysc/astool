package rule

import (
	"fmt"
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
	comments := strings.Split(field.Comment, "\n")
	for _, comment := range comments {
		comment = strings.TrimSpace(comment)
		if strings.HasPrefix(comment, "@iv") || strings.HasPrefix(comment, "@invalid") {
			// 去掉 @iv、空格
			comment = strings.TrimPrefix(comment, "@iv")
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
	fmt.Printf("rules: %v\n", rules)
	return ""
}
