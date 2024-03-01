package support

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGoService(t *testing.T) {
	task := &GoTask{}
	task.Run = func(ctx context.Context) (interface{}, error) {
		for task.Status() == STARTED {
			task.WaitForRunning(1 * time.Second)
			fmt.Println("running")
		}
		return 1, nil
	}
	task.Start()
	time.Sleep(3 * time.Second)
	task.Shutdown(NOW)

	// 异步获取任务结果
	fmt.Println(task.AsyncOutput())

	// 同步获取任务执行结果
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := task.Output(ctx)
	if result.OK {
		fmt.Println(result.Data)
	}
}
