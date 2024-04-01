package strutil

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Compare func(str string, target string) bool

func CompareAll(compare Compare, target string, searchStrings ...string) bool {
	for i := range searchStrings {
		if !compare(searchStrings[i], target) {
			return false
		}
	}
	return true
}

func CompareAny(compare Compare, target string, searchStrings ...string) bool {
	for i := range searchStrings {
		if compare(searchStrings[i], target) {
			return true
		}
	}
	return false
}

func ContainsAny(str string, searchStrings ...string) bool {
	return CompareAny(func(str1, str2 string) bool { return str1 == str2 }, str, searchStrings...)
}

func ContainsAnyIgnoreCase(str string, searchStrings ...string) bool {
	return CompareAny(func(str1, str2 string) bool { return strings.ToLower(str1) == strings.ToLower(str2) }, str, searchStrings...)
}

// 简单模式匹配（*）
func SimpleMatch(pattern string, target string) bool {
	// 检查 第一个'*'
	firstPos := strings.Index(pattern, "*")
	if firstPos == -1 {
		return pattern == target
	}
	// 检查 '*' 开头
	if firstPos == 0 {
		// 检查是否只有一个*
		if len(pattern) == 1 {
			return true
		}
		// 检查下一个 '*'
		nextPos := strings.Index(pattern[1:], "*")
		// 不存在下一个'*'
		if nextPos == -1 {
			return strings.HasSuffix(target, pattern[1:])
		}
		nextPos++

		part := pattern[1:nextPos]
		if part == "" {
			return SimpleMatch(pattern[1:], target)
		}

		// 检查中间部分
		offset := 0
		for partPos := strings.Index(target[offset:], part); partPos != -1; partPos = strings.Index(target[offset:], part) {
			if SimpleMatch(pattern[nextPos:], target[partPos+offset+len(part):]) {
				return true
			}
			offset = offset + partPos + 1
		}
		return false
	}
	return len(target) >= firstPos && strings.HasPrefix(target, pattern[:firstPos]) && SimpleMatch(pattern[firstPos:], target[firstPos:])
}

func EqualIgnoreCase(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

func IsHttp(path string) bool {
	return strings.HasPrefix(path, "http") || strings.HasPrefix(path, "https")
}

func ReplaceEach(input string, old []string, new []string) string {
	// 构建正则表达式
	regex := regexp.MustCompile("(" +
		"(" + regexp.QuoteMeta(old[0]) + ")" +
		"|(" + regexp.QuoteMeta(old[1]) + ")" +
		"|(" + regexp.QuoteMeta(old[2]) + ")" +
		"|\\." +
		"|:" +
		")")

	// 替换匹配的部分
	result := regex.ReplaceAllStringFunc(input, func(match string) string {
		switch match {
		case old[0], old[1], old[2]:
			return ""
		case ".", ":":
			return "/"
		default:
			return match
		}
	})

	return result
}

func Strs2Int64(strs string) (res []int64) {
	for _, param := range strings.Split(strs, ",") {
		num, _ := strconv.ParseInt(param, 10, 64)
		res = append(res, num)
	}
	return
}

func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func Uncapitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// 转化驼峰命名法（将下划线表示法转化为驼峰表示法）
func ToCamelCase(s string) string {
	if s == "" || !strings.Contains(s, string(SEPARATOR)) {
		return s
	}
	s = strings.ToLower(s)
	var buffer strings.Builder
	upperCase := false
	for _, r := range s {
		if r == SEPARATOR {
			upperCase = true
		} else if upperCase {
			buffer.WriteRune(unicode.ToUpper(r))
			upperCase = false
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

// 转化驼峰命名法（将下划线大写命名的字符串转换为驼峰式，HELLO_WORLD->HelloWorld)
func ConvertToCamelCase(name string) string {
	if name == "" {
		return ""
	} else if strings.Index(name, "_") == -1 {
		return Capitalize(name)
	}
	var builder strings.Builder
	camels := strings.Split(name, "_")
	for i := range camels {
		if camels[i] == "" {
			continue
		}
		builder.WriteString(Capitalize(camels[i]))
	}
	return builder.String()
}

// 转化成蛇形命名法
func CamelCaseToSnakeCase(s string) string {
	var buf strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			buf.WriteRune('_')
		}
		buf.WriteRune(unicode.ToLower(r))
	}
	return buf.String()
}

// 取子串
func SubStringBetween(str, open, close string) string {
	if str == "" || open == "" || close == "" {
		return ""
	}

	start := strings.Index(str, open)
	end := strings.Index(str, close)
	if start == -1 || end == -1 {
		return ""
	}
	return str[start+1 : end]
}

func EndsWithIgnoreCase(str string, suffix string) bool {
	return strings.HasSuffix(strings.ToLower(str), suffix)
}
