package support

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestManagetAll(t *testing.T) {
	// 启动5个 任务
	mg := NewTaskManager()
	for i := 0; i < 5; i++ {
		task := &GoTask{}
		task.Run = func(ctx context.Context) (interface{}, error) {
			return 1, nil
		}
		mg.Add(task)
	}
	results := mg.All()
	for i := range results {
		fmt.Println(results[i])
	}
}

func TestManagerAny(t *testing.T) {
	// 启动5个 任务
	mg := NewTaskManager()
	for i := 0; i < 5; i++ {
		task := &GoTask{}
		task.Run = func(ctx context.Context) (interface{}, error) {
			return 1, nil
		}
		mg.Add(task)
	}

	result := mg.Any()
	fmt.Println(result)
	fmt.Println(runtime.NumGoroutine())
	time.Sleep(3 * time.Second)
	fmt.Println(runtime.NumGoroutine())
}
