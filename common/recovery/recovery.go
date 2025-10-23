package recovery

import (
	"fmt"
	"rename-tool/common/applog"
	"rename-tool/setting/i18n"
	"runtime/debug"

	"golang.org/x/sys/windows"
)

// RecoverPanic 捕获并记录 panic，用于 GUI 程序（无控制台）
func RecoverPanic() {
	if r := recover(); r != nil {
		stack := string(debug.Stack())

		// 写入日志
		applog.Logger.Printf("[PANIC] %v\nStack:\n%s", r, stack)

		// 弹窗提示（因为没有控制台）
		message := fmt.Sprintf("%s\n\n%v", i18n.LogTr("ProgramCrashed"), r)
		windows.MessageBox(0, windows.StringToUTF16Ptr(message), windows.StringToUTF16Ptr("程序异常"), windows.MB_ICONERROR)

		// 退出程序
		windows.ExitProcess(1)
	}
}
