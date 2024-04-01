package models

import (
	"fmt"
	"github.com/robfig/cron/v3"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	// 注册两条
	var cnt int
	var id cron.EntryID
	id = RegisterCronTask("*/1 * * * * ?", func() {
		fmt.Printf("test: %d\n", cnt)
		cnt++
	}, v1.SysJob{JobName: "test"})
	RegisterCronTask("* * * * *", func() {
		fmt.Println("hello, world")
	}, v1.SysJob{})
	ResumeJob(id)

	// 测试getinfo
	t.Logf("%v", GetJobInfo(id))
	t.Logf("注册好的条目: %d", len(Schedule.Entries()))
	time.Sleep(5 * time.Second)

	// 暂停任务
	PauseJob(id)
	time.Sleep(2 * time.Second)
	// 恢复任务
	ResumeJob(id)
	time.Sleep(2 * time.Second)
}

func TestCronSub(t *testing.T) {
	Schedule = cron.New(cron.WithSeconds(),
		cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)))
	tasks = make(map[cron.EntryID]*cronTask)
	Schedule.Start()

	t.Run("test-Register", TestCron)
}
