package log

import (
	"testing"
)

func TestStructLog(t *testing.T) {
	Infof("something")
}

func TestWithValues(t *testing.T) {
	logger := WithValues("test-key", "test-value")
	logger.Infof("something")
}
