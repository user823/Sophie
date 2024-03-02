package strutil

import (
	"regexp"
	"strconv"
	"strings"
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
