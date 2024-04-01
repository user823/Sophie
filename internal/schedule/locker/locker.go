package locker

import (
	"github.com/user823/Sophie/pkg/errors"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"sync"
)

// 分布式锁前缀
const (
	// 任务调度锁
	JOB_SCHEDULE = "sophie-schedule-locker-schedule"
	JOB_PREFIX   = "sophie-schedule-jobs"
)

var (
	// job调度信息读写锁
	mu     *recipe.RWMutex
	muonce sync.Once
)

// 初始化分布式锁
func InitMutex(session *concurrency.Session) {
	muonce.Do(func() {
		mu = recipe.NewRWMutex(session, JOB_PREFIX)
	})
}

func Lock() error {
	if mu == nil {
		return errors.New("etcd concurrent locker is nil!")
	}
	return mu.Lock()
}

func UnLock() error {
	if mu == nil {
		return errors.New("etcd concurrent locker is nil!")
	}
	return mu.Unlock()
}

func RLock() error {
	if mu == nil {
		return errors.New("etcd concurrent locker is nil!")
	}
	return mu.RLock()
}

func RUnLock() error {
	if mu == nil {
		return errors.New("etcd concurrent locker is nil!")
	}
	return mu.RUnlock()
}
