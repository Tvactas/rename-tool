package admin

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"rename-tool/setting/i18n"

	"golang.org/x/sys/windows"
)

var (
	logger  *log.Logger
	logFile *os.File
)

// Initialize logger
func init() {
	logOutput := getLogWriter()
	logger = log.New(logOutput, "", log.LstdFlags)
	log.SetOutput(io.Discard)
}

// getLogWriter 按优先级尝试写入日志：用户目录 -> D盘 -> 丢弃
func getLogWriter() io.Writer {
	// 1. 尝试用户目录
	if home, err := os.UserHomeDir(); err == nil {
		logPath := filepath.Join(home, "tvacats_rename.log")
		if file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
			logFile = file
			return file
		}
	}

	// 2. 尝试 D 盘
	logPath := "D:\\tvacats_rename.log"
	if file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		logFile = file
		return file
	}

	// 3. 都失败则不保存
	return io.Discard
}

// Close 关闭日志文件（在 main 函数用 defer 调用）
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// Log 记录格式化日志
func Log(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// IsAdmin 检查是否以管理员身份运行
func IsAdmin() bool {
	sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		Log("%s: %v", i18n.LogTr("CreateWellKnownSidFail"), err)
		return false
	}

	token := windows.GetCurrentProcessToken()
	// GetCurrentProcessToken 返回伪句柄，无需 Close

	isMember, err := token.IsMember(sid)
	if err != nil {
		Log("%s: %v", i18n.LogTr("CheckIsMember"), err)
		return false
	}

	Log("%s: %v", i18n.LogTr("LoginIdentity"), isMember)
	return isMember
}
