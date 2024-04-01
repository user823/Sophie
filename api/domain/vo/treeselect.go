package vo

type TreeSelect struct {
	// 节点id
	Id int64 `json:"id,omitempty"`
	// 节点名称
	Label string `json:"label,omitempty"`
	// 子节点
	Children []TreeSelect `json:"children"`
}
