package errors

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	s := WithCallers()
	fmt.Printf("%+v", s)
}

func WithCallers() *stack {
	return callers()
}
