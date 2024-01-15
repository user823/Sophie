package shutdownmanagers

import (
	"github.com/user823/Sophie/pkg/shutdown"
	"os"
	"os/signal"
	"syscall"
)

// 默认关停管理者
// 监听
const name = "DefaultShutdownManager"

type DefaultShutdownManager struct {
	signals []os.Signal
}

func NewDefaultShutdownManager() *DefaultShutdownManager {
	return &DefaultShutdownManager{
		signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
	}
}

func (d *DefaultShutdownManager) Start(gs shutdown.GSInstance) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, d.signals...)
		<-c
		gs.StartShutdown(d)
	}()
	return nil
}

func (d *DefaultShutdownManager) GetName() string { return name }

func (d *DefaultShutdownManager) BeforeShutdown() error {
	return nil
}

func (d *DefaultShutdownManager) AfterShutdown() error {
	return nil
}
