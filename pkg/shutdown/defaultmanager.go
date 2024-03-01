package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

// 默认关停管理者
// 监听
const name = "DefaultShutdownManager"

type PosixSignalShutdownManager struct {
	signals []os.Signal
}

func DefaultShutdownManager() *PosixSignalShutdownManager {
	return &PosixSignalShutdownManager{
		signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
	}
}

func (d *PosixSignalShutdownManager) Start(gs GSInstance) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, d.signals...)
		<-c
		gs.StartShutdown(d)
	}()
	return nil
}

func (d *PosixSignalShutdownManager) GetName() string { return name }

func (d *PosixSignalShutdownManager) BeforeShutdown() error {
	return nil
}

func (d *PosixSignalShutdownManager) AfterShutdown() error {
	return nil
}
