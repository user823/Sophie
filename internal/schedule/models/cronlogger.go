package models

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/store"
	"github.com/user823/Sophie/pkg/log"
	"time"
)

// 实现cron 的Logger接口
// 将cron 的调度信息输出到sql 的jobLog 中
type SqlLogger struct {
	store store.JobLogStore
}

func NewLogger(s store.JobLogStore) cron.Logger {
	return &SqlLogger{s}
}

func (l *SqlLogger) Info(msg string, keysAndValues ...any) {
	sysJobLog := &v1.SysJobLog{}
	buildMessage(sysJobLog, msg, keysAndValues...)
	if sysJobLog.Status == v1.PAUSE || sysJobLog.JobMessage == "" {
		return
	}
	if err := l.store.InsertJobLog(context.Background(), sysJobLog, &api.CreateOptions{}); err != nil {
		log.Errorf("调度日志创建失败: %s", err.Error())
	}
}

func (l *SqlLogger) Error(err error, msg string, keysAndValues ...any) {
	sysJobLog := &v1.SysJobLog{ExceptionInfo: err.Error()}
	buildMessage(sysJobLog, msg, keysAndValues...)
	if sysJobLog.Status == v1.PAUSE || sysJobLog.JobMessage == "" {
		return
	}
	if err = l.store.InsertJobLog(context.Background(), sysJobLog, &api.CreateOptions{}); err != nil {
		log.Errorf("调度日志创建失败: %s", err.Error())
	}
}

func buildMessage(jobLog *v1.SysJobLog, msg string, kv ...any) {
	switch msg {
	case "schedule", "added":
		jobLog.JobMessage = "任务开始调度"
	case "run":
		jobLog.JobMessage = "任务被调度执行"
	case "removed":
		jobLog.JobMessage = "任务被移除"
	default:
		return
	}
	for i := 0; i < len(kv); i = i + 2 {
		if i+1 == len(kv) {
			break
		}
		if k, ok := kv[i].(string); ok {
			mapKeys(jobLog, k, kv[i+1])
		}
	}
}

// 将cron 的key 映射到joblog
func mapKeys(jobLog *v1.SysJobLog, key string, val any) {
	switch key {
	case "now":
		if t, ok := val.(time.Time); ok {
			jobLog.CreatedAt = t
		}
	case "entry":
		if id, ok := val.(cron.EntryID); ok {
			job := GetJobInfo(int64(id))
			fillJobLog(jobLog, job)
		}
	}
}

func fillJobLog(jobLog *v1.SysJobLog, job *v1.SysJob) {
	if job == nil {
		return
	}
	jobLog.JobName = job.JobName
	jobLog.JobGroup = job.JobGroup
	jobLog.InvokeTarget = job.InvokeTarget
	jobLog.Status = job.Status
}
