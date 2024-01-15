package shutdownmanagers

import (
	"fmt"
	"github.com/user823/Sophie/pkg/shutdown"
	"testing"
)

func TestShutdown(t *testing.T) {
	mg := NewDefaultShutdownManager()
	gs := shutdown.NewGracefulShutdownInstance("test")
	fn := func(msg string) error {
		fmt.Println(msg)
		return nil
	}
	gs.AddShutdownCallbacks(shutdown.ShutdownHelper(fn))
	gs.AddShutdownManagers(mg)
	gs.SetErrHandler(shutdown.DefaultErrHandler{})
	gs.Start()
}
