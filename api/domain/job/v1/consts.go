package v1

const (
	ServiceName = "Sophie Schedule"
)

// 任务调度通用常量
const (
	TASK_CLASS_NAME = "TASK_CLASS_NAME"
	// 执行目标key
	TASK_PROPERTIES = "TASK_PROPERTIES"
	// 默认
	MISFIRE_DEFAULT = "0"
	// 立刻触发执行
	MISFIRE_IGNORE_MISFIRES = "1"
	// 触发一次执行
	MISFIRE_FIRE_AND_PROCEED = "2"
	// 不触发立刻执行
	MISFIRE_DO_NOTHING = "3"
)
