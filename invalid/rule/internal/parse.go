package internal

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/inysc/astool/utils"
)

type (
	Pair = utils.Pair[string]
)

type Rule struct {
	Iden  string
	Pairs [4]Pair
}

func (r *Rule) Rule(iden, typ string) string {
	rule := r.Get("rule")
	if iden == Enum && rule[0] == '{' {
		return fmt.Sprintf("[]%s%s", remove_slice_ptr(typ), rule)
	}
	return rule
}

func (r *Rule) Message() string {
	message := r.Get("message")
	if strings.HasSuffix(message, ".Error()") {
		return strings.TrimSuffix(message, ".Error()")
	}
	if message == "" {
		return ""
	}
	return fmt.Sprintf("errors.New(%s)", message)
}

func (r *Rule) Tags() string {
	tags := r.Get("tags")

	if len(tags) > 0 && tags[0] == '{' {
		return "[]string" + tags
	}
	return tags
}

func (r *Rule) Get(key string) string {
	for _, v := range r.Pairs {
		if v.First == key {
			return strings.TrimSpace(v.Second)
		}
	}
	return ""
}

// for test
func (r *Rule) Equal(r2 *Rule) bool {
	if r == nil && r2 == nil {
		return true
	} else if r == nil || r2 == nil ||
		r.Iden != r2.Iden || len(r.Pairs) != len(r2.Pairs) {
		fmt.Printf("r1<%+v>, r2<%+v>\n", r, r2)
		return false
	}
	for idx := range r.Pairs {
		if r.Pairs[idx].First != r2.Pairs[idx].First ||
			r.Pairs[idx].Second != r2.Pairs[idx].Second {

			fmt.Printf("idx<%d>, r1.Pair<%+v>, r2.Pair<%+v>\n", idx, r.Pairs[idx], r2.Pairs[idx])
			return false
		}
	}
	return true
}

var REG = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+\s*=`) // 匹配以 `xxx =` 开头

func ParseComment(rule string) (r *Rule) {
	r = &Rule{}
	rule = strings.TrimSpace(rule)
	idx := 0
	for ; idx < len(rule); idx++ {
		if !isWord(rule[idx]) {
			break
		}
	}

	r.Iden = rule[:idx]
	if idx == len(rule) {
		return
	}

	rule = strings.TrimSpace(rule[idx:]) // 去掉标识符
	if rule[0] != '(' || rule[len(rule)-1] != ')' {
		panic("After the identifier, expect an open parenthesis to begin and a close parenthesis to end")
	}

	rule = rule[1 : len(rule)-1] // 去掉左右括号
	r.splitKeyValues(0, []byte(rule))
	return
}

func (r *Rule) splitKeyValues(index int, rule []byte) {

	rule = bytes.TrimSpace(rule)
	if len(rule) == 0 {
		return
	}

	var idx int
	if REG.Match(rule) {
		// 开始查找 key
		for idx = range rule {
			if isWord(rule[idx]) {
				continue
			}
			r.Pairs[index].First = string(rule[:idx])
			rule = rule[idx:]
			break
		}

		for idx = range rule {
			if unicode.IsSpace(rune(rule[idx])) {
				continue
			}
			//  == '=' // 由正则匹配保证
			rule = rule[idx+1:]
			break
		}
	}

	// 开始查找 value
	idx = FindValue(rule)
	if idx == -1 {
		panic("expected a value")
	}
	r.Pairs[index].Second = string(bytes.TrimSpace(rule[:idx]))
	if r.Pairs[index].First == "" {
		r.Pairs[index].First = r.KeyDefault(index)
	}
	if r.Pairs[index].First != r.KeyDefault(index) {
		log.Fatalf("unknown key<%s> in rule<%s>", r.Pairs[index].First, r.Iden)
	}

	if idx < len(rule) {
		r.splitKeyValues(index+1, rule[idx+1:])
	}
}

// 此处返回的 idx 是 `,` 的下标
// 如果没有，那就是 rule 的长度
func FindValue(rule []byte) int {
	inRaw := false // ` 在 raw 中
	inStr := false // " 在双引号范围中
	inBae := false // {
	idx := 0
	for idx = range rule {
		switch rule[idx] {
		case '`':
			if !inStr {
				inRaw = !inRaw
			}
		case '"':
			if !inRaw || (idx > 0 && rule[idx-1] != '\\') {
				inStr = !inStr
			}
		case '{', '}':
			if !inRaw && !inBae {
				inBae = !inBae
			}
		case ',':
			if !inRaw && !inStr && !inBae {
				return idx
			}
		}
	}
	if idx == len(rule)-1 {
		return len(rule)
	}

	return -1
}

func isWord(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

func (r *Rule) KeyDefault(idx int) string {
	if r.Iden == "RangeTime" && idx > 0 {
		if idx == 1 {
			return "format"
		}
		idx--
	} else if r.Iden == "Duck" || r.Iden == "Unique" {
		idx++
	}

	return [...]string{"rule", "message", "tags"}[idx]
}
