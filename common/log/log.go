package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"rename-tool/common/dirpath"
	"rename-tool/setting/config"
	"time"
)

// 修改日志记录函数
func LogError(err error) {
	if err == nil {
		return
	}

	logPath := GetErrorLogPath()
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprint("[", timestamp, "] ", err, "\n")
	if _, err := f.WriteString(logEntry); err != nil {
		// 记录写入错误，但不返回错误以避免循环
		fmt.Fprintf(os.Stderr, "Failed to write to error log: %v\n", err)
	}
}

// 获取错误日志文件路径
func GetErrorLogPath() string {
	appDir := dirpath.GetAppDataDir()
	logDir := filepath.Join(appDir, config.LogDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return filepath.Join(appDir, "error.log")
	}
	return filepath.Join(logDir, "error.log")
}

// 获取日志文件路径
func GetLogPath() string {
	appDir := dirpath.GetAppDataDir()
	logDir := filepath.Join(appDir, config.LogDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return filepath.Join(appDir, "rename.log")
	}
	return filepath.Join(logDir, "rename.log")
}

// getLogWriter returns a writer to a log file only (no console)
func GetLogWriter() io.Writer {
	var logFile string

	// Write user directories first
	if home, err := os.UserHomeDir(); err == nil {
		logFile = filepath.Join(home, "tvacats_rename.log")
	} else {

		// If the user directory cannot be obtained, write to the D drive
		logFile = "D:\\tvacats_rename.log"
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		// If the file cannot be written, you can choose to panic
		panic("无法写入日志文件: " + err.Error())
	}

	return file
}
