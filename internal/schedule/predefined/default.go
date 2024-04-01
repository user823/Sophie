package predefined

import (
	"fmt"
	"sync"
)

// 默认预定义的函数
var (
	Targets map[string]any
	once    sync.Once
)

func TestFunc() {
	fmt.Println("test success!")
}

func TestFuncWithParams(s string, b bool, l int64, d float64, i int) {
	fmt.Printf("执行多参方法: 字符串类型%s, 布尔类型%t, 长整形%d, 浮点型%f, 整形%d\n", s, b, l, d, i)
}

func TargetsInit() {
	once.Do(func() {
		Targets = make(map[string]any)
		Targets["TestFunc"] = TestFunc
		Targets["TestFuncWithParams"] = TestFuncWithParams
	})
}
