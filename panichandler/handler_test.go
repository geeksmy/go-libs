package panichandler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestZapHandler_Handle(t *testing.T) {
	logger := zap.L().With(zap.String("logger_from", "unit_test"))
	assert.NotPanics(t, func() {
		defer ZapHandler(logger).Handle()
		panic("这是一个测试 PanicHandler 的 UT, 看见终端输出时不要 panic")
	})
}
