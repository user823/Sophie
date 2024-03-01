package evict

// 淘汰策略
type EvictPolicy interface {
	// 添加key，如果超过了限制则返回要淘汰的key
	Add(key string)
	Remove(key string)
	// 剩余空间大小
	Remains() int
	// 清空key
	Clean()
}
