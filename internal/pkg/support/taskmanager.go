package support

import (
	"context"
	"sync"
)

// 简单任务管理器（非线程安全）
// 支持添加任务
// 支持等待所有任务执行完毕
// 支持等待其中一个任务执行完毕
// 支持统计任务完成数
type SimpleTaskManager struct {
	tasks []*GoTask
	wg    sync.WaitGroup
}

func NewTaskManager() *SimpleTaskManager {
	return &SimpleTaskManager{}
}

func (m *SimpleTaskManager) Add(tasks ...*GoTask) {
	m.tasks = append(m.tasks, tasks...)
}

// 同步等待所有任务执行完毕
func (m *SimpleTaskManager) All() (results []TaskResult) {
	m.wg.Add(len(m.tasks))
	results = make([]TaskResult, len(m.tasks))
	for i := range m.tasks {
		go func(i int) {
			defer m.wg.Done()
			m.tasks[i].Start()
			// 阻塞获取结果
			results[i] = m.tasks[i].Output(context.Background())
		}(i)
	}
	m.wg.Wait()
	return
}

// 同步等待任意一个任务执行完毕
func (m *SimpleTaskManager) Any() (res TaskResult) {
	ch := make(chan TaskResult, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m.wg.Add(len(m.tasks))
	for i := range m.tasks {
		go func(i int) {
			defer m.wg.Done()
			m.tasks[i].Start()

			for {
				select {
				case <-ctx.Done():
					m.tasks[i].Shutdown(NOW)
					return
				default:
					result := m.tasks[i].AsyncOutput()
					if result.OK {
						ch <- result
						return
					}
				}
			}
		}(i)
	}

	res = <-ch
	cancel()
	m.wg.Wait()
	return
}
