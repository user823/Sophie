package models

import (
	"github.com/user823/Sophie/pkg/ds"
	"sync"
)

var (
	// job死亡队列，表示因某种原因调度失败的任务
	// 后台程序可根据任务状态，重新调度或者丢弃
	DeadJobQueue *ds.LinkedHashSet[int64]
	donce        sync.Once
)

func InitDeadJobQueue() {
	donce.Do(func() {
		DeadJobQueue = ds.NewLinkedHashSet[int64]()
	})
}
