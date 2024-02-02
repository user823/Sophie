package support

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type MyTask struct {
	GoService
}

func TestGoService(t *testing.T) {
	task := &MyTask{}
	task.Run = func(ctx context.Context) error {
		for task.Status() == STARTED {
			task.WaitForRunning(1 * time.Second)
			fmt.Println("running")
		}
		return nil
	}
	task.Start()
	time.Sleep(3 * time.Second)
	task.Shutdown(NOW)
}
