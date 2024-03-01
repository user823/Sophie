package app

import "github.com/user823/Sophie/pkg/ds"

// 获取命令行选项信息
type CliOptions interface {
	Flags() *ds.FlagGroup
}

// 程序运行前进行配置检查
type ValidatableOptions interface {
	Validate() error
}

type CompletableOptions interface {
	Complete() error
}

type PrintableOptions interface {
	String() string
}
