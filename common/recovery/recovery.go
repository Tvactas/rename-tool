package recovery

import (
	"fmt"
	"rename-tool/common/log"
)

// RecoverPanic 用于优雅地处理程序中的panic
// 它会捕获panic并将错误信息记录到日志中
func RecoverPanic() {
	if r := recover(); r != nil {
		log.LogError(fmt.Errorf("panic: %v", r))
	}
}
