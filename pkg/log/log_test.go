package log

import (
	"testing"
)

func TestEzap(t *testing.T) {
	logger := &ezapLogger{
		logger: std,
		env:    []any{"test-key", "test-value"},
	}
	logger.Info("something")
}
