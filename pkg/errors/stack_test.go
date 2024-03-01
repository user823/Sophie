package errors

import (
	"fmt"
	"github.com/user823/Sophie/pkg/ds"
	"testing"
)

func TestStack(t *testing.T) {
	s := WithCallers()
	fmt.Printf("%+v", s)
}

func WithCallers() *ds.stack {
	return ds.callers()
}
