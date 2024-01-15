package support

import (
	"context"
	"github.com/user823/Sophie/pkg/log"
	"sync/atomic"
	"time"
)

// Example:
//type MyService struct {
//	ServiceName string
//	Manager     ServiceManager
//}
//
//func (s *MyService) GetServiceName() string {
//	return s.ServiceName
//}
//
//func (s *MyService) Run(ctx context.Context) error {
//	if s.Manager == nil {
//		log.Fatalf("Run this service need a manager. ")
//		return nil
//	}
//	for s.Manager.Status() == STARTED {
//		fmt.Printf("Service %s is running\n", s.ServiceName)
//		s.Manager.WaitForRunning(1 * time.Second)
//	}
//	return nil
//}
//
//func (s *MyService) OnWaitEnd() {}
//
//func (s *MyService) GetManager() ServiceManager {
//	return s.Manager
//}
//

//	manager := DefaultServiceManager()
//	service := &MyService{"Test Service", manager}
//	manager.Start(service)
//	time.Sleep(3 * time.Second)
//	manager.Shutdown(0)
//	time.Sleep(10 * time.Second)

const (
	INITIAL uint32 = iota
	STARTED
	STOPING
	NORMALEXIT
	EXCEPTIONEXIT
	TIMEOUTEXIT
)

const (
	JOINTIME = 90 * 1000
)

type GoService struct {
	ServiceName string
	JoinTime    time.Duration
	hasNotify   atomic.Bool
	cancelFn    context.CancelFunc
	stopCh      chan error
	status      atomic.Uint32
	Run         func(ctx context.Context) error
	OnWaitEnd   func()
}

func (s *GoService) Start() {
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
	s.stopCh = make(chan error, 1)
	go func() {
		s.stopCh <- s.Run(ctx)
	}()
	log.Infof("Start task successful: %s", s.ServiceName)
}

func (s *GoService) Status() uint32 {
	return s.status.Load()
}

func (s *GoService) GetJoinTime() time.Duration {
	return s.JoinTime
}

func (s *GoService) Wakeup() {
	if !s.hasNotify.CompareAndSwap(false, true) {
		return
	}
}

func (s *GoService) WaitForRunning(interval time.Duration) {
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

func (s *GoService) Shutdown(delay time.Duration) {
	<-time.After(delay)
	log.Infof("Try to stop service: %s", s.ServiceName)
	if !s.status.CompareAndSwap(STARTED, STOPING) {
		return
	}
	s.Wakeup()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.JoinTime)
		defer cancel()
		select {
		case err := <-s.stopCh:
			if err != nil {
				s.status.Store(EXCEPTIONEXIT)
				log.Warnf("Task run exception:[%s %s] ", s.ServiceName, err.Error())
				return
			}
			s.status.Store(NORMALEXIT)
			log.Infof("Task has exit: %s", s.ServiceName)
		case <-ctx.Done():
			log.Warnf("Task exit timeout: %s", s.ServiceName)
			s.status.Store(TIMEOUTEXIT)
		}
	}()
}
