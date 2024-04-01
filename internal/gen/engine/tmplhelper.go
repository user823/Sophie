package engine

import (
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strings"
	"text/template"
)

/*
	包含模板引擎需要用到的功能函数
	它们以自己的函数名为key 注册到engine中
*/

// 加载引擎功能函数
func loadTemplateFuncs(parser *template.Template) {
	parser.Funcs(template.FuncMap{
		"getColumnComment":     getColumnComment,
		"require":              require,
		"replace":              replace,
		"removePrefix":         removePrefix,
		"add":                  add,
		"sub":                  sub,
		"toThrift":             toThrift,
		"Capitalize":           strutil.Capitalize,
		"CamelCaseToSnakeCase": strutil.CamelCaseToSnakeCase,
		"Uncapitalize":         strutil.Uncapitalize,
	})
}

// 获取列注释
func getColumnComment(comment string) string {
	if strings.Contains(comment, "(") {
		index := strings.Index(comment, "(")
		return comment[:index]
	}
	return comment
}

// 判断列属性是否必须
func require(prop string) bool {
	return prop == v1.REQUIRE
}

func replace(str string, old string, new string) string {
	return strings.ReplaceAll(str, old, new)
}

// 去除业务前缀
func removePrefix(str string) string {
	if str == "" {
		return str
	}
	prefix := "Sys"
	if strings.HasPrefix(str, prefix) {
		return str[len(prefix)+1:]
	}
	return str
}

// 数字类型加
func add(num, c int) int {
	return num + c
}

// 数字类型减法
func sub(num, c int) int {
	return num - c
}

// golang类型转化成thrift类型
func toThrift(goType string) string {
	if goType == "int64" {
		return "i64"
	} else if goType == "float64" {
		return "double"
	}
	return "string"
}

// 额外注册：
// strutil.Capitalize
// strutil.CamelCaseToSnakeCase
// strutil.Uncapitalize
