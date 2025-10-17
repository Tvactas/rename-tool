package admin

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"rename-tool/setting/i18n"

	"golang.org/x/sys/windows"
)

var (
	initOnce sync.Once
	logger   *log.Logger
)

// Initialize logger
func init() {
	logOutput := getLogWriter()

	// 不使用 log.Default()，直接用独立 logger
	logger = log.New(logOutput, "", log.LstdFlags)

	// 可选：防止默认 log 打印到控制台
	log.SetOutput(io.Discard)
}

// getLogWriter returns a writer to a log file only (no console)
func getLogWriter() io.Writer {
	var logFile string

	// 优先写用户目录
	if home, err := os.UserHomeDir(); err == nil {
		logFile = filepath.Join(home, "tvacats_rename.log")
	} else {
		// 如果取不到用户目录，写到 D 盘
		logFile = "D:\\tvacats_rename.log"
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		// 如果无法写文件，可以选择 panic 或 fallback
		panic("无法写入日志文件: " + err.Error())
	}

	return file
}

// Log logs a formatted message using the package logger
func Log(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// IsAdmin returns true if running as administrator, false otherwise
func IsAdmin() bool {
	var isAdmin bool
	initOnce.Do(func() {
		sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
		if err != nil {
			Log("%s: %v", i18n.LogTr("CreateWellKnownSidFail"), err)
			return
		}

		token := windows.GetCurrentProcessToken()
		defer token.Close()

		isMember, err := token.IsMember(sid)
		if err != nil {
			Log("%s: %v", i18n.LogTr("CheckIsMember"), err)
			return
		}

		isAdmin = isMember
		Log("%s: %v", i18n.LogTr("LoginIdentity"), isAdmin)
	})
	return isAdmin
}
