package shutdown

import (
	"fmt"
	"github.com/user823/Sophie/pkg/log"
	"sync"
)

// 优雅关停回调函数
type ShutdownCallback func(string) error

// 错误处理
// 对于处理不了的错误要么传递、要么panic
type ErrHandler interface {
	OnError(err error)
	SetDeliver(eh ErrHandler)
}

type EmptyErrHandler struct{}

func (e EmptyErrHandler) OnError(err error)        {}
func (e EmptyErrHandler) SetDeliver(eh ErrHandler) {}

// GraceFulShutdown实例接口
// 负责添加回调函数、关停同一的错误处理等
// 访问者模式
type GSInstance interface {
	// 由ms执行关停逻辑
	StartShutdown(ms ShutdownManager)
	SetErrHandler(eh ErrHandler)
	// callbacks执行具体关停动作
	AddShutdownCallbacks(shutdowns ...ShutdownCallback)
}

// GSInstance 管理者
// 执行关停逻辑
type ShutdownManager interface {
	GetName() string
	BeforeShutdown() error
	AfterShutdown() error
	// 监听需要关停的组件， 当事件发生时执行关停动作
	Start(gs GSInstance) error
}

// GSInstance 实现类
type GracefulShutdown struct {
	eh        ErrHandler
	name      string
	callbacks []ShutdownCallback
	managers  []ShutdownManager
	inOrder   bool
}

func NewGracefulShutdownInstance(name string) *GracefulShutdown {
	return &GracefulShutdown{
		name: name,
	}
}

// 使用按顺序关停
func (gs *GracefulShutdown) SetInOrder() {
	gs.inOrder = true
}

// 启动所有管理者的监听服务
func (gs *GracefulShutdown) Start() error {
	for _, mg := range gs.managers {
		if err := mg.Start(gs); err != nil {
			return err
		}
	}
	return nil
}

func (gs *GracefulShutdown) AddShutdownCallbacks(cbs ...ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, cbs...)
}

// 添加管理者
func (gs *GracefulShutdown) AddShutdownManagers(mgs ...ShutdownManager) {
	gs.managers = append(gs.managers, mgs...)
}

func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	// 如果按顺序关停
	if gs.inOrder {
		gs.StartShutdownInOrder(sm)
		return
	}

	gs.HandleErr(sm.BeforeShutdown())

	var wg sync.WaitGroup
	wg.Add(len(gs.callbacks))
	for _, callback := range gs.callbacks {
		go func(callback ShutdownCallback) {
			defer wg.Done()
			msg := fmt.Sprintf("%s is working", sm.GetName())
			gs.HandleErr(callback(msg))
		}(callback)
	}
	wg.Wait()

	gs.HandleErr(sm.AfterShutdown())
}

func (gs *GracefulShutdown) StartShutdownInOrder(sm ShutdownManager) {
	gs.HandleErr(sm.BeforeShutdown())
	for _, callback := range gs.callbacks {
		msg := fmt.Sprintf("%s is working", sm.GetName())
		gs.HandleErr(callback(msg))
	}
	gs.HandleErr(sm.AfterShutdown())
}

func (gs *GracefulShutdown) SetErrHandler(eh ErrHandler) {
	gs.eh = eh
}

func (gs *GracefulShutdown) HandleErr(err error) {
	if err != nil {
		if gs.eh == nil {
			log.Warn("ErrHandler is not set!, ignore err ")
			return
		}
		gs.eh.OnError(err)
	}
}
