// package panichandler 方便用来在搜集 panic 的日志,
// example:
//  go func() {
//      defer ZapHandler(logger).Handler
//      // your function code here
//  }
package panichandler

import (
	"fmt"
	"runtime/debug"

	"go.uber.org/zap"
)

type zapHandler struct {
	logger *zap.Logger
}

func (h zapHandler) Handle() {
	if r := recover(); r != nil {
		fmt.Println("got panic: ", r)
		fmt.Println("stacktrace from panic: \n", string(debug.Stack()))
		h.logger.Error("got panic", zap.ByteString("stacktrace", debug.Stack()))
	}
}

// ZapHandler 使用 zap 搜集 panic 信息
// @param logger  logger == nil 会使用全局 Logger
func ZapHandler(logger *zap.Logger) zapHandler {
	if logger == nil {
		logger = zap.L()
	}
	return zapHandler{logger: logger}
}
