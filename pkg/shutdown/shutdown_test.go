package shutdown

import (
	"fmt"
	"testing"
)

func TestShutdown(t *testing.T) {
	c := make(chan struct{})
	mg := DefaultShutdownManager()
	gs := NewGracefulShutdownInstance("test")
	fn := func(msg string) error {
		fmt.Println(msg)
		c <- struct{}{}
		return nil
	}
	gs.AddShutdownCallbacks(fn)
	gs.AddShutdownManagers(mg)
	gs.SetErrHandler(&EmptyErrHandler{})
	gs.Start()
	<-c
}
