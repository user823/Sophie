package vo

type TreeSelect struct {
	// 节点id
	Id int64
	// 节点名称
	Label string
	// 子节点
	Children []TreeSelect
}
