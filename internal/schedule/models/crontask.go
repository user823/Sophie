package models

import (
	"context"
	"github.com/robfig/cron/v3"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/store"
	"github.com/user823/Sophie/pkg/log"
	"sync"
	"sync/atomic"
)

type RunFn func()

type cronTask struct {
	// 设置暂停状态
	isPause atomic.Bool
	run     RunFn
	// 任务其他信息
	job v1.SysJob
	// cronTaskId
	cronId cron.EntryID
}

func (c *cronTask) Pause() {
	c.isPause.Store(true)
}

func (c *cronTask) Resume() {
	c.isPause.Store(false)
}

func (c *cronTask) ChangeStatus() {
	if c.isPause.CompareAndSwap(true, false) {
		c.job.Status = v1.NORMAL
		return
	}
	c.isPause.CompareAndSwap(false, true)
	c.job.Status = v1.PAUSE
}

func (c *cronTask) Run() {
	if !c.isPause.Load() {
		c.run()
	}
}

var (
	// 将jobid 和task进行绑定
	tasks map[int64]*cronTask
	rw    sync.RWMutex
	// 全局调度器
	Schedule *cron.Cron
	once     sync.Once
)

// 获取所有任务
func GetAllJobs() (res []int64) {
	rw.RLock()
	defer rw.RUnlock()
	for k := range tasks {
		res = append(res, k)
	}
	return
}

// 获取任务数量
func GetJobNums() int {
	rw.RLock()
	defer rw.RUnlock()
	return len(tasks)
}

// 获取任务信息
func GetJobInfo(id int64) *v1.SysJob {
	rw.RLock()
	defer rw.RUnlock()
	if a, ok := tasks[id]; ok {
		return &a.job
	}
	return nil
}

// 创建定时任务
func RegisterCronTask(spec string, run RunFn, job v1.SysJob) cron.EntryID {
	task := &cronTask{run: run, job: job}
	cronId, err := Schedule.AddJob(spec, task)
	if err != nil {
		log.Error("创建定时任务失败: %s", err)
		return -1
	}
	if job.Status == v1.PAUSE {
		task.isPause.Store(true)
	}

	rw.Lock()
	defer rw.Unlock()
	task.cronId = cronId
	tasks[job.JobId] = task
	return cronId
}

// 暂停任务状态
func PauseJob(id int64) {
	rw.Lock()
	defer rw.Unlock()
	if a, ok := tasks[id]; ok {
		a.Pause()
		a.job.Status = v1.PAUSE
	}
}

// 恢复任务状态
func ResumeJob(id int64) {
	rw.Lock()
	defer rw.Unlock()
	if a, ok := tasks[id]; ok {
		a.Resume()
		a.job.Status = v1.NORMAL
	}
}

// 修改任务状态
func ChangeJobStatus(id int64) {
	rw.Lock()
	defer rw.Unlock()
	tasks[id].ChangeStatus()
}

// 删除任务
func DeleteJob(id int64) {
	rw.Lock()
	defer rw.Unlock()
	if a, ok := tasks[id]; ok {
		Schedule.Remove(a.cronId)
		delete(tasks, id)
	}
}

// 立刻执行一次任务
func RunJob(id int64) {
	rw.RLock()
	defer rw.RUnlock()
	// 如果存在任务，则获取调度器中的任务
	if a, ok := tasks[id]; ok {
		e := Schedule.Entry(a.cronId)
		if e.Valid() {
			e.WrappedJob.Run()
		}
	}
}

// 清除任务
func CleanJob() {
	rw.Lock()
	defer rw.Unlock()
	// 等待任务运行完毕
	<-Schedule.Stop().Done()

	// 重新初始化调度器
	tasks = make(map[int64]*cronTask)
	logger := NewLogger(store.Client().JobLogs())
	Schedule = cron.New(cron.WithSeconds(),
		cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
		cron.WithLogger(logger))
	Schedule.Start()
}

func Stop() context.Context {
	// 清空任务
	tasks = make(map[int64]*cronTask)
	return Schedule.Stop()
}

// 暂停所有任务
func PauseAll() {
	rw.Lock()
	defer rw.Unlock()

	for _, v := range tasks {
		v.Pause()
		v.job.Status = v1.PAUSE
	}
}

// 初始化调度器
func InitSchedule() {
	once.Do(func() {
		tasks = make(map[int64]*cronTask)
		//可用秒
		logger := NewLogger(store.Client().JobLogs())
		Schedule = cron.New(cron.WithSeconds(),
			cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
			cron.WithLogger(logger))
		Schedule.Start()
	})
}
