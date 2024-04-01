package excelutil

// xlsx 标签使用格式：
// key1:value1;key2:value2

// excel 格式控制属性
const (
	NAME   = "n"
	WIDTH  = "w"
	HEIGHT = "h"
	// 文字后缀
	SUFFIX = "s"
	// 读取内容时转换表达式
	// 格式: k=v,k=v,k=v
	READCONVERTEXP = "exp"
	// 结构体递归处理
	INLINE = "inline"
)

// excel 默认值
const (
	SEPARATOR = ";"
)
