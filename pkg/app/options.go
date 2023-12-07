package app

// 获取命令行选项信息
type CliOptions interface {
	Flags() FlagGroup
}

type ValidatableOptions interface {
	Validate() error
}

type CompletableOptions interface {
	Complete() error
}

type PrintableOptions interface {
	String() string
}
