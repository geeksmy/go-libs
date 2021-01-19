package test_tool

import (
	"testing"

	"go.uber.org/zap"
)

func MockedZAPForTest(t *testing.T) *zap.Logger {
	return zap.L().With(zap.String("TEST", "mock for test"))
}
