package main

import (
	"github.com/user823/Sophie/internal/schedule/worker"
)

func main() {
	worker.NewApp("sophie-schedule-worker").Run()
}
