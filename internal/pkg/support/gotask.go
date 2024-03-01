package support

import (
	"context"
	"github.com/user823/Sophie/pkg/log"
	"sync/atomic"
	"time"
)

const (
	INITIAL uint32 = iota
	STARTED
	STOPING
	EXIT
)

const (
	NOW = 0
)

// 用户执行Start前要提供Run方法
type GoTask struct {
	ServiceName string
	hasNotify   atomic.Bool
	cancelFn    context.CancelFunc
	status      atomic.Uint32
	Run         func(ctx context.Context) (interface{}, error)
	// 用于任务阻塞结束时回调
	OnWaitEnd   func()
	OnErrHandle func(err error)
	output      interface{}
}

type TaskResult struct {
	Data interface{}
	OK   bool
}

func (s *GoTask) Start() {
	if !s.status.CompareAndSwap(INITIAL, STARTED) {
		return
	}

	if s.Run == nil {
		log.Fatalf("Run function is empty, can't start task ")
		return
	}
	log.Infof("Try to start service: %s", s.ServiceName)
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFn = cancel

	go func() {
		defer cancel()
		result, err := s.Run(ctx)
		if err != nil && s.OnErrHandle != nil {
			s.OnErrHandle(err)
		}
		s.output = result
		s.status.Store(EXIT)
		log.Infof("Service Exit: %s", s.ServiceName)
	}()

	log.Infof("Start task successful: %s", s.ServiceName)
}

func (s *GoTask) Status() uint32 {
	return s.status.Load()
}

func (s *GoTask) Stoped() bool {
	return s.Status() == EXIT
}

func (s *GoTask) Wakeup() {
	if !s.hasNotify.CompareAndSwap(false, true) {
		return
	}
}

func (s *GoTask) WaitForRunning(interval time.Duration) {
	if s.hasNotify.CompareAndSwap(true, false) {
		if s.OnWaitEnd != nil {
			s.OnWaitEnd()
		}
		return
	}
	time.Sleep(interval)
	if s.OnWaitEnd != nil {
		s.OnWaitEnd()
	}
}

func (s *GoTask) Shutdown(delay time.Duration) {
	<-time.After(delay)
	log.Infof("Try to stop service: %s", s.ServiceName)
	if !s.status.CompareAndSwap(STARTED, STOPING) {
		return
	}
	s.Wakeup()
	s.cancelFn()
}

// 获取任务执行结果
func (s *GoTask) Output(ctx context.Context) TaskResult {
	for !s.Stoped() {
		select {
		case <-ctx.Done():
			return TaskResult{nil, false}
		default:
		}
	}
	return TaskResult{s.output, true}
}

// 异步获取任务执行结果
func (s *GoTask) AsyncOutput() TaskResult {
	return TaskResult{s.output, s.Stoped()}
}
